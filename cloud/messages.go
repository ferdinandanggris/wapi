package cloud

import (
	"context"
	"fmt"

	"github.com/ferdinandanggris/wapi/types"
)

// SendMessage sends a WhatsApp message. Use types.New*Message builders to create the message.
func (c *CloudClient) SendMessage(ctx context.Context, phoneNumberID string, msg *types.Message) (*types.SendResponse, error) {
	msg.MessagingProduct = "whatsapp"

	if msg.To == "" {
		return nil, fmt.Errorf("send message: recipient 'to' is required")
	}

	path := fmt.Sprintf("%s/messages", phoneNumberID)
	var resp types.SendResponse
	if err := c.do(ctx, "POST", path, msg, &resp); err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	return &resp, nil
}

// MarkAsRead marks an incoming message as read. messageID must be an incoming message ID.
func (c *CloudClient) MarkAsRead(ctx context.Context, phoneNumberID string, messageID string) error {
	msg := types.NewMarkAsRead(messageID)

	path := fmt.Sprintf("%s/messages", phoneNumberID)
	return c.do(ctx, "POST", path, msg, nil)
}
