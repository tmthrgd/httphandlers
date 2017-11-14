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
	calledNotFound := false
	hs := &SafeHostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	hs.ServeHTTP(httptest.NewRecorder(), &http.Request{Host: "example.com"})

	assert.True(t, calledNotFound, "SafeHostSwitch did not call NotFound")
}

func TestSafeHostSwitch(t *testing.T) {
	calledNotFound := false
	hs := &SafeHostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	calledExampleCom := false
	assert.NoError(t, hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleCom = true
	})))

	calledExampleOrg := false
	assert.NoError(t, hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleOrg = true
	})))

	hs.ServeHTTP(httptest.NewRecorder(), &http.Request{Host: "example.com"})

	assert.False(t, calledNotFound, "HostSwitch did not call correct handler: NotFound")
	assert.True(t, calledExampleCom, "HostSwitch did not call correct handler: example.com")
	assert.False(t, calledExampleOrg, "HostSwitch did not call correct handler: example.org")
}
