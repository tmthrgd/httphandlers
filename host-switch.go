// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/url"
)

// HostSwitch is a http.Handler that routes
// the request based on the Host header.
type HostSwitch struct {
	m map[string]http.Handler

	// NotFound is invoked for hosts
	// that have not been added to the
	// host switch.
	NotFound http.Handler
}

// Add adds a http.Handler to the host switch.
//
// It panics if the host has already been added.
func (hs *HostSwitch) Add(host string, h http.Handler) {
	if hs.m == nil {
		hs.m = make(map[string]http.Handler)
	}

	if _, dup := hs.m[host]; dup {
		panic("handlers: a handle is already registered for host '" + host + "'")
	}

	hs.m[host] = h
}

// ServeHTTP implements http.Handler.
func (hs *HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := (&url.URL{Host: r.Host}).Hostname()

	if h, ok := hs.m[host]; ok {
		h.ServeHTTP(w, r)
		return
	}

	if hs.NotFound != nil {
		hs.NotFound.ServeHTTP(w, r)
	} else {
		http.Error(w, forbiddenText, http.StatusForbidden)
	}
}

var forbiddenText = http.StatusText(http.StatusForbidden)
