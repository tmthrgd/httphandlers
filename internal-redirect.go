// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/url"
)

// parseURL is just url.Parse. It exists only so that
// url.Parse can be called in places where url is
// shadowed for godoc.
var parseURL = url.Parse

// InternalRedirect replaces the requests url with the
// given string. It is like http.Redirect but is handled
// internally.
func InternalRedirect(h http.Handler, url string) Handler {
	u, err := parseURL(url)
	if err != nil {
		panic(err)
	}

	if u.Fragment != "" {
		panic("handlers: fragment must be empty in InternalRedirect")
	}

	return &internalRedirect{h, u}
}

// InternalRedirectWrap returns a Middleware that calls
// InternalRedirect.
func InternalRedirectWrap(url string) Middleware {
	return func(h http.Handler) http.Handler {
		return InternalRedirect(h, url)
	}
}

type internalRedirect struct {
	h   http.Handler
	url *url.URL
}

func (ir *internalRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := new(http.Request)
	*req = *r

	url := new(url.URL)
	*url = *ir.url
	req.URL = url

	ir.h.ServeHTTP(w, req)
}
