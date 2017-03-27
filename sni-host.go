// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SNIHost sets (*http.Request).Host to the TLS
// SNI extension value if the request is TLS and
// the HTTP Host header was absent.
type SNIHost struct {
	http.Handler
}

// ServeHTTP implements http.Handler.
func (h *SNIHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "" && r.TLS != nil {
		r.Host = r.TLS.ServerName
	}

	h.Handler.ServeHTTP(w, r)
}
