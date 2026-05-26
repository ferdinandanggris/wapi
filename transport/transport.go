// Package transport provides HTTP middleware: retry, rate limiting, and composable chain.
package transport

import "net/http"

// Middleware wraps an http.RoundTripper with cross-cutting behavior (retry, rate-limit, etc).
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain composes middleware into a single RoundTripper.
// Middleware are applied left-to-right (outermost first).
func Chain(base http.RoundTripper, mws ...Middleware) http.RoundTripper {
	rt := base
	for i := len(mws) - 1; i >= 0; i-- {
		rt = mws[i](rt)
	}
	return rt
}
