// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// ErrorMessage calls http.Error with the given
// HTTP status code and message.
type ErrorMessage struct {
	Code    int
	Message string
}

// ServeHTTP implements http.Handler.
func (msg *ErrorMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, msg.Message, msg.Code)
}
