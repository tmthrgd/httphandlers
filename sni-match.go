// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/url"
)

// SNIMatch verifies that the TLS SNI extension
// matches the HTTP Host header.
//
// It does nothing for HTTP/2.0 to allow for
// connection coalescing.
//
// mismatch is invoked on requests that contain
// a mismatch between the TLS SNI extension and
// the HTTP Host header.
//
// If Mismatch is nil, a 400 Bad Request error
// will be returned instead.
func SNIMatch(h http.Handler, mismatch http.Handler) Handler {
	if mismatch == nil {
		mismatch = ErrorCode(http.StatusBadRequest)
	}

	return &sniMatch{h, mismatch}
}

// SNIMatchWrap returns a Middleware that calls SNIMatch.
func SNIMatchWrap(mismatch http.Handler) Middleware {
	return func(h http.Handler) http.Handler {
		return SNIMatch(h, mismatch)
	}
}

type sniMatch struct {
	h        http.Handler
	mismatch http.Handler
}

func (sm *sniMatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil || r.TLS.ServerName == "" || r.ProtoMajor == 2 ||
		r.TLS.ServerName == (&url.URL{Host: r.Host}).Hostname() {
		sm.h.ServeHTTP(w, r)
	} else {
		sm.mismatch.ServeHTTP(w, r)
	}
}
