// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package handlers

import "net/http"

// SecurityHeaders sets several recommended security related
// headers to sane defaults.
//
// It sets:
//  - X-Frame-Options: SAMEORIGIN,
//  - X-XSS-Protection: 1; mode=block,
//  - X-Content-Type-Options: nosniff,
//  - Referrer-Policy: strict-origin-when-cross-origin.
//
// It also optionally sets Content-Security-Policy and
// Strict-Transport-Security to user specified values.
type SecurityHeaders struct {
	http.Handler

	// The value of the Content-Security-Policy
	// header to set.
	//
	// It is recommended to set this to
	//  default-src 'none'; sandbox
	// and use less restrictive policies for each
	// resource as needed.
	ContentSecurityPolicy string

	// The value of the Strict-Transport-Security
	// header to set.
	//
	// It takes a max-age parameter with time in
	// seconds, which should be set to at least six
	// months, like so:
	//  max-age=15768000
	//
	// It also optionally takes two other flags:
	//  - includeSubDomains which applies the policy
	//    to all subdomains, and
	//  - preload which .
	//
	// This header should be used with caution, but
	// it is strongly recommend for all HTTPS sides.
	StrictTransportSecurity string
}

// ServeHTTP implements http.Handler.
func (sh *SecurityHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h["X-Frame-Options"] = []string{"SAMEORIGIN"}
	h["X-XSS-Protection"] = []string{"1; mode=block"}
	h["X-Content-Type-Options"] = []string{"nosniff"}
	h["Referrer-Policy"] = []string{"strict-origin-when-cross-origin"}

	if sh.ContentSecurityPolicy != "" {
		h["Content-Security-Policy"] = []string{sh.ContentSecurityPolicy}
	}

	if sh.StrictTransportSecurity != "" {
		h["Strict-Transport-Security"] = []string{sh.StrictTransportSecurity}
	}

	sh.Handler.ServeHTTP(w, r)
}
