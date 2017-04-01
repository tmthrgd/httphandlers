// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

type errorResponseWriter struct {
	http.ResponseWriter
	request *http.Request

	errors         map[int]*StaticError
	disablePadding bool

	didWrite  bool
	skipWrite bool
}

func (w *errorResponseWriter) WriteHeader(code int) {
	page, ok := w.errors[code]
	if !ok || w.didWrite || w.skipWrite {
		w.ResponseWriter.WriteHeader(code)
		return
	}

	w.skipWrite = true

	var padding string
	if code >= http.StatusBadRequest && !w.disablePadding {
		if ua := w.request.Header.Get("User-Agent"); ua != "" {
			if msie := strings.Index(ua, "MSIE "); msie != -1 && msie+7 < len(ua) && !strings.Contains(ua, "Opera") {
				padding = `
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->
<!-- a padding to disable MSIE and Chrome friendly error page -->`
			}
		}
	}

	h := w.Header()
	delete(h, "Cache-Control")
	delete(h, "Etag")
	delete(h, "Last-Modified")
	delete(h, "Content-Encoding")
	h["Content-Type"] = []string{"text/html; charset=utf-8"}

	for k, v := range page.Headers {
		h[k] = v
	}

	h["Content-Length"] = []string{strconv.FormatInt(int64(len(page.Body)+len(padding)), 10)}

	w.ResponseWriter.WriteHeader(code)

	if w.request.Method == http.MethodHead {
		return
	}

	_, err := w.ResponseWriter.Write(page.Body)

	if err == nil && padding != "" {
		_, err = io.WriteString(w.ResponseWriter, padding)
	}

	if err != nil {
		server := w.request.Context().Value(http.ServerContextKey).(*http.Server)
		if server.ErrorLog != nil {
			server.ErrorLog.Println(err)
		}
	}
}

func (w *errorResponseWriter) Write(p []byte) (int, error) {
	if w.skipWrite {
		return len(p), nil
	}

	w.didWrite = true
	return w.ResponseWriter.Write(p)
}

// StaticErrors intercepts calls to
// http.ResponseWriter.WriteHeader and replaces
// the response with a statically rendered error
// page.
type StaticErrors struct {
	http.Handler

	// Errors is a map of HTTP status code
	// (for example http.StatusNotFound) to
	// a statically rendered error page.
	Errors map[int]*StaticError

	// DisablePadding disables the MSIE/Chrome
	// friendly error padding.
	DisablePadding bool
}

// ServeHTTP implements http.Handler.
func (e *StaticErrors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Handler.ServeHTTP(&errorResponseWriter{
		ResponseWriter: w,
		request:        r,

		errors:         e.Errors,
		disablePadding: e.DisablePadding,
	}, r)
}

// StaticError represents a statically rendered
// error page.
type StaticError struct {
	// Body is the rendered page.
	Body []byte

	// Headers is a map of headers to set
	// on the response for this error page.
	Headers http.Header
}

// DefaultErrorMessages is a list of suggested error
// codes and messages to render.
var DefaultErrorMessages = map[int]struct{ Name, Message string }{
	http.StatusBadRequest: {
		"Bad Request",
		"Your user agent sent a request that this server could not understand.",
	},
	http.StatusForbidden: {
		"Forbidden",
		"You do not have permission to access this resource.",
	},
	http.StatusNotFound: {
		"File Not Found",
		"The link you followed may be broken, or the page may have been removed.",
	},
	http.StatusMethodNotAllowed: {
		"Method Not Allowed",
		"The specified HTTP method is not allowed for the requested resource.",
	},
	http.StatusInternalServerError: {
		"Internal Server Error",
		"An internal server error has occurred.",
	},
}
