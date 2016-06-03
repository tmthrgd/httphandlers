// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net"
	"net/http"
)

// HostRedirector redirects clients to the same
// request URL but with a different host.
type HostRedirector struct {
	// The host to redirect to.
	Host string

	// The HTTP status code to use when
	// redirecting, defaults to 301 Moved
	// Permanently.
	Code int
}

// ServeHTTP implements http.Handler.
func (hr *HostRedirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := *r.URL

	if url.Scheme == "" {
		if r.TLS != nil {
			url.Scheme = "https"
		} else {
			url.Scheme = "http"
		}
	}

	if _, port, err := net.SplitHostPort(r.Host); err == nil {
		url.Host = net.JoinHostPort(hr.Host, port)
	} else {
		url.Host = hr.Host
	}

	code := http.StatusMovedPermanently
	if hr.Code != 0 {
		code = hr.Code
	}

	http.Redirect(w, r, url.String(), code)
}
