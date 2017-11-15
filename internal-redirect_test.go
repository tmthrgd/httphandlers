// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalRedirect(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)

	var url *url.URL
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL
	})

	w := httptest.NewRecorder()
	InternalRedirect(h, "/new/path/to/other/file").ServeHTTP(w, r)

	assert.Equal(t, url.String(), "/new/path/to/other/file")

	w = httptest.NewRecorder()
	InternalRedirect(h, "/").ServeHTTP(w, r)

	assert.Equal(t, url.String(), "/")

	w = httptest.NewRecorder()
	InternalRedirect(h, "/?abc").ServeHTTP(w, r)

	assert.Equal(t, url.String(), "/?abc")
}

func TestInternalRedirectInvalidPanics(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	assert.Panics(t, func() {
		InternalRedirect(h, ":")
	})

	assert.Panics(t, func() {
		InternalRedirect(h, "/?abc#def")
	})
}

func TestInternalRedirectCopiesRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)

	var req *http.Request
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
	})

	w := httptest.NewRecorder()
	InternalRedirect(h, "/new/path/to/other/file").ServeHTTP(w, r)

	assert.False(t, r == req, "InternalRedirect did not copy http.Request")
}

func TestInternalRedirectCopiesURL(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)

	var url *url.URL
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL
	})

	w := httptest.NewRecorder()
	ir := InternalRedirect(h, "/new/path/to/other/file")
	ir.ServeHTTP(w, r)

	assert.False(t, ir.(*internalRedirect).url == url, "InternalRedirect did not copy url.URL")
}
