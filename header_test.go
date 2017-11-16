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

func TestSetHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	w := httptest.NewRecorder()
	SetHeader(h, "X-Test", "test").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	SetHeader(h, "x-test", "test").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	h1 := SetHeader(h, "X-Test", "test1")
	SetHeader(h1, "x-test", "test2").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test1"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	h1 = SetHeader(h, "X-Test1", "test1")
	SetHeader(h1, "X-Test2", "test2").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test1": {"test1"},
		"X-Test2": {"test2"},
	}, w.Result().Header)
}

func TestAddHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	w := httptest.NewRecorder()
	AddHeader(h, "X-Test", "test").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	AddHeader(h, "x-test", "test").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	h1 := AddHeader(h, "X-Test", "test1")
	AddHeader(h1, "x-test", "test2").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test": {"test2", "test1"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	h1 = AddHeader(h, "X-Test1", "test1")
	AddHeader(h1, "X-Test2", "test2").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test1": {"test1"},
		"X-Test2": {"test2"},
	}, w.Result().Header)
}

func TestDeleteHeader(t *testing.T) {
	newRecorder := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		w.HeaderMap = http.Header{
			"X-Test1": {"test1", "test2"},
			"X-Test2": {"test1", "test2"},
		}
		return w
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	w := newRecorder()
	DeleteHeader(h, "X-Test1").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test2": {"test1", "test2"},
	}, w.Result().Header)

	w = newRecorder()
	DeleteHeader(h, "x-test1").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Test2": {"test1", "test2"},
	}, w.Result().Header)

	w = newRecorder()
	h1 := DeleteHeader(h, "X-Test1")
	DeleteHeader(h1, "x-test2").ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{}, w.Result().Header)
}
