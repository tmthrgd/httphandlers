// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

// Package handlers is a collection of small utility
// http.Handler's for Golang.
package handlers

import "net/http"

// Handler is an alias to http.Handler for godoc.
type Handler = http.Handler

// Middleware represents a function that wraps an
// http.Handler.
type Middleware = func(http.Handler) http.Handler
