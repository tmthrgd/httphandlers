// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SetHeader sets a header to a given value in the
// response.
func SetHeader(h http.Handler, name, value string) http.Handler {
	return &setHeader{
		Handler: h,
		name:    http.CanonicalHeaderKey(name),
		value:   value,
	}
}

type setHeader struct {
	http.Handler

	name, value string
}

func (sh *setHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()[sh.name] = []string{sh.value}

	sh.Handler.ServeHTTP(w, r)
}

// AddHeader adds a header with a given value to the
// response.
func AddHeader(h http.Handler, name, value string) http.Handler {
	return &addHeader{
		Handler: h,
		name:    http.CanonicalHeaderKey(name),
		value:   value,
	}
}

type addHeader struct {
	http.Handler

	name, value string
}

func (ah *addHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h[ah.name] = append(h[ah.name], ah.value)

	ah.Handler.ServeHTTP(w, r)
}

// DeleteHeader removes a header from the response.
func DeleteHeader(h http.Handler, name string) http.Handler {
	return &deleteHeader{
		Handler: h,
		name:    http.CanonicalHeaderKey(name),
	}
}

type deleteHeader struct {
	http.Handler

	name string
}

func (dh *deleteHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	delete(w.Header(), dh.name)

	dh.Handler.ServeHTTP(w, r)
}
