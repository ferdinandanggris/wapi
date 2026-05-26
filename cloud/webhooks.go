package cloud

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ferdinandanggris/wapi/types"
)

// SubscribeToWebhooks subscribes the app to webhook events for the WABA.
func (c *CloudClient) SubscribeToWebhooks(ctx context.Context, wabaID string) (*types.SubscribedApp, error) {
	path := fmt.Sprintf("%s/subscribed_apps", wabaID)
	var app types.SubscribedApp
	if err := c.do(ctx, "POST", path, nil, &app); err != nil {
		return nil, fmt.Errorf("subscribe webhooks: %w", err)
	}
	return &app, nil
}

// UnsubscribeFromWebhooks unsubscribes the app from all webhook events.
func (c *CloudClient) UnsubscribeFromWebhooks(ctx context.Context, wabaID string) error {
	path := fmt.Sprintf("%s/subscribed_apps", wabaID)
	return c.doDelete(ctx, path)
}

// GetWebhookSubscription returns the current webhook subscription status.
func (c *CloudClient) GetWebhookSubscription(ctx context.Context, wabaID string) (*types.SubscribedApp, error) {
	path := fmt.Sprintf("%s/subscribed_apps", wabaID)
	var app types.SubscribedApp
	if err := c.do(ctx, "GET", path, nil, &app); err != nil {
		return nil, fmt.Errorf("get webhook subscription: %w", err)
	}
	return &app, nil
}

// SetWebhookFields updates which webhook event fields are subscribed.
func (c *CloudClient) SetWebhookFields(ctx context.Context, appID string, fields ...string) error {
	path := fmt.Sprintf("%s/subscriptions", appID)
	data := url.Values{
		"object": {"whatsapp_business_account"},
		"fields": {strings.Join(fields, ",")},
	}
	return c.doPostForm(ctx, path, data, nil)
}

// SetWebhookCallback sets the callback URL, verify token, and subscribed fields.
func (c *CloudClient) SetWebhookCallback(ctx context.Context, appID, callbackURL, verifyToken string, fields ...string) error {
	path := fmt.Sprintf("%s/subscriptions", appID)
	data := url.Values{
		"object":       {"whatsapp_business_account"},
		"callback_url": {callbackURL},
		"verify_token": {verifyToken},
		"fields":       {strings.Join(fields, ",")},
	}
	return c.doPostForm(ctx, path, data, nil)
}
