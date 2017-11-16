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

func TestHostRedirector(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	w := httptest.NewRecorder()
	HostRedirect("example.com", http.StatusSeeOther).ServeHTTP(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "http://example.com/path/to/file", w.HeaderMap.Get("Location"))

	r = httptest.NewRequest(http.MethodGet, "https://example.org/path/to/file", nil)
	w = httptest.NewRecorder()
	HostRedirect("example.com", http.StatusSeeOther).ServeHTTP(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "https://example.com/path/to/file", w.HeaderMap.Get("Location"))

	r = httptest.NewRequest(http.MethodGet, "https://example.org:1234/path/to/file", nil)
	w = httptest.NewRecorder()
	HostRedirect("example.com", http.StatusSeeOther).ServeHTTP(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "https://example.com:1234/path/to/file", w.HeaderMap.Get("Location"))
}
