package types

// PhoneNumber represents a WhatsApp Business phone number with its quality and limits.
type PhoneNumber struct {
	ID                     string `json:"id"`
	DisplayPhoneNumber     string `json:"display_phone_number"`
	VerifiedName           string `json:"verified_name"`
	QualityRating          string `json:"quality_rating"`
	MessagingLimit         string `json:"messaging_limit"`
	PlatformType           string `json:"platform_type"`
	CodeVerificationStatus string `json:"code_verification_status"`
	PinEnabled             bool   `json:"pin_enabled"`
	Throughput             *Throughput `json:"throughput,omitempty"`
}

type Throughput struct {
	Level string `json:"level"`
}

// BusinessProfile holds the WhatsApp Business profile information.
type BusinessProfile struct {
	MessagingProduct string   `json:"messaging_product"`
	Address          string   `json:"address,omitempty"`
	Description      string   `json:"description,omitempty"`
	Vertical         string   `json:"vertical,omitempty"`
	Email            string   `json:"email,omitempty"`
	Websites         []string `json:"websites,omitempty"`
	ProfilePictureURL string  `json:"profile_picture_url,omitempty"`
}

type RegisterRequest struct {
	MessagingProduct string `json:"messaging_product"`
	Pin              string `json:"pin"`
}
