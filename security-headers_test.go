// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(999)
	})

	w := httptest.NewRecorder()
	(&SecurityHeaders{Handler: h}).ServeHTTP(w, r)

	assert.Equal(t, w.Code, 999, "http.Handler not invoked")
	assert.Equal(t, w.HeaderMap, http.Header{
		"X-Frame-Options":        {"SAMEORIGIN"},
		"X-XSS-Protection":       {"1; mode=block"},
		"X-Content-Type-Options": {"nosniff"},
		"Referrer-Policy":        {"strict-origin-when-cross-origin"},
	})

	w = httptest.NewRecorder()
	(&SecurityHeaders{
		Handler: h,

		ContentSecurityPolicy:   "test1",
		StrictTransportSecurity: "test2",
		ExpectCT:                "test3",
	}).ServeHTTP(w, r)

	assert.Equal(t, w.Code, 999, "http.Handler not invoked")
	assert.Equal(t, w.HeaderMap, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-XSS-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"test1"},
		"Strict-Transport-Security": {"test2"},
		"Expect-CT":                 {"test3"},
	})

	w = httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"X-Frame-Options":           {"fail"},
		"X-XSS-Protection":          {"fail"},
		"X-Content-Type-Options":    {"fail"},
		"Referrer-Policy":           {"fail"},
		"Content-Security-Policy":   {"fail"},
		"Strict-Transport-Security": {"fail"},
		"Expect-CT":                 {"fail"},
	}

	(&SecurityHeaders{
		Handler: h,

		ContentSecurityPolicy:   "test1",
		StrictTransportSecurity: "test2",
		ExpectCT:                "test3",
	}).ServeHTTP(w, r)

	assert.Equal(t, w.Code, 999, "http.Handler not invoked")
	assert.Equal(t, w.HeaderMap, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-XSS-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"test1"},
		"Strict-Transport-Security": {"test2"},
		"Expect-CT":                 {"test3"},
	})

	w = httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"X-Frame-Options":           {"fail"},
		"X-XSS-Protection":          {"fail"},
		"X-Content-Type-Options":    {"fail"},
		"Referrer-Policy":           {"fail"},
		"Content-Security-Policy":   {"leave"},
		"Strict-Transport-Security": {"leave"},
		"Expect-CT":                 {"leave"},
	}

	(&SecurityHeaders{Handler: h}).ServeHTTP(w, r)

	assert.Equal(t, w.Code, 999, "http.Handler not invoked")
	assert.Equal(t, w.HeaderMap, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-XSS-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"leave"},
		"Strict-Transport-Security": {"leave"},
		"Expect-CT":                 {"leave"},
	})
}
