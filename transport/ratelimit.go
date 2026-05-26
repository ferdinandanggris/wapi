package transport

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// TokenBucket implements a generic token bucket rate limiter.
type TokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64
	last     time.Time
}

// NewTokenBucket creates a token bucket that refills at rate (tokens/sec) up to burst capacity.
func NewTokenBucket(rate float64, burst int) *TokenBucket {
	return &TokenBucket{
		tokens: float64(burst),
		max:    float64(burst),
		rate:   rate,
		last:   time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.tokens = tb.tokens + elapsed*tb.rate
	if tb.tokens > tb.max {
		tb.tokens = tb.max
	}
	tb.last = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

type rateLimitRT struct {
	next   http.RoundTripper
	bucket *TokenBucket
}

// RateLimit creates middleware that blocks until a token is available from the bucket.
// Respects context cancellation.
func RateLimit(rate float64, burst int) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &rateLimitRT{
			next:   next,
			bucket: NewTokenBucket(rate, burst),
		}
	}
}

func (r *rateLimitRT) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		if r.bucket.Allow() {
			break
		}
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
	return r.next.RoundTrip(req)
}

type contextKey string

const bucketKey contextKey = "rate_limit_bucket"

// WithBucket associates a shared token bucket with a context for cross-client rate limiting.
func WithBucket(ctx context.Context, bucket *TokenBucket) context.Context {
	return context.WithValue(ctx, bucketKey, bucket)
}

// FromContext retrieves a token bucket previously stored by WithBucket.
func FromContext(ctx context.Context) *TokenBucket {
	b, _ := ctx.Value(bucketKey).(*TokenBucket)
	return b
}
