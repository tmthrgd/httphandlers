// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// NeverModified wraps a http.Handler and returns a
// 304 Not Modified HTTP status code to the browser
// for all If-Modified-Since conditional requests.
//
// It is intended for resources that are guaranteed
// to never change, primarily content addressable
// resources.
//
// See the article
//  'This browser tweak saved 60% of requests to Facebook'
//   https://code.facebook.com/posts/557147474482256/
// for an overview of how this method works.
func NeverModified(h http.Handler) Handler {
	return &neverModified{h}
}

type neverModified struct {
	h http.Handler
}

func (nm *neverModified) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ims := r.Header["If-Modified-Since"]; len(ims) != 0 && ims[0] != "" &&
		(r.Method == http.MethodGet || r.Method == http.MethodHead) {
		if _, err := http.ParseTime(ims[0]); err == nil {
			writeNotModified(w)
			return
		}
	}

	nm.h.ServeHTTP(w, r)
}
