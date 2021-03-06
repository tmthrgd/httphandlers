// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"strconv"
)

// ServeError returns a http.Handler that serves
// content with a given HTTP status code.
//
// If mimeType is empty, it will be sniffed from
// content.
func ServeError(code int, content []byte, mimeType string) Handler {
	if mimeType == "" {
		mimeType = http.DetectContentType(content)
	}

	return &serveError{
		content,
		strconv.FormatInt(int64(len(content)), 10),
		mimeType,
		code,
	}
}

type serveError struct {
	content []byte
	size    string
	mime    string
	code    int
}

func (se *serveError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()

	if h.Get("Content-Encoding") == "" {
		h.Set("Content-Length", se.size)
	}

	if _, hasType := h["Content-Type"]; !hasType {
		h.Set("Content-Type", se.mime)
	}

	w.WriteHeader(se.code)

	if r.Method != http.MethodHead {
		w.Write(se.content)
	}
}
