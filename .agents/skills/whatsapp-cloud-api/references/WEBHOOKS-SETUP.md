# WhatsApp Cloud API - Webhook Setup Reference

> API Version: v21.0

## Overview

Webhook setup involves two parts:
1. **Subscribe** your WhatsApp Business Account (WABA) to your Meta App
2. **Configure** the subscription to specify which events to receive

## Endpoints

### Subscribe WABA to App

Subscribes your WhatsApp Business Account to receive webhook events.

```
POST /{waba-id}/subscribed_apps
Authorization: Bearer {access-token}
```

**Success Response:**
```json
{
  "id": "APP_ID",
  "name": "My App"
}
```

### Unsubscribe WABA from App

Stops all webhook delivery to your app.

```
DELETE /{waba-id}/subscribed_apps
Authorization: Bearer {access-token}
```

**Success Response:**
```json
{
  "success": true
}
```

### Get Subscription Status

Check if your WABA is currently subscribed.

```
GET /{waba-id}/subscribed_apps
Authorization: Bearer {access-token}
```

### Set Subscription Fields

Specify which event types to receive. Available fields:
- `messages` — Incoming messages and status updates
- `message_template_status` — Template approval/rejection updates
- `account_review_update` — Business account review status changes
- `message_template_quality_update` — Template quality changes

```
POST /{app-id}/subscriptions
Content-Type: application/x-www-form-urlencoded
Authorization: Bearer {access-token}

object=whatsapp_business_account
fields=messages,message_template_status
```

### Set Callback URL + Fields

Configure where Meta sends webhook events and which events to subscribe to.

```
POST /{app-id}/subscriptions
Content-Type: application/x-www-form-urlencoded
Authorization: Bearer {access-token}

object=whatsapp_business_account
callback_url=https://your-server.com/webhook
verify_token=your-verify-token
fields=messages,message_template_status
```

## Library Usage

```go
// Subscribe WABA to receive webhooks
app, _ := client.SubscribeToWebhooks(ctx, "waba-id")

// Unsubscribe
client.UnsubscribeFromWebhooks(ctx, "waba-id")

// Check subscription status
app, _ := client.GetWebhookSubscription(ctx, "waba-id")

// Set which fields to receive
client.SetWebhookFields(ctx, "app-id", "messages", "message_template_status")

// Set callback URL + verify token + fields in one call
client.SetWebhookCallback(ctx, "app-id",
    "https://myserver.com/webhook",
    "my-verify-token",
    "messages", "message_template_status",
)
```

## Notes

- The `waba-id` is your WhatsApp Business Account ID (from Business Manager)
- The `app-id` is your Meta App ID (from Meta Developer Dashboard)
- If using on-premise, webhook subscription is managed differently
- Even after subscribing, Meta must verify your callback URL via the GET challenge flow
