package cloud

import (
	"context"
	"fmt"

	"github.com/ferdinandanggris/wapi/types"
)

// GetBusinessProfile returns the WhatsApp Business profile for a phone number.
func (c *CloudClient) GetBusinessProfile(ctx context.Context, phoneNumberID string) (*types.BusinessProfile, error) {
	path := fmt.Sprintf("%s/whatsapp_business_profile", phoneNumberID)
	var result struct {
		Data []*types.BusinessProfile `json:"data"`
	}
	if err := c.do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get business profile: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("get business profile: no profile found")
	}
	return result.Data[0], nil
}

// UpdateBusinessProfile updates the WhatsApp Business profile for a phone number.
func (c *CloudClient) UpdateBusinessProfile(ctx context.Context, phoneNumberID string, profile *types.BusinessProfile) error {
	profile.MessagingProduct = "whatsapp"
	path := fmt.Sprintf("%s/whatsapp_business_profile", phoneNumberID)
	return c.do(ctx, "POST", path, profile, nil)
}
