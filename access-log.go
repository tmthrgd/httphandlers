// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Redeclared to avoid importing
// github.com/tmthrgd/go-server-push.
const sentinelH2Push = "X-H2-Push"

var stdLogger = log.New(os.Stderr, "", 0)

// AccessLogHandler logs HTTP requests to a
// *log.Logger.
type AccessLogHandler struct {
	http.Handler

	// The log to write log entries to.
	// Defaults to os.Stderr with no flags.
	AccessLog *log.Logger

	// The format string to use when logging
	// request start times. Defaults to
	// 2006/01/02 15:04:05.
	DateFormat string
}

// ServeHTTP implements http.Handler.
func (h *AccessLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	u := *r.URL
	u.Host = r.Host

	if r.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	lw := &logResponseWriter{
		ResponseWriter: w,

		code: http.StatusOK,
	}
	h.Handler.ServeHTTP(lw, r)

	var tlsVers, resumed string
	if r.TLS != nil {
		if tlsVers = tlsVersionToLogName[r.TLS.Version]; tlsVers == "" {
			tlsVers = " TLS:?"
		}

		if r.TLS.DidResume {
			resumed = " resumed"
		}
	}

	var pushed string
	if _, isPush := r.Header[sentinelH2Push]; isPush {
		pushed = " h2-pushed"
	}

	dateFormat := "2006/01/02 15:04:05"
	if h.DateFormat != "" {
		dateFormat = h.DateFormat
	}

	logger := stdLogger
	if h.AccessLog != nil {
		logger = h.AccessLog
	}

	logger.Printf("%s %s%s %s %s %s %d %d %d%s%s\n",
		start.Format(dateFormat),
		(&url.URL{Host: r.RemoteAddr}).Hostname(),
		tlsVers,
		r.Proto,
		r.Method,
		u.String(),
		lw.code,
		lw.size,
		time.Since(start)/time.Microsecond,
		resumed,
		pushed)
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
