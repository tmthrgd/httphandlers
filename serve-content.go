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
func ServeBytes(name string, modtime time.Time, data []byte) http.Handler {
	return &serveBytes{name, modtime, data}
}

type serveBytes struct {
	name    string
	modtime time.Time
	data    []byte
}

func (sb *serveBytes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, sb.name, sb.modtime, bytes.NewReader(sb.data))
}

// ServeString returns a http.Handler that calls
// http.ServeContent with a strings.Reader.
func ServeString(name string, modtime time.Time, data string) http.Handler {
	return &serveString{name, modtime, data}
}

type serveString struct {
	name    string
	modtime time.Time
	data    string
}

func (sb *serveString) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, sb.name, sb.modtime, strings.NewReader(sb.data))
}
