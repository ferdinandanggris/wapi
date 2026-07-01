package cloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	wapi "github.com/ferdinandanggris/wapi"
	"github.com/ferdinandanggris/wapi/types"
)

// ListWhatsAppBusinessAccounts returns all WABAs owned by the given business ID, with optional pagination.
func (c *CloudClient) ListWhatsAppBusinessAccounts(ctx context.Context, businessID string, opts ...wapi.ListOption) (*types.WABAList, error) {
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

	path := fmt.Sprintf("%s/owned_whatsapp_business_accounts", businessID)
	var list types.WABAList
	if err := c.doGet(ctx, path, v, &list); err != nil {
		return nil, fmt.Errorf("list waba accounts: %w", err)
	}
	return &list, nil
}
