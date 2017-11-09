// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// ServeTemplate returns a http.Handler that calls
// http.ServeContent with the executed template.
func ServeTemplate(name string, modtime time.Time, tmpl *template.Template, data interface{}) (http.Handler, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return &serveBytes{name, modtime, buf.Bytes()}, nil
}
