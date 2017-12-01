// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

var logBufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Redeclared to avoid importing
// github.com/tmthrgd/go-server-push.
const sentinelH2Push = "X-H2-Push"

// AccessLog wraps a http.Handler and logs
// all HTTP requests to an io.Writer that
// defaults to os.Stderr.
//
// The log format is intended for human
// debugging and may not be stable.
func AccessLog(h http.Handler, out io.Writer) http.Handler {
	if out == nil {
		out = os.Stderr
	}

	return &accessLog{h, out}
}

type accessLog struct {
	h   http.Handler
	out io.Writer
}

func (al *accessLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	buf := logBufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	var scratch [20]byte

	buf.Write(start.AppendFormat(scratch[:0], "2006/01/02 15:04:05 "))
	buf.WriteString((&url.URL{Host: r.RemoteAddr}).Hostname())

	if r.TLS == nil {
		buf.WriteByte(' ')
	} else if vers := tlsVersionToLogName[r.TLS.Version]; vers != "" {
		buf.WriteString(vers)
	} else {
		buf.WriteString(" TLS:? ")
	}

	buf.WriteString(r.Proto)
	buf.WriteByte(' ')
	buf.WriteString(r.Method)

	uri := *r.URL
	uri.Host = r.Host

	if r.TLS != nil {
		uri.Scheme = "https"
	} else {
		uri.Scheme = "http"
	}

	buf.WriteByte(' ')
	buf.WriteString(uri.String())

	lw := &logResponseWriter{
		ResponseWriter: w,
	}

	var rw http.ResponseWriter = lw

	_, cok := w.(http.CloseNotifier)
	_, hok := w.(http.Hijacker)
	_, pok := w.(http.Pusher)

	switch {
	case cok && hok:
		hj := hijackLogResponseWriter{lw}
		rw = closeNotifyHijackLogResponseWriter{hj}
	case cok && pok:
		rw = closeNotifyPusherLogResponseWriter{lw}
	case cok:
		rw = closeNotifyLogResponseWriter{lw}
	case hok:
		rw = hijackLogResponseWriter{lw}
	case pok:
		rw = pusherLogResponseWriter{lw}
	}

	al.h.ServeHTTP(rw, r)

	if lw.code == 0 {
		lw.code = http.StatusOK
	}

	buf.WriteByte(' ')
	buf.Write(strconv.AppendInt(scratch[:0], int64(lw.code), 10))
	buf.WriteByte(' ')
	buf.Write(strconv.AppendInt(scratch[:0], lw.size, 10))
	buf.WriteByte(' ')
	buf.Write(strconv.AppendInt(scratch[:0], int64(time.Since(start)/time.Microsecond), 10))

	if r.TLS != nil && r.TLS.DidResume {
		buf.WriteString(" resumed")
	}

	if _, isPush := r.Header[sentinelH2Push]; isPush {
		buf.WriteString(" h2-pushed")
	}

	buf.WriteByte('\n')
	buf.WriteTo(al.out)

	logBufferPool.Put(buf)
}

type logResponseWriter struct {
	http.ResponseWriter

	code int
	size int64
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)

	if w.code == 0 {
		w.code = code
	}
}

func (w *logResponseWriter) Write(p []byte) (n int, err error) {
	if w.code == 0 {
		w.code = http.StatusOK
	}

	n, err = w.ResponseWriter.Write(p)
	w.size += int64(n)
	return
}

func (w *logResponseWriter) WriteString(s string) (n int, err error) {
	if w.code == 0 {
		w.code = http.StatusOK
	}

	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += int64(n)
	return
}

func (w *logResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// This struct is intentionally small (1 pointer wide) so as to
// fit inside an interface{} without causing an allocaction.
type hijackLogResponseWriter struct {
	*logResponseWriter
}

func (w hijackLogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := w.ResponseWriter.(http.Hijacker).Hijack()

	if err == nil && w.code == 0 {
		// The status will be StatusSwitchingProtocols if there was no
		// error and WriteHeader has not been called yet.
		w.code = http.StatusSwitchingProtocols
	}

	return conn, rw, err
}

var _ http.Hijacker = hijackLogResponseWriter{}

var tlsVersionToLogName = map[uint16]string{
	tls.VersionSSL30: " SSL3.0 ",
	tls.VersionTLS10: " TLS1.0 ",
	tls.VersionTLS11: " TLS1.1 ",
	tls.VersionTLS12: " TLS1.2 ",
	0x0304:           " TLS1.3 ",
	0x7f00 | 18:      " TLS1.3-d18 ",
	0x7f00 | 22:      " TLS1.3-d22 ",
}

type (
	// Each of these structs is intentionally small (1 pointer wide) so
	// as to fit inside an interface{} without causing an allocaction.
	closeNotifyLogResponseWriter       struct{ *logResponseWriter }
	pusherLogResponseWriter            struct{ *logResponseWriter }
	closeNotifyHijackLogResponseWriter struct{ hijackLogResponseWriter }
	closeNotifyPusherLogResponseWriter struct{ *logResponseWriter }
)

var (
	_ http.CloseNotifier = closeNotifyLogResponseWriter{}
	_ http.CloseNotifier = closeNotifyHijackLogResponseWriter{}
	_ http.CloseNotifier = closeNotifyPusherLogResponseWriter{}
	_ http.Hijacker      = closeNotifyHijackLogResponseWriter{}
	_ http.Pusher        = pusherLogResponseWriter{}
	_ http.Pusher        = closeNotifyPusherLogResponseWriter{}
)

func (w closeNotifyLogResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyHijackLogResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyPusherLogResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w pusherLogResponseWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}

func (w closeNotifyPusherLogResponseWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}
