package wapi_test

import (
	"fmt"
	"testing"

	wapi "github.com/ferdinandanggris/wapi"
)

func TestError_Error(t *testing.T) {
	e := &wapi.Error{
		Code:      131030,
		Message:   "Recipient not valid WhatsApp user",
		FBTraceID: "ABC123",
	}

	want := "wapi: [131030] Recipient not valid WhatsApp user (trace: ABC123)"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{&wapi.Error{Code: 130429, HTTPCode: 429}, true},
		{&wapi.Error{Code: 131056, HTTPCode: 429}, true},
		{&wapi.Error{Code: 500, HTTPCode: 500}, true},
		{&wapi.Error{Code: 503, HTTPCode: 503}, true},
		{&wapi.Error{Code: 131030, HTTPCode: 400}, false},
		{&wapi.Error{Code: 132012, HTTPCode: 400}, false},
		{&wapi.Error{Code: 190, HTTPCode: 401}, false},
		{fmt.Errorf("random error"), false},
		{nil, false},
	}

	for _, tt := range tests {
		got := wapi.IsRetryable(tt.err)
		if got != tt.want {
			t.Errorf("IsRetryable(%+v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestIsRateLimit(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{&wapi.Error{Code: 130429}, true},
		{&wapi.Error{Code: 131056}, false},
		{&wapi.Error{Code: 500}, false},
		{nil, false},
	}

	for _, tt := range tests {
		got := wapi.IsRateLimit(tt.err)
		if got != tt.want {
			t.Errorf("IsRateLimit(%+v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}
