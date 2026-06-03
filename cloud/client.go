// Package cloud implements the WhatsApp Cloud API client (wapi.Client interface).
package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	wapi "github.com/ferdinandanggris/wapi"
	"github.com/ferdinandanggris/wapi/transport"
)

const defaultBaseURL = "https://graph.facebook.com"
const defaultAPIVersion = "v21.0"

// Option configures a CloudClient.
type Option func(*CloudClient)

// CloudClient implements wapi.Client for the WhatsApp Cloud API.
// Create one with New(), configure with Option values.
type CloudClient struct {
	baseURL     string
	apiVersion  string
	accessToken string
	httpClient  *http.Client
}

// New creates a CloudClient with retry middleware enabled by default.
// Required: WithAccessToken.
func New(opts ...Option) *CloudClient {
	c := &CloudClient{
		baseURL:    defaultBaseURL,
		apiVersion: defaultAPIVersion,
		httpClient: &http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				transport.DefaultRetry(),
			),
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithAccessToken sets the WABA access token (required).
func WithAccessToken(token string) Option {
	return func(c *CloudClient) { c.accessToken = token }
}

// WithAPIVersion overrides the default API version (v21.0).
func WithAPIVersion(v string) Option {
	return func(c *CloudClient) { c.apiVersion = v }
}

// WithBaseURL overrides the default Meta Graph API base URL.
func WithBaseURL(u string) Option {
	return func(c *CloudClient) { c.baseURL = strings.TrimRight(u, "/") }
}

// WithHTTPClient sets a custom http.Client (disables default retry middleware).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *CloudClient) { c.httpClient = hc }
}

// WithRetry configures the retry middleware with the specified max attempts.
func WithRetry(maxAttempts int) Option {
	return func(c *CloudClient) {
		c.httpClient.Transport = transport.Chain(
			http.DefaultTransport,
			transport.Retry(transport.RetryConfig{
				MaxAttempts: maxAttempts,
				MinWait:     time.Second,
				MaxWait:     60 * time.Second,
			}),
		)
	}
}

func (c *CloudClient) apiURL(path string) string {
	if c.apiVersion != "" {
		return fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiVersion, strings.TrimLeft(path, "/"))
	}
	return fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(path, "/"))
}

func (c *CloudClient) do(ctx context.Context, method, path string, body, out interface{}) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.apiURL(path), reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	return c.sendRequest(req, out)
}

func (c *CloudClient) doUpload(ctx context.Context, path, filename, mimeType string, file io.Reader, out interface{}) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if err := w.WriteField("messaging_product", "whatsapp"); err != nil {
		return fmt.Errorf("write field: %w", err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", mimeType)
	fw, err := w.CreatePart(h)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(fw, file); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	if err := w.WriteField("type", mimeType); err != nil {
		return fmt.Errorf("write mime type: %w", err)
	}

	w.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL(path), &buf)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return c.sendRequest(req, out)
}

func (c *CloudClient) doGet(ctx context.Context, path string, params url.Values, out interface{}) error {
	u := c.apiURL(path)
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return fmt.Errorf("create get request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	return c.sendRequest(req, out)
}

func (c *CloudClient) doDelete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.apiURL(path), nil)
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	return c.sendRequest(req, nil)
}

type metaErrorResponse struct {
	Error *metaError `json:"error"`
}

type metaError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Subcode   int    `json:"error_subcode,omitempty"`
	ErrorData *struct {
		Details string `json:"details"`
	} `json:"error_data,omitempty"`
	FBTraceID string `json:"fbtrace_id"`
}

func (c *CloudClient) sendRequest(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.parseError(respBody, resp.StatusCode)
	}

	var merr metaErrorResponse
	_ = json.Unmarshal(respBody, &merr)
	if merr.Error != nil {
		return c.buildError(merr.Error, resp.StatusCode)
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *CloudClient) parseError(respBody []byte, httpCode int) error {
	var merr metaErrorResponse
	if err := json.Unmarshal(respBody, &merr); err != nil || merr.Error == nil {
		return &wapi.Error{
			HTTPCode: httpCode,
			Code:     httpCode,
			Message:  fmt.Sprintf("HTTP %d: %s", httpCode, string(respBody)),
		}
	}
	return c.buildError(merr.Error, httpCode)
}

func (c *CloudClient) buildError(merr *metaError, httpCode int) error {
	details := ""
	if merr.ErrorData != nil {
		details = merr.ErrorData.Details
	}

	errType := wapi.ErrUnknown
	switch merr.Type {
	case "OAuthException":
		errType = wapi.ErrOAuth
	case "GraphMethodException":
		errType = wapi.ErrGraphMethod
	}
	if merr.Code == 130429 {
		errType = wapi.ErrRateLimit
	}

	return &wapi.Error{
		Code:      merr.Code,
		Subcode:   merr.Subcode,
		Message:   merr.Message,
		Type:      errType,
		FBTraceID: merr.FBTraceID,
		HTTPCode:  httpCode,
		Details:   details,
	}
}

// SetAccessToken updates the access token at runtime.
// All subsequent API calls use the new token.
func (c *CloudClient) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *CloudClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

func (c *CloudClient) doPostForm(ctx context.Context, path string, data url.Values, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL(path), strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create post form request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.sendRequest(req, out)
}
