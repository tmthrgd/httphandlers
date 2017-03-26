// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net"
	"net/http"
	"net/url"
)

// RedirectToHTTPS redirects clients to the
// same URL but with the scheme set to https.
type RedirectToHTTPS struct {
	// Optionally specifies a host to
	// redirect to for clients that do
	// not set the HTTP Host header.
	//
	// If Host is an empty string, a 400
	// Bad Request error will be returned
	// instead.
	Host string

	// Optionally specifies a port to
	// add to the URL.
	Port string

	// The HTTP status code to use when
	// redirecting, defaults to 301 Moved
	// Permanently.
	Code int
}

// ServeHTTP implements http.Handler.
func (h *RedirectToHTTPS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := *r.URL
	u.Scheme = "https"

	if u.Host = (&url.URL{Host: r.Host}).Hostname(); u.Host == "" {
		if h.Host == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		u.Host = h.Host
	}

	switch h.Port {
	case "", "443", "https":
	default:
		u.Host = net.JoinHostPort(u.Host, h.Port)
	}

	code := http.StatusMovedPermanently
	if h.Code != 0 {
		code = h.Code
	}

	http.Redirect(w, r, u.String(), code)
}
