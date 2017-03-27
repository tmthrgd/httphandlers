// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// ErrorMessage calls http.Error with the given
// HTTP status code and message.
func ErrorMessage(msg string, code int) http.Handler {
	return &errorHandler{msg, code}
}

type errorHandler struct {
	msg  string
	code int
}

func (h *errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, h.msg, h.code)
}
