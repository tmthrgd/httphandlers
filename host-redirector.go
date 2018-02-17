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

// HostRedirect returns a request handler that
// redirects each request it receives to the
// same url, but with a different host, using
// the given status code.
//
// The provided code should be in the 3xx range
// and is usually http.StatusMovedPermanently.
func HostRedirect(host string, code int) Handler {
	return &hostRedirector{host, code}
}

type hostRedirector struct {
	host string
	code int
}

func (hr *hostRedirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := *r.URL

	if r.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if port := (&url.URL{Host: r.Host}).Port(); port != "" {
		u.Host = net.JoinHostPort(hr.host, port)
	} else {
		u.Host = hr.host
	}

	http.Redirect(w, r, u.String(), hr.code)
}
