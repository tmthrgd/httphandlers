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

func TestSetHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	w := httptest.NewRecorder()
	SetHeaders(h, map[string]string{
		"X-Test": "test",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	SetHeaders(h, map[string]string{
		"x-test": "test",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	SetHeaders(h, map[string]string{
		"X-Test": "test1",
		"x-test": "test2",
		"x-Test": "test3",
		"X-test": "test4",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test1"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	SetHeaders(h, map[string]string{
		"X-Test1": "test1",
		"x-test2": "test2",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test1": {"test1"},
		"X-Test2": {"test2"},
	}, w.Result().Header)
}
