// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

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
	s.Handler.ServeHTTP(&statusCodeResponseWriter{
		ResponseWriter: w,
		request:        r,

		handlers: s.handlers,
	}, r)
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	request *http.Request

	handlers map[int]http.Handler

	didWrite  bool
	skipWrite bool
}

func (w *statusCodeResponseWriter) WriteHeader(code int) {
	handler, ok := w.handlers[code]
	if !ok || w.didWrite || w.skipWrite {
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

	handler.ServeHTTP(w.ResponseWriter, w.request)
}

func (w *statusCodeResponseWriter) Write(p []byte) (int, error) {
	if w.skipWrite {
		return len(p), nil
	}

	w.didWrite = true
	return w.ResponseWriter.Write(p)
}
