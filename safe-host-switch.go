// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// SafeHostSwitch is a http.Handler that routes
// the request based on the Host header. It is
// thread safe and requites no external locks.
type SafeHostSwitch struct {
	m sync.Map

	// NotFound is invoked for hosts
	// that have not been added to the
	// host switch.
	NotFound http.Handler
}

// Add adds a http.Handler to the host switch.
//
// It returns an error if the host has already been
// added.
func (hs *SafeHostSwitch) Add(host string, h http.Handler) error {
	if _, dup := hs.m.LoadOrStore(host, h); dup {
		return fmt.Errorf("handlers: a handle is already registered for host %q", host)
	}

	return nil
}

// Remove removes a http.Handler from the host switch.
func (hs *SafeHostSwitch) Remove(host string) {
	hs.m.Delete(host)
}

// ServeHTTP implements http.Handler.
func (hs *SafeHostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := (&url.URL{Host: r.Host}).Hostname()

	if h, ok := hs.m.Load(host); ok {
		h.(http.Handler).ServeHTTP(w, r)
		return
	}

	if hs.NotFound != nil {
		hs.NotFound.ServeHTTP(w, r)
	} else {
		http.Error(w, forbiddenText, http.StatusForbidden)
	}
}
