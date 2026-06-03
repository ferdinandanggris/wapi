package types

// WebhookPayload is the root object received from Meta webhooks.
type WebhookPayload struct {
	Object string          `json:"object"`
	Entry  []*WebhookEntry `json:"entry"`
}

// WebhookEntry contains all changes for a specific WABA.
type WebhookEntry struct {
	ID      string           `json:"id"`
	Changes []*WebhookChange `json:"changes"`
}

// WebhookChange holds the actual message or status update payload.
type WebhookChange struct {
	Value *WebhookValue `json:"value"`
	Field string        `json:"field"`
}

// WebhookValue contains the actual payload: messages, statuses, contacts, and metadata.
type WebhookValue struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         *Metadata        `json:"metadata"`
	Contacts         []*WaContact     `json:"contacts,omitempty"`
	Messages         []*IncomingMsg   `json:"messages,omitempty"`
	Statuses         []*StatusUpdate  `json:"statuses,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WaContact struct {
	Profile *Profile `json:"profile"`
	WaID    string   `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

// IncomingMsg represents a message received via webhook (text, image, interactive reply, etc).
type IncomingMsg struct {
	From        string             `json:"from"`
	ID          string             `json:"id"`
	Timestamp   string             `json:"timestamp"`
	Type        string             `json:"type"`
	Text        *IncomingText      `json:"text,omitempty"`
	Image       *IncomingMedia     `json:"image,omitempty"`
	Video       *IncomingMedia     `json:"video,omitempty"`
	Audio       *IncomingMedia     `json:"audio,omitempty"`
	Document    *IncomingDocument  `json:"document,omitempty"`
	Location    *IncomingLocation  `json:"location,omitempty"`
	Contacts    []*Contact         `json:"contacts,omitempty"`
	Interactive *IncomingInteractive `json:"interactive,omitempty"`
	Button      *IncomingButton    `json:"button,omitempty"`
	Context     *IncomingContext   `json:"context,omitempty"`
	Referral    *IncomingReferral  `json:"referral,omitempty"`
	Order       *IncomingOrder     `json:"order,omitempty"`
	Reaction    *IncomingReaction  `json:"reaction,omitempty"`
}

type IncomingText struct {
	Body string `json:"body"`
}

type IncomingMedia struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type IncomingDocument struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename"`
}

type IncomingLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type IncomingInteractive struct {
	Type            string               `json:"type"`
	InButtonReply   *IncomingButtonReply  `json:"button_reply,omitempty"`
	InListReply     *IncomingListReply    `json:"list_reply,omitempty"`
}

type IncomingButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type IncomingListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type IncomingButton struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

type IncomingContext struct {
	From      string `json:"from,omitempty"`
	ID        string `json:"id,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
}

type IncomingReferral struct {
	SourceURL     string `json:"source_url,omitempty"`
	SourceType    string `json:"source_type,omitempty"`
	SourceID      string `json:"source_id,omitempty"`
	Headline      string `json:"headline,omitempty"`
	Body          string `json:"body,omitempty"`
	MediaType     string `json:"media_type,omitempty"`
	ImageURL      string `json:"image_url,omitempty"`
	VideoURL      string `json:"video_url,omitempty"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	CTAPayload    string `json:"cta_payload,omitempty"`
}

type IncomingOrder struct {
	CatalogID string          `json:"catalog_id"`
	ProductItems []*OrderItem `json:"product_items"`
	Text      string          `json:"text,omitempty"`
}

type OrderItem struct {
	ProductRetailerID string `json:"product_retailer_id"`
	Quantity          string `json:"quantity,omitempty"`
	ItemPrice         string `json:"item_price,omitempty"`
	Currency          string `json:"currency,omitempty"`
}

// StatusUpdate is a delivery or read receipt for a sent message.
type StatusUpdate struct {
	ID            string              `json:"id"`
	Status        string              `json:"status"`
	Timestamp     string              `json:"timestamp"`
	RecipientID   string              `json:"recipient_id"`
	Conversation  *StatusConversation `json:"conversation,omitempty"`
	Pricing       *StatusPricing      `json:"pricing,omitempty"`
	Errors        []*StatusError      `json:"errors,omitempty"`
}

type StatusConversation struct {
	ID     string              `json:"id"`
	Origin *ConversationOrigin `json:"origin"`
}

type ConversationOrigin struct {
	Type string `json:"type"`
}

type StatusPricing struct {
	Billable      bool   `json:"billable"`
	PricingModel  string `json:"pricing_model"`
	Category      string `json:"category"`
}

type StatusError struct {
	Code      int              `json:"code"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	ErrorData *StatusErrorData `json:"error_data,omitempty"`
}

type StatusErrorData struct {
	Details string `json:"details"`
}

// SubscribedApp represents the app subscribed to webhook events.
type SubscribedApp struct {
	Name                 string   `json:"name,omitempty"`
	ID                   string   `json:"id,omitempty"`
	OverrideCallbackURI  string   `json:"override_callback_uri,omitempty"`
}

// Subscription represents a webhook subscription configuration.
type Subscription struct {
	Object             string              `json:"object"`
	CallbackURL        string              `json:"callback_url,omitempty"`
	Fields             []*SubscriptionField `json:"fields,omitempty"`
	VerifyToken        string              `json:"verify_token,omitempty"`
}

type SubscriptionField struct {
	Name string `json:"name"`
}

// IncomingReaction represents a reaction to a message received via webhook.
type IncomingReaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji,omitempty"`
}

// SubscriptionResponse contains the list of subscribed webhook fields.
type SubscriptionResponse struct {
	Data     []*Subscription `json:"data,omitempty"`
	Fields   string          `json:"fields,omitempty"`
}
