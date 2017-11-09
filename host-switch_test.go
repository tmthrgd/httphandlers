// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"testing"
)

type fakeResponseWriter struct {
	Headers http.Header
	Code    int
}

func (rw *fakeResponseWriter) Header() http.Header {
	return rw.Headers
}

func (*fakeResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (rw *fakeResponseWriter) WriteHeader(code int) {
	rw.Code = code
}

func TestHostSwitchAdd(t *testing.T) {
	var hs HostSwitch
	hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	defer func() {
		if err := recover(); err != nil {
			if err != `handlers: a handle is already registered for host 'example.com'` {
				panic(err)
			}
		} else {
			t.Error("(*HostSwitch).Add did not panic on duplicate")
		}
	}()
	hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
}

func TestHostSwitchNotFound(t *testing.T) {
	calledNotFound := false
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	hs.ServeHTTP(new(fakeResponseWriter), &http.Request{Host: "example.com"})

	if !calledNotFound {
		t.Error("HostSwitch did not call NotFound")
	}
}

func TestHostSwitch(t *testing.T) {
	calledNotFound := false
	hs := &HostSwitch{
		NotFound: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			calledNotFound = true
		}),
	}

	calledExampleCom := false
	hs.Add("example.com", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleCom = true
	}))

	calledExampleOrg := false
	hs.Add("example.org", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledExampleOrg = true
	}))

	hs.ServeHTTP(new(fakeResponseWriter), &http.Request{Host: "example.com"})

	if calledNotFound || !calledExampleCom || calledExampleOrg {
		t.Error("HostSwitch did not call correct handler")
	}
}
