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
//  - X-Content-Type-Options: nosniff, and
//  - Referrer-Policy: strict-origin-when-cross-origin.
//
// It also optionally sets the Content-Security-Policy,
// Strict-Transport-Security and Expect-CT to user
// specified values.
type SecurityHeaders struct {
	Handler http.Handler

	// The value of the Content-Security-Policy
	// header to set.
	//
	// It is recommended to set this to
	//  default-src 'none'; sandbox
	// and use less restrictive policies for each
	// resource as needed.
	//
	// This header may require caution to use safely,
	// but it is strongly recommend for all sites.
	//
	// See the article
	//  'Content Security Policy - An Introduction'
	//   https://scotthelme.co.uk/content-security-policy-an-introduction/
	// for more information.
	ContentSecurityPolicy string

	// The value of the Strict-Transport-Security
	// header to set.
	//
	// It takes a max-age directive, with time in
	// seconds, which indicate how long browsers
	// should cache the policy. It should be set to
	// at least six months, like so:
	//  max-age=15768000
	//
	// It also optionally takes two other directives:
	//  - includeSubDomains, which applies the policy
	//    to all subdomains, and
	//  - preload, which signals to browser
	//    manufacturers that this policy may be
	//    preloaded into the browser to prevent
	//    them from ever connecting to the site
	//    without TLS. Visit https://hstspreload.org/
	//    to request preloading.
	//
	// This header may require caution to use safely,
	// but it is strongly recommend for all HTTPS
	// only sites.
	//
	// See the article
	//  'HSTS - The missing link in Transport Layer Security'
	//   https://scotthelme.co.uk/hsts-the-missing-link-in-tls/
	// for more information.
	StrictTransportSecurity string

	// The value of the Expect-CT header to set.
	//
	// It takes a max-age directive, with time in
	// seconds, which indicate how long browsers
	// should cache the policy.
	//
	// It also optionally takes two other directives:
	//  - enforce, which indicates that browsers
	//    should enforce the policy or treat it as
	//    a report-only policy, and
	//  - report-uri, which specifies a URI that
	//    a browser should send a report to, if it
	//    doesn't receive valid CT information.
	//
	// See the article
	//  'A new security header: Expect-CT'
	//   https://scotthelme.co.uk/a-new-security-header-expect-ct/
	// for more information.
	ExpectCT string

	// The value of the Report-To header to set.
	//
	// It should be a json object containing a group,
	// max_age, endpoints and include_subdomains fields.
	//
	// See the article
	//  'Introducing the Reporting API, Network Error Logging and other major upgrades to Report URI'
	//   https://scotthelme.co.uk/introducing-the-reporting-api-nel-other-major-changes-to-report-uri/
	// for more information.
	ReportTo string

	// The value of the NEL header to set.
	//
	// It should be a json object containing a report_to,
	// max_age and include_subdomains fields.
	//
	// See the article
	//  'Network Error Logging: Deep Dive'
	//   https://scotthelme.co.uk/network-error-logging-deep-dive/
	// for more information.
	NEL string

	// The value of the Feature-Policy header to set.
	//
	// This header allows a site to control what browser
	// features are allowed to be used.
	//
	// See the article
	//  'A new security header: Feature Policy'
	//   https://scotthelme.co.uk/a-new-security-header-feature-policy/
	// for more information.
	FeaturePolicy string
}

// SecurityHeadersWrap returns a Middleware that produces a
// *SecurityHeaders handler.
func SecurityHeadersWrap(c *SecurityHeaders) Middleware {
	return func(h http.Handler) http.Handler {
		sh := new(SecurityHeaders)

		if c != nil {
			*sh = *c
		}

		sh.Handler = h
		return sh
	}
}

// ServeHTTP implements http.Handler.
func (sh *SecurityHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("X-Frame-Options", "SAMEORIGIN")
	h.Set("X-Xss-Protection", "1; mode=block")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	if sh.ContentSecurityPolicy != "" {
		h.Set("Content-Security-Policy", sh.ContentSecurityPolicy)
	}

	if sh.StrictTransportSecurity != "" {
		h.Set("Strict-Transport-Security", sh.StrictTransportSecurity)
	}

	if sh.ExpectCT != "" {
		h.Set("Expect-Ct", sh.ExpectCT)
	}

	if sh.ReportTo != "" {
		h.Set("Report-To", sh.ReportTo)
	}

	if sh.NEL != "" {
		h.Set("Nel", sh.NEL)
	}

	if sh.FeaturePolicy != "" {
		h.Set("Feature-Policy", sh.FeaturePolicy)
	}

	sh.Handler.ServeHTTP(w, r)
}
