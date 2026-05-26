package cloud

import (
	"context"
	"fmt"

	"github.com/ferdinandanggris/wapi/types"
)

// RegisterPhone registers a phone number with a 6-digit PIN.
func (c *CloudClient) RegisterPhone(ctx context.Context, phoneNumberID, pin string) error {
	body := map[string]string{
		"messaging_product": "whatsapp",
		"pin":               pin,
	}
	path := fmt.Sprintf("%s/register", phoneNumberID)
	return c.do(ctx, "POST", path, body, nil)
}

// DeregisterPhone deregisters a phone number from the WABA.
func (c *CloudClient) DeregisterPhone(ctx context.Context, phoneNumberID string) error {
	path := fmt.Sprintf("%s/deregister", phoneNumberID)
	return c.do(ctx, "POST", path, map[string]string{"messaging_product": "whatsapp"}, nil)
}

// GetPhoneNumber returns details for a specific phone number.
func (c *CloudClient) GetPhoneNumber(ctx context.Context, phoneNumberID string) (*types.PhoneNumber, error) {
	var pn types.PhoneNumber
	if err := c.do(ctx, "GET", phoneNumberID, nil, &pn); err != nil {
		return nil, fmt.Errorf("get phone number: %w", err)
	}
	return &pn, nil
}

// ListPhoneNumbers returns all phone numbers associated with a WABA.
func (c *CloudClient) ListPhoneNumbers(ctx context.Context, wabaID string) ([]*types.PhoneNumber, error) {
	path := fmt.Sprintf("%s/phone_numbers", wabaID)
	var result struct {
		Data []*types.PhoneNumber `json:"data"`
	}
	if err := c.do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("list phone numbers: %w", err)
	}
	return result.Data, nil
}

// SetTwoStepPIN enables or changes the 6-digit PIN for two-step verification.
func (c *CloudClient) SetTwoStepPIN(ctx context.Context, phoneNumberID, pin string) error {
	body := map[string]string{"pin": pin}
	return c.do(ctx, "POST", phoneNumberID, body, nil)
}
