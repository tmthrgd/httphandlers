// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSNIHost(t *testing.T) {
	var host string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host = r.Host

		w.WriteHeader(999)
	})

	r := httptest.NewRequest(http.MethodGet, "http://example.org/path/to/file", nil)
	w := httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, "example.org", host)

	r = httptest.NewRequest(http.MethodGet, "https://example.org/path/to/file", nil)
	w = httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, "example.org", host)

	r = httptest.NewRequest(http.MethodGet, "https://example.org/path/to/file", nil)
	r.Host = ""
	w = httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, "example.org", host)

	r = httptest.NewRequest(http.MethodGet, "https://example.org/path/to/file", nil)
	r.Host, r.TLS.ServerName = "", ""
	w = httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, "", host)

	r = httptest.NewRequest(http.MethodGet, "https://example.org/path/to/file", nil)
	r.Host, r.TLS.ServerName = "", "example.net"
	w = httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, "example.net", host)
}

func TestSNIHostCopiesRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.Host = ""

	var req *http.Request
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
	})

	w := httptest.NewRecorder()
	SNIHost(h).ServeHTTP(w, r)

	assert.False(t, r == req, "SNIHost did not copy http.Request")
}
