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

func TestRedirectToHTTPS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)

	w := httptest.NewRecorder()
	(&RedirectToHTTPS{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "https://example.com/path/to/file", w.HeaderMap.Get("Location"))

	w = httptest.NewRecorder()
	(&RedirectToHTTPS{Code: http.StatusSeeOther}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "https://example.com/path/to/file", w.HeaderMap.Get("Location"))
}

func TestRedirectToHTTPSWithPort(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)

	w := httptest.NewRecorder()
	(&RedirectToHTTPS{Port: "1234"}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "https://example.com:1234/path/to/file", w.HeaderMap.Get("Location"))

	w = httptest.NewRecorder()
	(&RedirectToHTTPS{Port: "443"}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "https://example.com/path/to/file", w.HeaderMap.Get("Location"))

	w = httptest.NewRecorder()
	(&RedirectToHTTPS{Port: "https"}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "https://example.com/path/to/file", w.HeaderMap.Get("Location"))
}

func TestRedirectToHTTPSNoHost(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Host = ""

	w := httptest.NewRecorder()
	(&RedirectToHTTPS{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusBadRequest))

	w = httptest.NewRecorder()
	(&RedirectToHTTPS{Host: "example.org"}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "https://example.org/path/to/file", w.HeaderMap.Get("Location"))
}
