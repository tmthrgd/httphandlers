// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

type stringWriter interface {
	WriteString(s string) (n int, err error)
}

type responseWriterFlusherSW interface {
	http.ResponseWriter
	http.Flusher
	stringWriter
}

type closeNotifyResponseWriter struct {
	responseWriterFlusherSW
	http.CloseNotifier
}

type hijackResponseWriter struct {
	responseWriterFlusherSW
	http.Hijacker
}

type pusherResponseWriter struct {
	responseWriterFlusherSW
	http.Pusher
}

type closeNotifyHijackResponseWriter struct {
	responseWriterFlusherSW
	http.CloseNotifier
	http.Hijacker
}

type closeNotifyPusherResponseWriter struct {
	responseWriterFlusherSW
	http.CloseNotifier
	http.Pusher
}
