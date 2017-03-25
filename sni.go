// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SNIHandler verifies that the TLS SNI extension
// matches the HTTP Host header. It also fills
// in (*http.Request).Host if it is blank and
// the request is TLS.
type SNIHandler struct {
	http.Handler

	// Mismatch is invoked on requests
	// that contain a mismatch betwee the
	// TLS SNI extension and the HTTP Host
	// header.
	//
	// If Mismatch is null, a 400 Bad
	// Request error will be returned instead.
	Mismatch http.Handler
}

// ServeHTTP implements http.Handler.
func (h *SNIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.TLS == nil || r.TLS.ServerName == "":
		h.Handler.ServeHTTP(w, r)
	case r.Host == "":
		r.Host = r.TLS.ServerName
		h.Handler.ServeHTTP(w, r)
	case r.TLS.ServerName == stripPort(r.Host):
		h.Handler.ServeHTTP(w, r)
	case h.Mismatch != nil:
		h.Mismatch.ServeHTTP(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
