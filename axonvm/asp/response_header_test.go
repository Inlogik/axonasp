/*
 * AxonASP Server
 * Copyright (C) 2026 G3pix Ltda. All rights reserved.
 */
package asp

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestResponseAddHeaderPreservesRepeatedSetCookieValues(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := NewResponse(recorder)
	response.AddHeader("Set-Cookie", "fixation=one; Path=/; HttpOnly")
	response.AddHeader("Set-Cookie", "auth=two; Path=/; HttpOnly")
	response.Flush()

	want := []string{
		"fixation=one; Path=/; HttpOnly",
		"auth=two; Path=/; HttpOnly",
	}
	if got := recorder.Header().Values("Set-Cookie"); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected Set-Cookie headers: got %#v, want %#v", got, want)
	}
}

func TestResponseAddHeaderReplacesRepeatedCookieName(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := NewResponse(recorder)
	response.AddHeader("Set-Cookie", "user=first; Path=/; HttpOnly")
	response.AddHeader("Set-Cookie", "other=value; Path=/; HttpOnly")
	response.AddHeader("Set-Cookie", "USER=latest; Path=/; HttpOnly")
	response.Flush()

	want := []string{
		"other=value; Path=/; HttpOnly",
		"USER=latest; Path=/; HttpOnly",
	}
	if got := recorder.Header().Values("Set-Cookie"); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected Set-Cookie headers: got %#v, want %#v", got, want)
	}
}

func TestResponseAddHeaderNormalizesDoubleSlashCookiePath(t *testing.T) {
	response := NewResponse(httptest.NewRecorder())
	response.AddHeader("Set-Cookie", "session=value; path=//; HttpOnly")
	response.Flush()

	cookies := response.Output.(*httptest.ResponseRecorder).Result().Header.Values("Set-Cookie")
	if len(cookies) != 1 || !strings.Contains(cookies[0], "path=/;") {
		t.Fatalf("expected normalized cookie path, got %#v", cookies)
	}
}
