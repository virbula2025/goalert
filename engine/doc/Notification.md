# GoAlert Notification System (`/notification`)

The notification system is responsible for the final "delivery" of alerts to users. It serves as an abstraction layer between GoAlert’s internal logic and external service providers like Twilio, Slack, and Mailgun.

## The Notification Workflow

The journey of a notification follows a specific lifecycle to ensure that no message is lost, even if a provider is temporarily down.

### 1. Generation
When the `escalationmanager` determines a user needs to be notified, it creates a record in the `outgoing_messages` table. This record is initially in a `pending` state.

### 2. The Engine Loop (`engine/message`)
The message engine constantly polls for pending messages. 
- It checks **User Contact Methods** to see how the user wants to be reached (SMS, Voice, Email, Slack).
- It checks **Notification Rules** to determine the timing (e.g., "SMS immediately, call if no response in 5 minutes").

### 3. Dispatch
The engine passes the message to the `/notification` package.
- **Provider Selection**: Based on the destination type, the system routes the message to the appropriate driver (`/twilio`, `/slack`, `/email`, etc.).
- **Message Formatting**: The `/nfymsg` package formats the raw alert data into a human-readable string or voice script tailored for the specific medium.

### 4. Handling Failures
If a provider returns an error (e.g., a "500 Internal Server Error" from an API), GoAlert follows this failure protocol:

* **Temporary Failures**: If the error is transient (network timeout, rate limit), the message state remains `pending` or moves to `retry`. The engine will attempt to send it again with exponential backoff.
* **Permanent Failures**: If the error is permanent (e.g., "Invalid Phone Number" or "User has blocked the app"), the message is marked as `failed`.
* **Alert Escalation**: A notification failure **does not** stop the alert. If the user doesn't acknowledge the alert (because they never got the message), the `escalationmanager` will eventually move the alert to the next person in the policy.



---

## Directory Structure & Responsibilities

| Directory | Purpose |
| :--- | :--- |
| `/twilio` | Handles SMS and Voice calls via the Twilio API. Includes TwiML generation for interactive phone menus. |
| `/slack` | Manages Slack DM notifications and interactive button responses (Ack/Close). |
| `/email` | Interfaces with SMTP or dedicated providers like Mailgun to send alert emails. |
| `/webhook` | Sends POST requests to external URLs for custom integrations. |
| `/nfydest` | Logic for validating and formatting "Destinations" (where the message goes). |
| `/nfymsg` | Templates for the content of the messages (e.g., "Alert #123: Database is down"). |

---

## Reliability Features

### Rate Limiting (`/util/calllimiter`)
To prevent overwhelming providers or being flagged as spam, GoAlert implements per-user and per-provider rate limits. If a system triggers 1,000 alerts at once, the engine will queue the notifications and bleed them out at a safe rate.

### Dedicated Retries
The `engine/npcyclemanager` (Notification Policy Cycle Manager) tracks the specific "step" a user is on. If a message fails, this manager ensures that the system doesn't just stall, but continues trying the user's secondary contact methods if configured.

### Feedback Loop
When a user interacts with a notification (e.g., pressing "4" on a phone call or clicking "Acknowledge" in Slack), the `/notification` drivers receive this callback and pass it back to the `engine/message/receiver.go` to update the alert status immediately.