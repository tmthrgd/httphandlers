// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"io"
	"net/http"
)

// StatusCodeSwitch intercepts calls to
// http.ResponseWriter.WriteHeader and redirects
// the request to a http.Handler based on the
// response status code.
//
// handlers is a map of HTTP status code
// (for example http.StatusNotFound) to
// a http.Handler to use for the response.
//
// It can be used with ServeError to statically
// render pretty error pages.
func StatusCodeSwitch(h http.Handler, handlers map[int]http.Handler) http.Handler {
	return &statusCodeSwitch{h, handlers}
}

type statusCodeSwitch struct {
	http.Handler

	handlers map[int]http.Handler
}

func (s *statusCodeSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc := &statusCodeResponseWriter{
		ResponseWriter: w,
		req:            r,

		handlers: s.handlers,
	}

	var rw http.ResponseWriter = sc

	c, cok := w.(http.CloseNotifier)
	h, hok := w.(http.Hijacker)
	p, pok := w.(http.Pusher)

	switch {
	case cok && hok:
		rw = &closeNotifyHijackResponseWriter{sc, c, h}
	case cok && pok:
		rw = &closeNotifyPusherResponseWriter{sc, c, p}
	case cok:
		rw = &closeNotifyResponseWriter{sc, c}
	case hok:
		rw = &hijackResponseWriter{sc, h}
	case pok:
		rw = &pusherResponseWriter{sc, p}
	}

	s.Handler.ServeHTTP(rw, r)
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	req *http.Request

	handlers map[int]http.Handler

	didWrite  bool
	skipWrite bool
}

func (w *statusCodeResponseWriter) WriteHeader(code int) {
	if w.skipWrite {
		return
	}

	handler, ok := w.handlers[code]
	if !ok || w.didWrite {
		w.ResponseWriter.WriteHeader(code)
		return
	}

	w.skipWrite = true

	h := w.Header()
	delete(h, "Cache-Control")
	delete(h, "Etag")
	delete(h, "Last-Modified")
	delete(h, "Content-Encoding")
	delete(h, "Content-Length")
	delete(h, "Content-Type")

	handler.ServeHTTP(w.ResponseWriter, w.req)
}

func (w *statusCodeResponseWriter) Write(p []byte) (int, error) {
	if w.skipWrite {
		return len(p), nil
	}

	w.didWrite = true
	return w.ResponseWriter.Write(p)
}

func (w *statusCodeResponseWriter) WriteString(s string) (int, error) {
	if w.skipWrite {
		return len(s), nil
	}

	w.didWrite = true
	return io.WriteString(w.ResponseWriter, s)
}

func (w *statusCodeResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
