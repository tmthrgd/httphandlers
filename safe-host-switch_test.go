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

func TestSafeHostSwitchAdd(t *testing.T) {
	var hs SafeHostSwitch

	assert.NoError(t, hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))
	assert.NoError(t, hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))

	assert.Error(t, hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))
}

func TestSafeHostSwitchRemove(t *testing.T) {
	var hs SafeHostSwitch

	assert.NoError(t, hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))
	assert.NoError(t, hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))

	hs.Remove("example.com")

	assert.NoError(t, hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))
	assert.Error(t, hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))
}

func TestSafeHostSwitchNotFound(t *testing.T) {
	hs := &SafeHostSwitch{
		NotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(999)
		}),
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "SafeHostSwitch did not call NotFound")
}

func TestSafeHostSwitchForbidden(t *testing.T) {
	var hs SafeHostSwitch

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusForbidden))
}

func TestSafeHostSwitch(t *testing.T) {
	hs := &SafeHostSwitch{
		NotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(999)
		}),
	}

	assert.NoError(t, hs.Add("example.com", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(997)
	})))
	assert.NoError(t, hs.Add("example.org", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(998)
	})))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, 997, w.Code, "SafeHostSwitch invoked incorrect http.Handler")
}
