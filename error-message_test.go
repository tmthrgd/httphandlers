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

func TestErrorCode(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	ErrorCode(http.StatusNotFound).ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusNotFound))

	w = httptest.NewRecorder()
	ErrorCode(http.StatusBadGateway).ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusBadGateway)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusBadGateway))
}

func TestErrorCodeEqual(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	ErrorCode(http.StatusNotFound).ServeHTTP(w1, r)
	http.Error(w2, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	assert.Equal(t, w1.Result(), w2.Result())
}

func TestErrorMessage(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	ErrorMessage("test1", http.StatusNotFound).ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Contains(t, w.Body.String(), "test1")

	w = httptest.NewRecorder()
	ErrorMessage("test2", http.StatusBadGateway).ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusBadGateway)
	assert.Contains(t, w.Body.String(), "test2")
}

func TestErrorMessageEqual(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	ErrorMessage("test", http.StatusNotFound).ServeHTTP(w1, r)
	http.Error(w2, "test", http.StatusNotFound)

	assert.Equal(t, w1.Result(), w2.Result())
}
