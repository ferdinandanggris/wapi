// Package wapi defines the Client interface and shared types for the WhatsApp Cloud API.
package wapi

import (
	"context"
	"io"

	"github.com/ferdinandanggris/wapi/types"
)

// Client is the WhatsApp Cloud API interface.
// Implementations: cloud.CloudClient.
type Client interface {
	SendMessage(ctx context.Context, phoneNumberID string, msg *types.Message) (*types.SendResponse, error)
	MarkAsRead(ctx context.Context, phoneNumberID string, messageID string) error

	UploadMedia(ctx context.Context, phoneNumberID string, filename string, file io.Reader, mimeType string) (*types.MediaUploadResponse, error)
	GetMediaURL(ctx context.Context, mediaID string) (*types.MediaInfo, error)
	DownloadMedia(ctx context.Context, mediaID string) (io.ReadCloser, error)
	DeleteMedia(ctx context.Context, mediaID string) error

	CreateTemplate(ctx context.Context, wabaID string, tpl *types.Template) (*types.Template, error)
	EditTemplate(ctx context.Context, wabaID, templateID string, tpl *types.Template) error
	DeleteTemplate(ctx context.Context, wabaID string, name string) error
	GetTemplate(ctx context.Context, templateID string) (*types.Template, error)
	ListTemplates(ctx context.Context, wabaID string, opts ...ListOption) (*types.TemplateList, error)

	RegisterPhone(ctx context.Context, phoneNumberID, pin string) error
	DeregisterPhone(ctx context.Context, phoneNumberID string) error
	GetPhoneNumber(ctx context.Context, phoneNumberID string) (*types.PhoneNumber, error)
	ListPhoneNumbers(ctx context.Context, wabaID string) ([]*types.PhoneNumber, error)
	SetTwoStepPIN(ctx context.Context, phoneNumberID, pin string) error

	GetBusinessProfile(ctx context.Context, phoneNumberID string) (*types.BusinessProfile, error)
	UpdateBusinessProfile(ctx context.Context, phoneNumberID string, profile *types.BusinessProfile) error

	SubscribeToWebhooks(ctx context.Context, wabaID string) (*types.SubscribedApp, error)
	UnsubscribeFromWebhooks(ctx context.Context, wabaID string) error
	GetWebhookSubscription(ctx context.Context, wabaID string) (*types.SubscribedApp, error)
	SetWebhookFields(ctx context.Context, appID string, fields ...string) error
	SetWebhookCallback(ctx context.Context, appID, callbackURL, verifyToken string, fields ...string) error

	Close() error
}

// ListOption configures pagination for list operations.
type ListOption func(*ListParams)

// ListParams holds pagination parameters.
type ListParams struct {
	Limit  int
	Offset string
}

// WithLimit sets the maximum number of items to return.
func WithLimit(n int) ListOption {
	return func(p *ListParams) { p.Limit = n }
}

// WithOffset sets the pagination cursor for the next page.
func WithOffset(token string) ListOption {
	return func(p *ListParams) { p.Offset = token }
}
