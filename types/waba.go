package types

// WhatsAppBusinessAccount represents a WhatsApp Business Account (WABA) returned by owned_whatsapp_business_accounts.
type WhatsAppBusinessAccount struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Currency                 string `json:"currency,omitempty"`
	TimezoneID               string `json:"timezone_id,omitempty"`
	MessageTemplateNamespace string `json:"message_template_namespace,omitempty"`
}

// WABAList is the paginated response from listing owned WABAs.
type WABAList struct {
	Data   []*WhatsAppBusinessAccount `json:"data"`
	Paging *Paging                    `json:"paging,omitempty"`
}
