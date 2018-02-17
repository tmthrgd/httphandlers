// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

// ServeBytes returns a http.Handler that calls
// http.ServeContent with a bytes.Reader.
func ServeBytes(name string, modtime time.Time, content []byte) Handler {
	return &serveBytes{name, modtime, content}
}

type serveBytes struct {
	name    string
	modtime time.Time
	content []byte
}

func (s *serveBytes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, s.name, s.modtime, bytes.NewReader(s.content))
}

// ServeString returns a http.Handler that calls
// http.ServeContent with a strings.Reader.
func ServeString(name string, modtime time.Time, content string) Handler {
	return &serveString{name, modtime, content}
}

type serveString struct {
	name    string
	modtime time.Time
	content string
}

func (s *serveString) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, s.name, s.modtime, strings.NewReader(s.content))
}
