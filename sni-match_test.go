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
	"github.com/stretchr/testify/require"
)

func TestSNIMatch(t *testing.T) {
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(998)
	})

	r := httptest.NewRequest(http.MethodGet, "http://example.com/path/to/file", nil)
	w := httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "SNIMatch invoked incorrect http.Handler")

	r = httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.TLS.ServerName = ""
	w = httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "SNIMatch invoked incorrect http.Handler")

	r = httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.ProtoMajor = 2
	w = httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "SNIMatch invoked incorrect http.Handler")

	r = httptest.NewRequest(http.MethodGet, "https://example.com:1234/path/to/file", nil)
	// httptest.NewRequest does not properly strip the
	// port before setting r.TLS.ServerName.
	r.Host, r.TLS.ServerName = "example.com:1234", "example.com"
	w = httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "SNIMatch invoked incorrect http.Handler")
}

func TestSNIMatchMismatch(t *testing.T) {
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(998)
	})

	r := httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.Host = "example.org"
	w := httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 998, w.Code, "SNIMatch invoked incorrect http.Handler")

	r = httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.TLS.ServerName = "example.org"
	w = httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 998, w.Code, "SNIMatch invoked incorrect http.Handler")

	r = httptest.NewRequest(http.MethodGet, "https://example.com/path/to/file", nil)
	r.Host = "example.org:1234"
	w = httptest.NewRecorder()
	SNIMatch(h1, h2).ServeHTTP(w, r)

	assert.Equal(t, 998, w.Code, "SNIMatch invoked incorrect http.Handler")
}

func TestSNIMatchMismatchNoHandler(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	sni := SNIMatch(h, nil)

	require.IsType(t, (*sniMatch)(nil), sni)
	require.IsType(t, (*errorHandler)(nil), sni.(*sniMatch).mismatch)
	assert.Equal(t, http.StatusBadRequest, sni.(*sniMatch).mismatch.(*errorHandler).code)
}
