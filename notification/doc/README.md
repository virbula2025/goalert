# GoAlert Notification System (`/notification`)

The `/notification` package is the "Dispatch Center" for GoAlert. It translates internal system events (like a service going down) into external communications (SMS, Voice, Slack, Email, or Webhooks) and manages the lifecycle of those messages.

## System Architecture

The notification system is built to be provider-agnostic. The core engine doesn't need to know how Twilio's API works; it simply asks the `Manager` to send a message to a "Destination."



### 1. The Manager & Registry
- **Manager (`manager.go`)**: Coordinates the flow of outgoing messages and incoming responses.
- **Registry (`nfydest/registry.go`)**: A central "phone book" that maps destination IDs (e.g., `builtin-webhook`) to the Go code that knows how to handle them.

### 2. Message Lifecycle
Messages are tracked through several states in the database via the `Store`:
- `Pending`: Queued for sending.
- `Sending`: Currently being transmitted to the provider.
- `Sent`: Accepted by the provider.
- `Delivered`: Confirmed arrival (if supported by provider).
- `Failed`: Permanent error (e.g., invalid number).

### 3. Interaction & Feedbacks
The system supports two-way communication. When a user presses a key on their phone or clicks a Slack button, the provider sends a callback to GoAlert. The `ResultReceiver` processes these actions to Acknowledge or Resolve alerts in real-time.

## Key Features

- **Rate Limiting**: Prevents "Notification Storms" from overwhelming providers or users.
- **Verification**: Handles the logic for sending and validating 6-digit codes to verify new contact methods.
- **Bundling**: Logic to group multiple alerts into a single notification to reduce noise.
- **Metrics**: Detailed Prometheus tracking for notification volume and delivery success.