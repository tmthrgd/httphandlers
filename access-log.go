// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bytes"
	"crypto/tls"
	"io"
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
	http.Handler

	out io.Writer
}

func (l *accessLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	buf := logBufferPool.Get().(*bytes.Buffer)
	defer logBufferPool.Put(buf)

	buf.Reset()

	buf.WriteString(start.Format("2006/01/02 15:04:05 "))
	buf.WriteString((&url.URL{Host: r.RemoteAddr}).Hostname())

	if r.TLS != nil {
		if vers := tlsVersionToLogName[r.TLS.Version]; vers != "" {
			buf.WriteString(vers)
		} else {
			buf.WriteString(" TLS:?")
		}
	}

	buf.WriteByte(' ')
	buf.WriteString(r.Proto)
	buf.WriteByte(' ')
	buf.WriteString(r.Method)

	if r.TLS != nil {
		buf.WriteString(" https://")
	} else {
		buf.WriteString(" http://")
	}

	buf.WriteString(r.Host)
	buf.WriteString(r.RequestURI)

	lw := &logResponseWriter{
		ResponseWriter: w,

		code: http.StatusOK,
	}
	l.Handler.ServeHTTP(lw, r)

	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(int64(lw.code), 10))
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(int64(lw.size), 10))
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(int64(time.Since(start)/time.Microsecond), 10))

	if r.TLS != nil && r.TLS.DidResume {
		buf.WriteString(" resumed")
	}

	if _, isPush := r.Header[sentinelH2Push]; isPush {
		buf.WriteString(" h2-pushed")
	}

	buf.WriteByte('\n')
	buf.WriteTo(l.out)
}

type logResponseWriter struct {
	http.ResponseWriter

	code int
	size int64
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.code = code
}

func (w *logResponseWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	w.size += int64(n)
	return
}

var tlsVersionToLogName = map[uint16]string{
	tls.VersionSSL30: " SSL3.0",
	tls.VersionTLS10: " TLS1.0",
	tls.VersionTLS11: " TLS1.1",
	tls.VersionTLS12: " TLS1.2",
	0x0304:           " TLS1.3",
	0x7f00 | 18:      " TLS1.3-d18",
}
