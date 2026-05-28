package types

import "encoding/json"

// TemplateMessage references a pre-approved message template with optional parameters.
type TemplateMessage struct {
	Name       string                `json:"name"`
	Language   *TemplateLanguage     `json:"language"`
	Components []*TemplateMsgComponent `json:"components,omitempty"`
}

// TemplateLanguage specifies the template language code (e.g., "en_US").
type TemplateLanguage struct {
	Code string `json:"code"`
}

// TemplateMsgComponent is a component of a template message (header, body, button).
type TemplateMsgComponent struct {
	Type       string                  `json:"type"`
	SubType    string                  `json:"sub_type,omitempty"`
	Index      int                     `json:"index,omitempty"`
	Parameters []*TemplateParameter    `json:"parameters,omitempty"`
}

// HasIndex reports whether the component has a meaningful index (for button components).
// Non-button components (body, header, footer) do not have an index.
func (c *TemplateMsgComponent) HasIndex() bool {
	return c.Type == "button"
}

// MarshalJSON customizes JSON output to include index 0 for button components.
func (c *TemplateMsgComponent) MarshalJSON() ([]byte, error) {
	type alias TemplateMsgComponent
	a := alias(*c)

	if c.HasIndex() {
		// Build manually to force index into output
		m := map[string]interface{}{
			"type":    c.Type,
			"sub_type": c.SubType,
			"index":   c.Index,
		}
		if c.Parameters != nil {
			m["parameters"] = c.Parameters
		}
		return json.Marshal(m)
	}
	return json.Marshal(a)
}

// TemplateParameter is a value for template variables (text, image, video, document).
type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Image *Media `json:"image,omitempty"`
	Video *Media `json:"video,omitempty"`
	Document *Document `json:"document,omitempty"`
}

// Template is a message template definition for CRUD operations.
type Template struct {
	ID         string              `json:"id,omitempty"`
	Name       string              `json:"name,omitempty"`
	Language   string              `json:"language,omitempty"`
	Category   string              `json:"category,omitempty"`
	Status     string              `json:"status,omitempty"`
	Components []*TemplateComponent `json:"components,omitempty"`
}

type TemplateComponent struct {
	Type    string             `json:"type"`
	Format  string             `json:"format,omitempty"`
	Text    string             `json:"text,omitempty"`
	Example json.RawMessage    `json:"example,omitempty"`
	Buttons []*TemplateButton  `json:"buttons,omitempty"`
}

type TemplateButton struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
	OTPType string `json:"otp_type,omitempty"`
}

// TemplateList is the paginated response from listing templates.
type TemplateList struct {
	Data   []*Template `json:"data"`
	Paging *Paging     `json:"paging,omitempty"`
}

type Paging struct {
	Cursors  *Cursors `json:"cursors,omitempty"`
	Next     string   `json:"next,omitempty"`
	Previous string   `json:"previous,omitempty"`
}

type Cursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

// NewTemplateMessage creates a template message. Use NewBodyComponent, NewHeaderComponent for parameters.
func NewTemplateMessage(to, name, lang string, components ...*TemplateMsgComponent) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: &TemplateMessage{
			Name:       name,
			Language:   &TemplateLanguage{Code: lang},
			Components: components,
		},
	}
}

// NewTextParameter creates a text parameter for template variable substitution.
func NewTextParameter(text string) *TemplateParameter {
	return &TemplateParameter{Type: "text", Text: text}
}

// NewBodyComponent creates a template body component with parameters for {{1}}, {{2}}, etc.
func NewBodyComponent(params ...*TemplateParameter) *TemplateMsgComponent {
	return &TemplateMsgComponent{Type: "body", Parameters: params}
}

// NewHeaderComponent creates a template header component with parameters.
func NewHeaderComponent(params ...*TemplateParameter) *TemplateMsgComponent {
	return &TemplateMsgComponent{Type: "header", Parameters: params}
}

// NewURLButtonComponent creates a URL button component with a dynamic path suffix.
func NewURLButtonComponent(index int, suffix string) *TemplateMsgComponent {
	return &TemplateMsgComponent{
		Type:    "button",
		SubType: "url",
		Index:   index,
		Parameters: []*TemplateParameter{
			{Type: "text", Text: suffix},
		},
	}
}
