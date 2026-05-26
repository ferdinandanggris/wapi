package types

// Message represents an outbound WhatsApp message.
// Use the New*Message builder functions to create typed messages.
type Message struct {
	MessagingProduct string             `json:"messaging_product"`
	RecipientType    string             `json:"recipient_type,omitempty"`
	To               string             `json:"to"`
	Type             string             `json:"type"`
	Context          *Context           `json:"context,omitempty"`
	Text             *Text              `json:"text,omitempty"`
	Image            *Media             `json:"image,omitempty"`
	Video            *Media             `json:"video,omitempty"`
	Audio            *Media             `json:"audio,omitempty"`
	Document         *Document          `json:"document,omitempty"`
	Sticker          *Media             `json:"sticker,omitempty"`
	Location         *Location          `json:"location,omitempty"`
	Contacts         []*Contact         `json:"contacts,omitempty"`
	Reaction         *Reaction          `json:"reaction,omitempty"`
	Interactive      *Interactive       `json:"interactive,omitempty"`
	Template         *TemplateMessage   `json:"template,omitempty"`
	Status           string             `json:"status,omitempty"`
}

// SendResponse is returned by SendMessage on success.
type SendResponse struct {
	MessagingProduct string       `json:"messaging_product"`
	Contacts         []*Contact   `json:"contacts"`
	Messages         []*MessageID `json:"messages"`
}

// MessageID contains the WhatsApp message ID.
type MessageID struct {
	ID string `json:"id"`
}

// Context creates a reply to an existing message in a conversation thread.
type Context struct {
	MessageID string `json:"message_id"`
}

// Text holds the body of a text message with optional URL preview.
type Text struct {
	Body       string `json:"body"`
	PreviewURL bool   `json:"preview_url,omitempty"`
}

// Media identifies uploaded media by ID or external URL, with optional caption.
type Media struct {
	ID      string `json:"id,omitempty"`
	Link    string `json:"link,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// Document extends Media with a display filename.
type Document struct {
	ID       string `json:"id,omitempty"`
	Link     string `json:"link,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// Location sends a map pin with optional name and address.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// Contact represents a contact message.
type Contact struct {
	Input     string     `json:"input,omitempty"`
	WaID      string     `json:"wa_id,omitempty"`
	Name      *Name      `json:"name,omitempty"`
	Phones    []*Phone   `json:"phones,omitempty"`
	Emails    []*Email   `json:"emails,omitempty"`
	URLs      []*URL     `json:"urls,omitempty"`
	Addresses []*Addr    `json:"addresses,omitempty"`
	Org       *Org       `json:"org,omitempty"`
	Birthday  string     `json:"birthday,omitempty"`
}

type Name struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
}

type Phone struct {
	Phone string `json:"phone"`
	Type  string `json:"type,omitempty"`
	WaID  string `json:"wa_id,omitempty"`
}

type Email struct {
	Email string `json:"email"`
	Type  string `json:"type,omitempty"`
}

type URL struct {
	URL  string `json:"url"`
	Type string `json:"type,omitempty"`
}

type Addr struct {
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Org struct {
	Company  string `json:"company,omitempty"`
	Title    string `json:"title,omitempty"`
	Department string `json:"department,omitempty"`
}

// Reaction sends an emoji reaction to an existing message.
// Use empty emoji to remove a reaction.
type Reaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

// NewTextMessage creates a text message. Set previewURL=true to enable link previews.
func NewTextMessage(to, body string, previewURL bool) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text:             &Text{Body: body, PreviewURL: previewURL},
	}
}

// NewImageMessage sends an image by uploaded media ID.
func NewImageMessage(to, linkOrID, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "image",
		Image:            &Media{ID: linkOrID, Caption: caption},
	}
}

// NewImageByLink sends an image by external URL.
func NewImageByLink(to, link, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "image",
		Image:            &Media{Link: link, Caption: caption},
	}
}

// NewVideoMessage sends a video by uploaded media ID.
func NewVideoMessage(to, linkOrID, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "video",
		Video:            &Media{ID: linkOrID, Caption: caption},
	}
}

// NewVideoByLink sends a video by external URL.
func NewVideoByLink(to, link, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "video",
		Video:            &Media{Link: link, Caption: caption},
	}
}

// NewAudioMessage sends audio by uploaded media ID.
func NewAudioMessage(to, linkOrID string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "audio",
		Audio:            &Media{ID: linkOrID},
	}
}

// NewAudioByLink sends audio by external URL.
func NewAudioByLink(to, link string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "audio",
		Audio:            &Media{Link: link},
	}
}

// NewDocumentMessage sends a document by uploaded media ID.
func NewDocumentMessage(to, linkOrID, filename, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "document",
		Document:         &Document{ID: linkOrID, Filename: filename, Caption: caption},
	}
}

// NewDocumentByLink sends a document by external URL.
func NewDocumentByLink(to, link, filename, caption string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "document",
		Document:         &Document{Link: link, Filename: filename, Caption: caption},
	}
}

// NewStickerMessage sends a sticker by uploaded media ID (WebP, 512x512 recommended).
func NewStickerMessage(to, linkOrID string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "sticker",
		Sticker:          &Media{ID: linkOrID},
	}
}

// NewStickerByLink sends a sticker by external WebP URL.
func NewStickerByLink(to, link string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "sticker",
		Sticker:          &Media{Link: link},
	}
}

// NewLocationMessage sends a map location with name and optional address.
func NewLocationMessage(to string, lat, lng float64, name, address string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "location",
		Location:         &Location{Latitude: lat, Longitude: lng, Name: name, Address: address},
	}
}

// NewContactMessage sends a formatted contact card.
func NewContactMessage(to string, contacts ...*Contact) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "contacts",
		Contacts:         contacts,
	}
}

// NewReactionMessage sends an emoji reaction to a message.
func NewReactionMessage(to, messageID, emoji string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "reaction",
		Reaction:         &Reaction{MessageID: messageID, Emoji: emoji},
	}
}

// NewRemoveReactionMessage removes an emoji reaction from a message.
func NewRemoveReactionMessage(to, messageID string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "reaction",
		Reaction:         &Reaction{MessageID: messageID, Emoji: ""},
	}
}

// NewMarkAsRead creates a "read" status update for an incoming message.
func NewMarkAsRead(messageID string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		Type:             "action",
		Status:           "read",
		Context:          &Context{MessageID: messageID},
	}
}

// WithContext sets the message ID to reply to. Returns the message for chaining.
func (m *Message) WithContext(messageID string) *Message {
	m.Context = &Context{MessageID: messageID}
	return m
}
