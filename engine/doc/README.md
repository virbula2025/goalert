# GoAlert Engine (`/engine`)

The Engine is the central processing hub of GoAlert. It is responsible for the background automation that turns static configurations (schedules, escalation policies, and services) into dynamic actions (sending notifications, escalating alerts, and cleaning up old data).

Unlike the API layer which is reactive (waiting for a user request), the Engine is **proactive**—it constantly monitors the state of the database to determine what needs to happen next.

## Architecture & Design

The Engine is composed of several independent **Managers**. Each manager focuses on a specific domain of the application. These managers typically follow a "Polled Loop" or "Worker Queue" pattern:
1.  **Look for Work**: Query the database for records that need processing (e.g., an alert that hasn't been acknowledged).
2.  **Acquire Lock**: Use a distributed locking mechanism (`processinglock`) to ensure that only one instance of the engine processes a specific item at a time.
3.  **Perform Action**: Update the database, send a message, or move an alert to the next step.
4.  **Release/Sleep**: Wait for the next interval.

---

## Core Managers & Important Files

### 1. Escalation Manager (`/escalationmanager`)
The most critical component. It handles the logic of moving an alert through an Escalation Policy.
- **`escalationmanager.go`**: The main loop that identifies active alerts that are "due" for their next step.
- **Logic**: If an alert remains unacknowledged beyond the configured delay, this manager identifies the next user or schedule in the chain and creates new outgoing notifications.

### 2. Message Manager (`/message`)
Handles the lifecycle of a notification (SMS, Voice, Email, Slack).
- **`sender.go`**: Orchestrates the sending of messages. It interfaces with the `/notification` package to talk to external providers like Twilio or Slack.
- **`receiver.go`**: Handles incoming responses (e.g., a user replying "4" to an SMS to acknowledge an alert).
- **`scheduler.go`**: Manages the rate-limiting and queuing of messages to prevent flooding providers or users.

### 3. Schedule & Rotation Managers (`/schedulemanager`, `/rotationmanager`)
Handles the "Who is on call?" logic.
- **`rotationmanager.go`**: Calculates the state of rotations. If a rotation is set to "Weekly," this manager calculates the handoff time and updates the current participant.
- **`schedulemanager.go`**: Flattens complex schedules and overrides into a set of "on-call" segments in the database, which other managers then use to find targets for notifications.

### 4. Status Manager (`/statusmgr`)
Keeps users informed about the state of alerts they have interacted with.
- **`statusmgr.go`**: Sends "Status Updates" to users. For example, if User A acknowledges an alert, the Status Manager ensures User B (who was also notified) receives a message saying the alert is no longer active.

### 5. Heartbeat Manager (`/heartbeatmanager`)
Monitors the "Heartbeat" integration keys.
- **`heartbeatmanager.go`**: Looks for services that haven't "checked in" within their configured threshold. If a heartbeat is missed, this manager automatically triggers a new alert.



---

## Supporting Infrastructure

- **`/processinglock`**: A critical utility used by almost all managers. It uses PostgreSQL advisory locks to ensure that in a high-availability setup (multiple GoAlert instances), work is not duplicated and race conditions are avoided.
- **`/cleanupmanager`**: Handles data retention. It deletes old alerts, logs, and expired sessions to keep the database size manageable.
- **`/metricsmanager`**: Aggregates system data into metrics (e.g., "Alert Count" or "Notification Latency") for internal reporting and Prometheus export.
- **`/npcyclemanager` (Notification Policy Cycle)**: Manages the retry logic and timing for user-specific notification rules (e.g., "SMS me immediately, then Call me 5 minutes later").

---

## How to Add New Engine Logic

1.  **Define the Work**: Create a SQL query that finds "Outdated" or "Pending" rows.
2.  **Implement the Manager**: Create a struct that implements a `Run(ctx)` method.
3.  **Use Locks**: Ensure you use `processinglock.NewConfig` to prevent multiple engine nodes from processing the same row.
4.  **Register**: Add the new manager to the initialization logic in `app/app.go`.

---

## Summary of File Responsibilities

| Directory | Responsibility |
| :--- | :--- |
| `escalationmanager/` | Moving alerts from step 1 to step 2 to step N. |
| `message/` | Sending, receiving, and rate-limiting notifications. |
| `rotationmanager/` | Shifting people through on-call rotations. |
| `statusmgr/` | Sending updates when alerts are closed or acknowledged. |
| `heartbeatmanager/` | Monitoring external services for "I'm alive" signals. |
| `cleanupmanager/` | Deleting expired data and logs. |
| `processinglock/` | Distributed locking to prevent duplicate work. |