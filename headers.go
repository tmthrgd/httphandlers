// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SetHeaders sets multiple response headers.
func SetHeaders(h http.Handler, headers map[string]string) http.Handler {
	canonical := make(map[string]string, len(headers))

	for k, v := range headers {
		canonical[http.CanonicalHeaderKey(k)] = v
	}

	return &setHeaders{
		Handler: h,

		headers: canonical,
	}
}

type setHeaders struct {
	http.Handler

	headers map[string]string
}

func (sh *setHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hdr := w.Header()

	for k, v := range sh.headers {
		hdr[k] = []string{v}
	}

	sh.Handler.ServeHTTP(w, r)
}
