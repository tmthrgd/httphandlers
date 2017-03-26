// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// ErrorCode calls http.Error with the given
// HTTP status code. It uses http.StatusText
// for the message.
type ErrorCode int

// ServeHTTP implements http.Handler.
func (code ErrorCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(int(code)), int(code))
}
