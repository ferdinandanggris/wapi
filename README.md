# wapi — WhatsApp Cloud API Go Library

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Zero Deps](https://img.shields.io/badge/dependencies-zero-0f0)](go.mod)

**wapi** is a production-ready Go library for [Meta WhatsApp Cloud API](https://developers.facebook.com/docs/whatsapp/cloud-api). Zero external dependencies, full API coverage, typed message builders, transport middleware, webhook handler, and reusable mock server.

## Features

- **No dependencies** — pure `net/http`, `encoding/json`, `mime/multipart`
- **Full API coverage** — messages (13 types), media, templates CRUD, phone management, business profile, webhook subscription
- **Typed builders** — `types.NewTextMessage()`, `types.NewInteractiveButton()`, etc.
- **Transport middleware** — retry with exponential backoff, token bucket rate limiter
- **Webhook handler** — `http.Handler` with signature verification, `OnMessage`/`OnStatus` callbacks
- **Reusable mock** — `cloud.MockServer` for integration testing without Meta
- **Clean interface** — `wapi.Client` in root package, implemented by `cloud.CloudClient`

## Installation

```bash
go get github.com/ferdinandanggris/wapi
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/ferdinandanggris/wapi/cloud"
    "github.com/ferdinandanggris/wapi/types"
)

func main() {
    ctx := context.Background()

    client := cloud.New(
        cloud.WithAccessToken(os.Getenv("WABA_TOKEN")),
    )

    phoneID := os.Getenv("PHONE_NUMBER_ID")
    msg := types.NewTextMessage("6281234567890", "Hello from wapi!", false)
    resp, err := client.SendMessage(ctx, phoneID, msg)
    if err != nil {
        fmt.Fprintf(os.Stderr, "send failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("sent! message id: %s\n", resp.Messages[0].ID)
}
```

## Usage Guide

### Send Messages

```go
// Text
types.NewTextMessage(to, "Hello!", false)
types.NewTextMessage(to, "Click here → https://example.com", true)

// Media by ID (uploaded) or Link
types.NewImageMessage(to, "media-id", "caption")
types.NewImageByLink(to, "https://...", "caption")
types.NewVideoMessage(to, "media-id", "caption")
types.NewAudioByLink(to, "https://...mp3")
types.NewDocumentByLink(to, "https://...pdf", "invoice.pdf", "Invoice")
types.NewStickerMessage(to, "media-id")
types.NewStickerByLink(to, "https://...webp")

// Location
types.NewLocationMessage(to, -6.2088, 106.8456, "Monas", "Jakarta Pusat")

// Reaction
types.NewReactionMessage(to, "original-msg-id", "👍")
types.NewRemoveReactionMessage(to, "original-msg-id")

// Reply to a specific message
msg := types.NewTextMessage(to, "This is a reply", false)
msg.WithContext("original-message-id")

// Mark as read
client.SendMessage(ctx, types.NewMarkAsRead("incoming-msg-id"))
```

### Interactive Messages

```go
// Reply buttons
msg := types.NewInteractiveButton(to, "Choose:",
    types.NewButton("yes", "Yes"),
    types.NewButton("no", "No"),
)
msg.Interactive.WithHeader("text", "Question").
    WithFooter("Reply within 24h")

// List
msg = types.NewInteractiveList(to, "View Menu", "Menu:",
    types.NewSection("Food",
        types.NewRow("nasi", "Nasi Goreng", "Rp 20,000"),
        types.NewRow("mie", "Mie Goreng", "Rp 18,000"),
    ),
    types.NewSection("Drinks",
        types.NewRow("kopi", "Kopi", "Rp 12,000"),
    ),
)

// CTA URL
msg = types.NewInteractiveCTA(to, "Visit Now",
    "https://example.com", "Click below:")
```

### Template Messages

```go
// Simple template (no variables)
msg := types.NewTemplateMessage(to, "hello_world", "en_US")

// Template with parameters
msg = types.NewTemplateMessage(to, "otp_template", "en_US",
    types.NewBodyComponent(
        types.NewTextParameter("123456"),
    ),
)

// List and send
templates, _ := client.ListTemplates(ctx, wabaID)
for _, t := range templates.Data {
    fmt.Printf("- %s (%s) [%s]\n", t.Name, t.Status, t.Category)
}
```

### Media Upload & Download

```go
// Upload (multipart/form-data)
file, _ := os.Open("image.png")
defer file.Close()

resp, _ := client.UploadMedia(ctx, phoneID, "image.png", file, "image/png")
fmt.Printf("media id: %s\n", resp.ID)

// Get download URL
info, _ := client.GetMediaURL(ctx, resp.ID)

// Download
reader, _ := client.DownloadMedia(ctx, resp.ID)
defer reader.Close()
io.Copy(os.Stdout, reader)

// Delete
client.DeleteMedia(ctx, resp.ID)
```

### Template CRUD

```go
// Create
tpl, _ := client.CreateTemplate(ctx, wabaID, &types.Template{
    Name:     "welcome",
    Language: "en_US",
    Category: "UTILITY",
    Components: []*types.TemplateComponent{
        {Type: "BODY", Text: "Welcome {{1}}!"},
    },
})

// Edit
client.EditTemplate(ctx, wabaID, tpl.ID, &types.Template{...})

// Get
tpl, _ = client.GetTemplate(ctx, tpl.ID)

// List with pagination
list, _ := client.ListTemplates(ctx, wabaID, wapi.WithLimit(20))

// Delete
client.DeleteTemplate(ctx, tpl.ID)
```

### Phone Number Management

```go
// Register
client.RegisterPhone(ctx, phoneID, "123456")

// Get details
phone, _ = client.GetPhoneNumber(ctx, phoneID)
fmt.Printf("quality: %s, limit: %s\n", phone.QualityRating, phone.MessagingLimit)

// List all numbers
phones, _ = client.ListPhoneNumbers(ctx, wabaID)

// Set 2-step PIN
client.SetTwoStepPIN(ctx, phoneID, "654321")

// Business profile
profile, _ = client.GetBusinessProfile(ctx, phoneID)
profile.Description = "New description"
client.UpdateBusinessProfile(ctx, phoneID, profile)
```

### Webhook Subscription

```go
// Subscribe
client.SubscribeToWebhooks(ctx, wabaID)

// Set callback URL + fields
client.SetWebhookCallback(ctx, appID,
    "https://your-server.com/webhook",
    "verify_token",
    "messages", "message_template_status",
)

// Set fields only
client.SetWebhookFields(ctx, appID, "messages")

// Check subscription
sub, _ = client.GetWebhookSubscription(ctx, wabaID)

// Unsubscribe
client.UnsubscribeFromWebhooks(ctx, wabaID)
```

## Webhook Handler

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/ferdinandanggris/wapi/types"
    "github.com/ferdinandanggris/wapi/webhook"
)

func main() {
    handler := &webhook.Handler{
        VerifyToken: os.Getenv("WEBHOOK_VERIFY_TOKEN"),
        AppSecret:   os.Getenv("WABA_APP_SECRET"),
        Logger:      log.Default(),
        OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
            fmt.Printf("[%s][%s] %s: %s\n", meta.DisplayPhoneNumber, msg.ID, contact.Profile.Name, msg.Text.Body)
            return nil
        },
        OnStatus: func(status *types.StatusUpdate, meta *types.Metadata) error {
            fmt.Printf("[%s] message %s: %s\n", meta.DisplayPhoneNumber, status.ID, status.Status)
            return nil
        },
    }

    mux := http.NewServeMux()
    mux.Handle("/webhook", handler)

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

The handler:
- **GET** `/webhook?hub.mode=subscribe&hub.verify_token=...&hub.challenge=...` — verification handshake
- **POST** `/webhook` — receives inbound messages and status updates
- Validates `X-Hub-Signature-256` when `AppSecret` is set

## Client Configuration

```go
client := cloud.New(
    cloud.WithAccessToken("EAAx..."),          // Required: WABA access token

    cloud.WithAPIVersion("v22.0"),             // Optional: default v21.0
    cloud.WithBaseURL("https://graph.facebook.com"), // Optional
    cloud.WithHTTPClient(customHTTPClient),    // Optional: custom *http.Client
    cloud.WithRetry(5),                        // Optional: retry up to 5 times
)

// phoneNumberID is passed per-call:
// client.SendMessage(ctx, phoneID, msg)
// client.MarkAsRead(ctx, phoneID, msgID)
// client.UploadMedia(ctx, phoneID, ...)
```

Recommended environment variables:

```
WABA_TOKEN=EAAT...
PHONE_NUMBER_ID=123456...
WABA_ID=123456...
APP_SECRET=abc123...
WEBHOOK_VERIFY_TOKEN=my_token
```

## Error Handling

```go
resp, err := client.SendMessage(ctx, msg)
if err != nil {
    var e *wapi.Error
    if errors.As(err, &e) {
        fmt.Printf("code=%d, type=%s, trace=%s\n", e.Code, e.Type, e.FBTraceID)
        if e.Type == wapi.ErrRateLimit {
            // back off
        }
    }
}

// Helper functions
wapi.IsRetryable(err)   // true for 5xx, 130429, 131056
wapi.IsRateLimit(err)    // true for code 130429
```

Error type classification:

| Type | Code | Meaning |
|---|---|---|
| `ErrOAuth` | varies | Token expired, invalid permissions |
| `ErrGraphMethod` | varies | Invalid request, missing fields |
| `ErrRateLimit` | 130429 | Too many requests |
| `ErrServer` | 500+ | Meta server error |
| `ErrUnknown` | — | Unclassified |

## Transport Middleware

Built-in middleware chain via `transport.Chain`:

```go
// Default: retry (3 attempts, 1s–60s backoff)
client := cloud.New(...)  // retry included by default

// Custom retry
client = cloud.New(cloud.WithRetry(5))

// Custom chain with rate limiter
hc := &http.Client{
    Transport: transport.Chain(
        http.DefaultTransport,
        transport.Retry(transport.RetryConfig{MaxAttempts: 5}),
        transport.NewRateLimiter(100, time.Second),
    ),
}
client = cloud.New(cloud.WithHTTPClient(hc))
```

## Mock Server

```go
import "github.com/ferdinandanggris/wapi/cloud"

func TestSomething(t *testing.T) {
    mock := cloud.NewMockServer(cloud.MockConfig{})
    defer mock.Close()

    client := cloud.New(
        cloud.WithBaseURL(mock.URL()),
        cloud.WithAccessToken("test-token"),
    )

    // Use client as normal — mock handles all endpoints
    // Phone number ID is passed per-call: client.SendMessage(ctx, "test-phone", msg)
}
```

Supports customizable HTTP status codes for testing error scenarios (rate limits, validation errors, etc.).

## Examples

| Example | Description |
|---|---|
| `examples/cloud_send` | Basic text and image messages |
| `examples/send_template` | Send a template message |
| `examples/send_interactive` | Interactive buttons and list |
| `examples/upload_media` | Multipart media upload |
| `examples/bulk_send` | Send messages to multiple recipients |
| `examples/webhook_server` | Full webhook server with ngrok |
| `examples/setup_webhook` | Subscribe + set callback URL |
| `examples/setup_callback` | Set callback URL only |

Run an example:

```bash
export WABA_TOKEN=EAAx...
export PHONE_NUMBER_ID=123456...
export TO_NUMBER=6281234567890

go run examples/cloud_send/main.go
```

## Project Structure

```
.
├── client.go          # Client interface
├── errors.go          # Error types + helpers
├── cloud/
│   ├── client.go      # CloudClient implementation
│   ├── messages.go    # SendMessage, MarkAsRead
│   ├── media.go       # UploadMedia, GetMediaURL, DownloadMedia, DeleteMedia
│   ├── templates.go   # Template CRUD
│   ├── phone.go       # Phone registration + profile
│   ├── webhooks.go    # Webhook subscription API
│   └── mock_test.go   # Reusable mock server
├── types/
│   ├── messages.go    # Message struct + builders
│   ├── interactive.go # Interactive + builders
│   ├── templates.go   # Template + components
│   ├── media.go       # Media response types
│   ├── phone.go       # Phone + business profile types
│   └── webhooks.go    # Inbound payload + subscription types
├── transport/
│   ├── transport.go   # Middleware chain
│   ├── retry.go       # Exponential backoff retry
│   └── ratelimit.go   # Token bucket rate limiter
├── webhook/
│   ├── handler.go     # HTTP handler (verify + handle)
│   └── signature.go   # X-Hub-Signature-256 verification
└── examples/          # Runnable example programs
```

## Development

```bash
make test    # CGO_ENABLED=1 go test -race ./...
make lint    # golangci-lint run
make vet     # go vet ./...
make fmt     # gofmt -w .
```

Install golangci-lint:
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.60.1
```

Lint config (`.golangci.yml`):
- `errcheck` — no unchecked errors
- `govet` — report suspicious constructs
- `staticcheck` — static analysis
- `exhaustive` — check exhaustiveness of type switches
- `noctx` — ensure context propagation in HTTP requests

## License

MIT
