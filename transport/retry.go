package transport

import (
	"math/rand"
	"net/http"
	"time"
)

// RetryConfig configures exponential backoff retry behavior.
type RetryConfig struct {
	MaxAttempts int
	MinWait     time.Duration
	MaxWait     time.Duration
}

type retryRT struct {
	next   http.RoundTripper
	config RetryConfig
}

// Retry creates middleware that retries on HTTP 429 and 5xx with exponential backoff + jitter.
// Resets request body via req.GetBody() on each retry.
func Retry(config RetryConfig) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &retryRT{next: next, config: config}
	}
}

// DefaultRetry returns a Retry middleware with 3 max attempts, 1s min wait, 60s max wait.
func DefaultRetry() Middleware {
	return Retry(RetryConfig{
		MaxAttempts: 3,
		MinWait:     time.Second,
		MaxWait:     60 * time.Second,
	})
}

func (r *retryRT) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := r.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	attempts := 1
	for shouldRetry(resp.StatusCode) && attempts < r.config.MaxAttempts {
		resp.Body.Close()

		wait := backoff(attempts, r.config.MinWait, r.config.MaxWait)
		time.Sleep(wait)

		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			req.Body = body
		}

		resp, err = r.next.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		attempts++
	}

	return resp, nil
}

func shouldRetry(code int) bool {
	return code == 429 || code >= 500
}

func backoff(attempt int, min, max time.Duration) time.Duration {
	d := time.Duration(1<<uint(attempt-1)) * min
	if d > max {
		d = max
	}
	if d <= 0 {
		d = time.Millisecond
	}
	jitter := time.Duration(rand.Int63n(int64(d) / 4))
	return d + jitter
}
