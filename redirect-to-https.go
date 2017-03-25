// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net"
	"net/http"
)

// RedirectToHTTPS redirects clients to the
// same URL but with the scheme set to https.
type RedirectToHTTPS struct {
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

	if u.Host = stripPort(r.Host); u.Host == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
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
