package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdinandanggris/wapi/transport"
)

type stubRT struct {
	attempt int
	code    int
}

func (s *stubRT) RoundTrip(req *http.Request) (*http.Response, error) {
	s.attempt++
	rec := httptest.NewRecorder()
	rec.WriteHeader(s.code)
	return rec.Result(), nil
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantStatus int
		wantAttempt int
	}{
		{"stops on 200", 200, 200, 1},
		{"retries on 429", 429, 429, 3},
		{"retries on 500", 500, 500, 3},
		{"does not retry on 400", 400, 400, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &stubRT{code: tt.statusCode}
			tr := transport.Retry(transport.RetryConfig{MaxAttempts: 3, MinWait: 1})(rt)

			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := tr.RoundTrip(req)

			if rt.attempt != tt.wantAttempt {
				t.Errorf("attempts = %d, want %d", rt.attempt, tt.wantAttempt)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
