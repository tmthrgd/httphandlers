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

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Frame-Options":        {"SAMEORIGIN"},
		"X-Xss-Protection":       {"1; mode=block"},
		"X-Content-Type-Options": {"nosniff"},
		"Referrer-Policy":        {"strict-origin-when-cross-origin"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	(&SecurityHeaders{
		Handler: h,

		ContentSecurityPolicy:   "test1",
		StrictTransportSecurity: "test2",
		ExpectCT:                "test3",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-Xss-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"test1"},
		"Strict-Transport-Security": {"test2"},
		"Expect-Ct":                 {"test3"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"X-Frame-Options":           {"fail"},
		"X-Xss-Protection":          {"fail"},
		"X-Content-Type-Options":    {"fail"},
		"Referrer-Policy":           {"fail"},
		"Content-Security-Policy":   {"fail"},
		"Strict-Transport-Security": {"fail"},
		"Expect-Ct":                 {"fail"},
	}

	(&SecurityHeaders{
		Handler: h,

		ContentSecurityPolicy:   "test1",
		StrictTransportSecurity: "test2",
		ExpectCT:                "test3",
	}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-Xss-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"test1"},
		"Strict-Transport-Security": {"test2"},
		"Expect-Ct":                 {"test3"},
	}, w.Result().Header)

	w = httptest.NewRecorder()
	w.HeaderMap = http.Header{
		"X-Frame-Options":           {"fail"},
		"X-Xss-Protection":          {"fail"},
		"X-Content-Type-Options":    {"fail"},
		"Referrer-Policy":           {"fail"},
		"Content-Security-Policy":   {"leave"},
		"Strict-Transport-Security": {"leave"},
		"Expect-Ct":                 {"leave"},
	}

	(&SecurityHeaders{Handler: h}).ServeHTTP(w, r)

	assert.Equal(t, 999, w.Code, "http.Handler not invoked")
	assert.Equal(t, http.Header{
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-Xss-Protection":          {"1; mode=block"},
		"X-Content-Type-Options":    {"nosniff"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Content-Security-Policy":   {"leave"},
		"Strict-Transport-Security": {"leave"},
		"Expect-Ct":                 {"leave"},
	}, w.Result().Header)
}
