package cloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	wapi "github.com/ferdinandanggris/wapi"
	"github.com/ferdinandanggris/wapi/types"
)

// CreateTemplate creates a new message template in the specified WABA.
func (c *CloudClient) CreateTemplate(ctx context.Context, wabaID string, tpl *types.Template) (*types.Template, error) {
	path := fmt.Sprintf("%s/message_templates", wabaID)
	var created types.Template
	if err := c.do(ctx, "POST", path, tpl, &created); err != nil {
		return nil, fmt.Errorf("create template: %w", err)
	}
	return &created, nil
}

// EditTemplate updates an existing message template by template ID.
func (c *CloudClient) EditTemplate(ctx context.Context, wabaID, templateID string, tpl *types.Template) error {
	return c.do(ctx, "POST", templateID, tpl, nil)
}

// DeleteTemplate deletes a message template by name.
func (c *CloudClient) DeleteTemplate(ctx context.Context, wabaID, name string) error {
	path := fmt.Sprintf("%s/message_templates?name=%s", wabaID, url.QueryEscape(name))
	return c.doDelete(ctx, path)
}

// GetTemplate retrieves a message template by ID.
func (c *CloudClient) GetTemplate(ctx context.Context, templateID string) (*types.Template, error) {
	var tpl types.Template
	if err := c.do(ctx, "GET", templateID, nil, &tpl); err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	return &tpl, nil
}

// ListTemplates returns all message templates with optional pagination (wapi.WithLimit, wapi.WithOffset).
func (c *CloudClient) ListTemplates(ctx context.Context, wabaID string, opts ...wapi.ListOption) (*types.TemplateList, error) {
	params := &wapi.ListParams{}
	for _, opt := range opts {
		opt(params)
	}

	v := url.Values{}
	if params.Limit > 0 {
		v.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != "" {
		v.Set("after", params.Offset)
	}

	path := fmt.Sprintf("%s/message_templates", wabaID)
	var list types.TemplateList
	if err := c.doGet(ctx, path, v, &list); err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return &list, nil
}
