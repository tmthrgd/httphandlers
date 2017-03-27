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
// It does nothing for HTTP/2.0 to allow
// for connection coalescing.
type SNIMatch struct {
	http.Handler

	// Mismatch is invoked on requests
	// that contain a mismatch betwee the
	// TLS SNI extension and the HTTP Host
	// header.
	//
	// If Mismatch is nil, a 400 Bad
	// Request error will be returned instead.
	Mismatch http.Handler
}

// ServeHTTP implements http.Handler.
func (h *SNIMatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.TLS == nil || r.TLS.ServerName == "" || r.ProtoMajor == 2 ||
		r.TLS.ServerName == (&url.URL{Host: r.Host}).Hostname():
		h.Handler.ServeHTTP(w, r)
	case h.Mismatch != nil:
		h.Mismatch.ServeHTTP(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
