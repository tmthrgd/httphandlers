// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNeverModifiedInvalidIMS(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	w := httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "NeverModified didn't invoke http.Handler")

	r = httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", "")
	w = httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "NeverModified didn't invoke http.Handler")

	r = httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", "!invalid!")
	w = httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "NeverModified didn't invoke http.Handler")
}

func TestNeverModifiedInvalidIMSHeaders(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	w := httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"Content-Type":   {"test1"},
		"Content-Length": {"test2"},
		"Last-Modified":  {"test3"},
	}
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "NeverModified didn't invoke http.Handler")
	assert.Equal(t, http.Header{
		"Content-Type":   {"test1"},
		"Content-Length": {"test2"},
		"Last-Modified":  {"test3"},
	}, w.Result().Header)
}

func TestNotModifiedValidIMS(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	now := time.Now()

	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", now.Format(http.TimeFormat))
	w := httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)

	r = httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", time.Unix(0, 0).Format(http.TimeFormat))
	w = httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)

	r = httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", time.Time{}.Format(http.TimeFormat))
	w = httptest.NewRecorder()
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)
}

func TestNotModifiedValidIMSHeaders(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	now := time.Now()

	r := httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", now.Format(http.TimeFormat))
	w := httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"Content-Type":   {"test1"},
		"Content-Length": {"test2"},
		"Last-Modified":  {"test3"},
	}
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)
	assert.Equal(t, http.Header{
		"Last-Modified": {"test3"},
	}, w.Result().Header)

	r = httptest.NewRequest(http.MethodGet, "/path/to/file", nil)
	r.Header.Set("If-Modified-Since", now.Format(http.TimeFormat))
	w = httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"Content-Type":   {"test1"},
		"Content-Length": {"test2"},
		"Last-Modified":  {"test3"},
		"Etag":           {"test4"},
	}
	NeverModified(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)
	assert.Equal(t, http.Header{
		"Etag": {"test4"},
	}, w.Result().Header)
}

func TestNotModifiedTimeFormats(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	now := time.Now()

	for _, format := range []string{
		http.TimeFormat,
		time.RFC850,
		time.ANSIC,
	} {
		r := httptest.NewRequest(http.MethodHead, "/path/to/file", nil)
		r.Header.Set("If-Modified-Since", now.Format(format))
		w := httptest.NewRecorder()
		NeverModified(h).ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotModified, w.Code,
			"NeverModified invoked http.Handler for time-format:%q", format)
	}
}

func TestNotModifiedMethods(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})
	now := time.Now()

	for _, method := range []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	} {
		r := httptest.NewRequest(method, "/path/to/file", nil)
		r.Header.Set("If-Modified-Since", now.Format(http.TimeFormat))
		w := httptest.NewRecorder()
		NeverModified(h).ServeHTTP(w, r)

		assert.Equal(t, 999, w.Code,
			"NeverModified didn't invoke http.Handler for method:%s", method)
	}
}
