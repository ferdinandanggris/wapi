package wapi

import "fmt"

// ErrorType classifies a wapi error.
type ErrorType string

const (
	ErrOAuth       ErrorType = "OAuthException"       // Token expired, invalid permissions
	ErrGraphMethod ErrorType = "GraphMethodException"  // Invalid request, missing fields
	ErrRateLimit   ErrorType = "RateLimit"             // Too many requests (code 130429)
	ErrServer      ErrorType = "Server"                // Meta server error (HTTP 5xx)
	ErrUnknown     ErrorType = "Unknown"               // Unclassified error
)

// Error represents a WhatsApp Cloud API error.
type Error struct {
	Code      int
	Subcode   int
	Message   string
	Type      ErrorType
	FBTraceID string
	HTTPCode  int
	Details   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("wapi: [%d] %s (trace: %s)", e.Code, e.Message, e.FBTraceID)
}

func (e *Error) Unwrap() error { return nil }

// IsRetryable returns true if the error is retryable (HTTP 5xx or rate limit code).
func IsRetryable(err error) bool {
	var e *Error
	if !as(err, &e) {
		return false
	}
	if e.HTTPCode >= 500 {
		return true
	}
	return e.Code == 130429 || e.Code == 131056
}

// IsRateLimit returns true if the error is a rate limit (code 130429).
func IsRateLimit(err error) bool {
	var e *Error
	if !as(err, &e) {
		return false
	}
	return e.Code == 130429
}

func as(err error, target interface{}) bool {
	switch t := target.(type) {
	case **Error:
		e, ok := err.(*Error)
		if !ok {
			return false
		}
		*t = e
		return true
	}
	return false
}
