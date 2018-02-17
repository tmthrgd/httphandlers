// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bytes"
	"io"
	"time"
)

// Template is an interface that represents both
// text/template and html/template for use with
// ServeTemplate and ServeErrorTemplate.
type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

// ServeTemplate returns a http.Handler that calls
// http.ServeContent with the executed template.
func ServeTemplate(name string, modtime time.Time, tmpl Template, data interface{}) (Handler, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return ServeBytes(name, modtime, buf.Bytes()), nil
}

// ServeErrorTemplate returns a http.Handler that serves
// the executed template with a given HTTP status code.
//
// If mimeType is empty, it will be sniffed from
// content.
func ServeErrorTemplate(code int, tmpl Template, data interface{}, mimeType string) (Handler, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return ServeError(code, buf.Bytes(), mimeType), nil
}
