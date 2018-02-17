// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type statusCodeSwitchedError struct{}

func (statusCodeSwitchedError) Error() string {
	return "handlers: response switched to other http.Handler"
}

func (statusCodeSwitchedError) StatusCodeSwitched() {}

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
func StatusCodeSwitch(h http.Handler, handlers map[int]http.Handler) Handler {
	return &statusCodeSwitch{h, handlers}
}

// StatusCodeSwitchWrap returns a Middleware that calls
// StatusCodeSwitch.
func StatusCodeSwitchWrap(handlers map[int]http.Handler) Middleware {
	return func(h http.Handler) http.Handler {
		return StatusCodeSwitch(h, handlers)
	}
}

type statusCodeSwitch struct {
	h        http.Handler
	handlers map[int]http.Handler
}

func (sw *statusCodeSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc := &statusCodeResponseWriter{
		rw: w,

		handlers: sw.handlers,
	}

	var rw http.ResponseWriter = sc

	_, cok := w.(http.CloseNotifier)
	_, hok := w.(http.Hijacker)
	_, pok := w.(http.Pusher)

	switch {
	case cok && hok:
		rw = closeNotifyHijackStatusCodeResponseWriter{sc}
	case cok && pok:
		rw = closeNotifyPusherStatusCodeResponseWriter{sc}
	case cok:
		rw = closeNotifyStatusCodeResponseWriter{sc}
	case hok:
		rw = hijackStatusCodeResponseWriter{sc}
	case pok:
		rw = pusherStatusCodeResponseWriter{sc}
	}

	sw.h.ServeHTTP(rw, r)
	sc.writeHeaders()

	if sc.skipWrite {
		sw.handlers[int(sc.code)].ServeHTTP(w, r)
	}
}

type statusCodeResponseWriter struct {
	rw http.ResponseWriter

	handlers map[int]http.Handler

	headers http.Header

	didWrite  bool
	skipWrite bool
	code      int32
}

func (w *statusCodeResponseWriter) Header() http.Header {
	if w.didWrite && !w.skipWrite {
		return w.rw.Header()
	}

	if w.headers == nil {
		w.headers = cloneHeader(w.rw.Header())
	}

	return w.headers
}

func (w *statusCodeResponseWriter) writeHeaders() {
	if w.headers == nil || w.skipWrite {
		return
	}

	hdr := w.rw.Header()
	for k := range hdr {
		delete(hdr, k)
	}
	for k, vv := range w.headers {
		hdr[k] = vv
	}
}

func (w *statusCodeResponseWriter) WriteHeader(code int) {
	if w.skipWrite || w.didWrite {
		return
	}

	if _, ok := w.handlers[code]; ok && int(int32(code)) == code {
		w.skipWrite = true
		w.code = int32(code)
		return
	}

	w.didWrite = true
	w.writeHeaders()
	w.rw.WriteHeader(code)
}

func (w *statusCodeResponseWriter) Write(p []byte) (int, error) {
	if w.skipWrite {
		return 0, statusCodeSwitchedError{}
	}

	w.didWrite = true
	return w.rw.Write(p)
}

func (w *statusCodeResponseWriter) WriteString(s string) (int, error) {
	if w.skipWrite {
		return 0, statusCodeSwitchedError{}
	}

	w.didWrite = true
	return io.WriteString(w.rw, s)
}

func (w *statusCodeResponseWriter) Flush() {
	if w.skipWrite {
		return
	}

	if f, ok := w.rw.(http.Flusher); ok {
		f.Flush()
	}
}

type (
	// Each of these structs is intentionally small (1 pointer wide) so
	// as to fit inside an interface{} without causing an allocaction.
	closeNotifyStatusCodeResponseWriter       struct{ *statusCodeResponseWriter }
	hijackStatusCodeResponseWriter            struct{ *statusCodeResponseWriter }
	pusherStatusCodeResponseWriter            struct{ *statusCodeResponseWriter }
	closeNotifyHijackStatusCodeResponseWriter struct{ *statusCodeResponseWriter }
	closeNotifyPusherStatusCodeResponseWriter struct{ *statusCodeResponseWriter }
)

var (
	_ http.CloseNotifier = closeNotifyStatusCodeResponseWriter{}
	_ http.CloseNotifier = closeNotifyHijackStatusCodeResponseWriter{}
	_ http.CloseNotifier = closeNotifyPusherStatusCodeResponseWriter{}
	_ http.Hijacker      = hijackStatusCodeResponseWriter{}
	_ http.Hijacker      = closeNotifyHijackStatusCodeResponseWriter{}
	_ http.Pusher        = pusherStatusCodeResponseWriter{}
	_ http.Pusher        = closeNotifyPusherStatusCodeResponseWriter{}
)

func (w closeNotifyStatusCodeResponseWriter) CloseNotify() <-chan bool {
	return w.rw.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyHijackStatusCodeResponseWriter) CloseNotify() <-chan bool {
	return w.rw.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyPusherStatusCodeResponseWriter) CloseNotify() <-chan bool {
	return w.rw.(http.CloseNotifier).CloseNotify()
}

func (w hijackStatusCodeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.skipWrite {
		return nil, nil, http.ErrNotSupported
	}

	return w.rw.(http.Hijacker).Hijack()
}

func (w closeNotifyHijackStatusCodeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.skipWrite {
		return nil, nil, http.ErrNotSupported
	}

	return w.rw.(http.Hijacker).Hijack()
}

func (w pusherStatusCodeResponseWriter) Push(target string, opts *http.PushOptions) error {
	if w.skipWrite {
		return http.ErrNotSupported
	}

	return w.rw.(http.Pusher).Push(target, opts)
}

func (w closeNotifyPusherStatusCodeResponseWriter) Push(target string, opts *http.PushOptions) error {
	if w.skipWrite {
		return http.ErrNotSupported
	}

	return w.rw.(http.Pusher).Push(target, opts)
}
