// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

type responseWriterFlusher interface {
	http.ResponseWriter
	http.Flusher
}

type hijackResponseWriter struct {
	responseWriterFlusher
	http.Hijacker
}
