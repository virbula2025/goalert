# Contributing to GoAlert: Telnyx & Advanced Routing

This guide outlines the steps to add Telnyx as a notification provider and implement multi-number call-in routing.

## 1. Adding Telnyx Support

To add Telnyx, you will need to create a new package `notification/telnyx` and register it in the application.

### Key Files to Reference:
- `notification/twilio/config.go`: See how Twilio's API credentials and settings are managed.
- `notification/twilio/client.go`: Reference this for the API client implementation.
- `notification/twilio/nfydestsms.go` & `nfydestvoice.go`: These implement the `nfydest.Sender` interface.

### Steps:
1. **Define Configuration**: Add Telnyx-specific fields (API Key, Public Key) to the global config in `config/config.go`.
2. **Implement Sender Interface**: Create `notification/telnyx/sender.go` that implements:
   - `Send(ctx context.Context, msg notification.Message) (*notification.SentMessage, error)`
3. **Registration**: In `app/app.go` or a new `app/inittelnyx.go`, initialize the Telnyx sender and register it with the notification manager.

## 2. Multi-Number Call-In Routing

"Call-in" refers to a user calling a GoAlert number to check or acknowledge alerts. Routing different numbers to different teams requires mapping incoming "To" numbers to specific Service IDs or Escalation Policies.

### How to Implement:
1. **Identify Incoming Request**: Twilio (and likely Telnyx) sends a POST request to your `/v1/twilio/voice` (or new `/v1/telnyx/voice`) endpoint.
2. **Look up Destination**: In the voice handler (see `notification/twilio/voice.go`), use the `To` phone number from the request to look up the associated GoAlert entity.
   - You may need to add a new table or configuration field to map "Phone Number -> Service/Group".
3. **Route Logic**:
   - If the number belongs to "Team A", fetch the current on-call person for Team A's schedule.
   - Use the `VoiceResponse` (TwiML equivalent) to prompt the user or bridge the call.

## 3. Recommended Learning Path

### Phase 1: Exploration
- [cite_start]**Study `notification/manager.go`**: Understand how messages are queued and dispatched. [cite: 1, 12, 13]
- [cite_start]**Examine `notification/twilio/twiml.go`**: Learn how GoAlert generates voice instructions. [cite: 21, 32, 33]

### Phase 2: Implementation
- Start by implementing **Telnyx SMS**. It is simpler than Voice because it doesn't require complex state machines or interactive menus.
- Move to **Telnyx Voice** once SMS is stable.
- Finally, implement the **Routing Table** for multiple incoming numbers.

## 4. Testing
GoAlert uses "Smoke Tests" for provider integrations.
- Reference `test/smoke/twiliosms_test.go` and `test/smoke/twiliovoice_test.go`.
- [cite_start]You should create `test/smoke/telnyxsms_test.go` using a similar harness to mock the Telnyx API. [cite: 1, 21]