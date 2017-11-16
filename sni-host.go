// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SNIHost sets (*http.Request).Host to the TLS
// SNI extension value if the request is TLS and
// the HTTP Host header was absent.
func SNIHost(h http.Handler) http.Handler {
	return &sniHost{h}
}

type sniHost struct {
	h http.Handler
}

func (sh *sniHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "" && r.TLS != nil && r.TLS.ServerName != "" {
		rr := *r
		rr.Host = r.TLS.ServerName
		r = &rr
	}

	sh.h.ServeHTTP(w, r)
}
