package types

// Interactive holds interactive message content: buttons, list, or CTA URL.
type Interactive struct {
	Type   string         `json:"type"`
	Header *InteractiveHeader `json:"header,omitempty"`
	Body   *InteractiveBody   `json:"body,omitempty"`
	Footer *InteractiveFooter `json:"footer,omitempty"`
	Action *InteractiveAction `json:"action"`
}

// InteractiveHeader is an optional header for interactive messages (text, image, video, document).
type InteractiveHeader struct {
	Type     string  `json:"type"`
	Text     string  `json:"text,omitempty"`
	Image    *Media  `json:"image,omitempty"`
	Video    *Media  `json:"video,omitempty"`
	Document *Document `json:"document,omitempty"`
}

// InteractiveBody is the body text of an interactive message.
type InteractiveBody struct {
	Text string `json:"text"`
}

// InteractiveFooter is optional footer text.
type InteractiveFooter struct {
	Text string `json:"text"`
}

// InteractiveAction defines what happens when the user interacts: buttons, sections, or CTA params.
type InteractiveAction struct {
	Button    string          `json:"button,omitempty"`
	Buttons   []*Button       `json:"buttons,omitempty"`
	Sections  []*Section      `json:"sections,omitempty"`
	Name      string          `json:"name,omitempty"`
	Parameters *ActionParams  `json:"parameters,omitempty"`
	CatalogID string          `json:"catalog_id,omitempty"`
	ProductRetailerID string `json:"product_retailer_id,omitempty"`
	ProductItems []*ProductItem `json:"product_items,omitempty"`
}

// ActionParams holds CTA URL parameters.
type ActionParams struct {
	DisplayText string `json:"display_text"`
	URL         string `json:"url"`
}

// Button is a reply button for interactive messages.
type Button struct {
	Type  string      `json:"type"`
	Reply *ButtonReply `json:"reply"`
}

type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Section is a group of rows in a list message.
type Section struct {
	Title        string         `json:"title,omitempty"`
	Rows         []*Row         `json:"rows,omitempty"`
	ProductItems []*ProductItem `json:"product_items,omitempty"`
}

// Row is a single item in a list section.
type Row struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type ProductItem struct {
	ProductRetailerID string `json:"product_retailer_id"`
}

// NewInteractiveButton creates a message with up to 3 reply buttons.
func NewInteractiveButton(to, body string, buttons ...*Button) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: &Interactive{
			Type: "button",
			Body: &InteractiveBody{Text: body},
			Action: &InteractiveAction{Buttons: buttons},
		},
	}
}

// NewInteractiveList creates a message with a single-select list menu.
func NewInteractiveList(to, buttonText, body string, sections ...*Section) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: &Interactive{
			Type: "list",
			Body: &InteractiveBody{Text: body},
			Action: &InteractiveAction{
				Button:   buttonText,
				Sections: sections,
			},
		},
	}
}

// NewInteractiveCTA creates a message with a call-to-action URL button.
func NewInteractiveCTA(to, displayText, url, body string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: &Interactive{
			Type: "cta_url",
			Body: &InteractiveBody{Text: body},
			Action: &InteractiveAction{
				Name: "cta_url",
				Parameters: &ActionParams{
					DisplayText: displayText,
					URL:         url,
				},
			},
		},
	}
}

// NewButton creates a reply button with a unique ID and display title.
func NewButton(id, title string) *Button {
	return &Button{
		Type:  "reply",
		Reply: &ButtonReply{ID: id, Title: title},
	}
}

// NewSection creates a list section with a title and rows.
func NewSection(title string, rows ...*Row) *Section {
	return &Section{Title: title, Rows: rows}
}

// NewRow creates a list row with ID, title, and optional description.
func NewRow(id, title, description string) *Row {
	return &Row{ID: id, Title: title, Description: description}
}

func (ia *Interactive) WithHeader(headerType, text string) *Interactive {
	ia.Header = &InteractiveHeader{Type: headerType, Text: text}
	return ia
}

func (ia *Interactive) WithFooter(text string) *Interactive {
	ia.Footer = &InteractiveFooter{Text: text}
	return ia
}

func (ia *Interactive) WithImageHeader(mediaIDOrLink, caption string) *Interactive {
	ia.Header = &InteractiveHeader{Type: "image", Image: &Media{ID: mediaIDOrLink}}
	return ia
}
