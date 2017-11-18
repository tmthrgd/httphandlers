// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	ht "html/template"
	"testing"
	tt "text/template"

	"github.com/stretchr/testify/assert"
)

func TestTemplateInterface(t *testing.T) {
	assert.Implements(t, (*Template)(nil), (*ht.Template)(nil))
	assert.Implements(t, (*Template)(nil), (*tt.Template)(nil))
}
