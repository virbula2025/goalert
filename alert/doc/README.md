# GoAlert Alert Domain (`/alert`)

The `/alert` directory contains the core business logic, data structures, and persistence layers for Alerts in GoAlert. This package is responsible for the entire lifecycle of an incident—from creation and deduplication to escalation and resolution.

## Alert Lifecycle

An alert in GoAlert isn't just a record; it is a state machine that transitions based on system events or human interaction.



### 1. Creation & Deduplication
When an external system (via Integration Key, Email, or Heartbeat) reports an issue, the `AlertStore` first checks for "Deduplication."
- **Deduplication Logic (`dedup.go`)**: If an open alert already exists for the same service with a matching `dedup_key`, GoAlert increments a count instead of creating a new noisy notification.

### 2. State Management (`state.go`, `status.go`)
Alerts exist in one of three primary states:
- **Triggered**: The alert is active and needs attention.
- **Active (Acknowledged)**: A human has acknowledged the alert; automatic escalation is paused.
- **Closed (Resolved)**: The issue is fixed; notifications stop.

### 3. Escalation Tracking
The alert maintains its own "Escalation State." This tracks which step of the Escalation Policy it is currently on and when it is due to move to the next step.

## Alert Logging (`/alertlog`)

Every single action taken on an alert is recorded in the Alert Log. This is critical for auditing and "Post-Mortem" analysis.
- **Events**: "Created," "Escalated," "Acknowledged," "Notification Sent," etc.
- **Metadata**: Logs store the "Subject" (Who did it?) and "Message" (What happened?).



## Key Concepts

* **Summary & Details**: The `summary` is a short description (e.g., "CPU High"), while `details` can contain large blobs of diagnostic information.
* **Service Association**: Every alert must belong to exactly one `Service`. The Service determines which Escalation Policy is used.
* **Metrics**: The package tracks "Time to Acknowledge" (TTA) and "Time to Resolve" (TTR) to help teams measure their performance.