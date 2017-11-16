// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SetHeader sets a header to a given value in the
// response.
func SetHeader(h http.Handler, name, value string) http.Handler {
	return &setHeader{h, http.CanonicalHeaderKey(name), value}
}

type setHeader struct {
	h     http.Handler
	name  string
	value string
}

func (sh *setHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()[sh.name] = []string{sh.value}

	sh.h.ServeHTTP(w, r)
}

// AddHeader adds a header with a given value to the
// response.
func AddHeader(h http.Handler, name, value string) http.Handler {
	return &addHeader{h, http.CanonicalHeaderKey(name), value}
}

type addHeader struct {
	h     http.Handler
	name  string
	value string
}

func (ah *addHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h[ah.name] = append(h[ah.name], ah.value)

	ah.h.ServeHTTP(w, r)
}

// DeleteHeader removes a header from the response.
func DeleteHeader(h http.Handler, name string) http.Handler {
	return &deleteHeader{h, http.CanonicalHeaderKey(name)}
}

type deleteHeader struct {
	h    http.Handler
	name string
}

func (dh *deleteHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	delete(w.Header(), dh.name)

	dh.h.ServeHTTP(w, r)
}
