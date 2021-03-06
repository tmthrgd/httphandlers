// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SetHeaders sets multiple response headers.
func SetHeaders(h http.Handler, headers map[string]string) Handler {
	canonical := make(map[string]string, len(headers))

	for k, v := range headers {
		canonical[http.CanonicalHeaderKey(k)] = v
	}

	// Always give preference to any header that appears
	// in canonical form in the headers map. The selection
	// between two separate headers, both of which are in
	// non-canonical form, is undefined.
	for k := range canonical {
		if v, ok := headers[k]; ok {
			canonical[k] = v
		}
	}

	return &setHeaders{h, canonical}
}

// SetHeadersWrap returns a Middleware that calls
// SetHeaders.
func SetHeadersWrap(headers map[string]string) Middleware {
	return func(h http.Handler) http.Handler {
		return SetHeaders(h, headers)
	}
}

type setHeaders struct {
	h       http.Handler
	headers map[string]string
}

func (sh *setHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hdr := w.Header()

	for k, v := range sh.headers {
		hdr[k] = []string{v}
	}

	sh.h.ServeHTTP(w, r)
}
