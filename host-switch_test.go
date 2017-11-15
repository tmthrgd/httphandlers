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
	calledNotFound := false
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	hs.ServeHTTP(httptest.NewRecorder(), r)

	assert.True(t, calledNotFound, "HostSwitch did not call NotFound")
}

func TestHostSwitchForbidden(t *testing.T) {
	var hs HostSwitch

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	w := httptest.NewRecorder()
	hs.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusForbidden)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusForbidden))
}

func TestHostSwitch(t *testing.T) {
	calledNotFound := false
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	calledExampleCom := false
	hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleCom = true
	}))

	calledExampleOrg := false
	hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleOrg = true
	}))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	hs.ServeHTTP(httptest.NewRecorder(), r)

	assert.False(t, calledNotFound, "HostSwitch did not call correct handler: NotFound")
	assert.True(t, calledExampleCom, "HostSwitch did not call correct handler: example.com")
	assert.False(t, calledExampleOrg, "HostSwitch did not call correct handler: example.org")
}
