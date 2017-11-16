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

func TestHostSwitchAdd(t *testing.T) {
	var hs HostSwitch

	assert.NotPanics(t, func() {
		hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	})

	assert.PanicsWithValue(t, `handlers: a handle is already registered for host 'example.com'`, func() {
		hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	})
}

func TestHostSwitchNotFound(t *testing.T) {
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(999)
		}),
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "HostSwitch did not call NotFound")
}

func TestHostSwitchForbidden(t *testing.T) {
	var hs HostSwitch

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusForbidden))
}

func TestHostSwitch(t *testing.T) {
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(999)
		}),
	}

	hs.Add("example.com", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(997)
	}))
	hs.Add("example.org", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(998)
	}))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, 997, w.Code, "HostSwitch invoked incorrect http.Handler")
}
