# File Summary for `/notification`

This list covers all source files identified in the `/notification` root and its core sub-packages.

## Root Package (`/notification`)

| File | Purpose |
| :--- | :--- |
| **`manager.go`** | The primary orchestrator. Routes outgoing messages to providers and incoming results to the engine. |
| **`store.go`** | Database access layer for messages, verification codes, and rate-limit counters. |
| **`search.go`** | Logic for querying notification history (used in the Admin logs and User profile). |
| **`messagetype.go`** | Defines constants for message categories (Alert, StatusUpdate, Verification, Test). |
| **`result.go`** | Definitions for user actions (Acknowledge, Resolve, Escalate). |
| **`resultreceiver.go`** | Interface for components that handle user feedback from notifications. |
| **`sender.go`** | Interfaces for message sending and setting up two-way communication. |
| **`destid.go`** | Utilities for parsing and generating unique Destination IDs. |
| **`namedreceiver.go`** | Helper to associate human-readable names with notification targets. |
| **`metrics.go`** | Prometheus instrumentation for tracking delivery performance. |
| **`compat.go`** | Compatibility layer for bridging various internal message formats. |
| **`result_string.go`** | Generated stringer methods for notification results. |

## Destination Logic (`/notification/nfydest`)

| File | Purpose |
| :--- | :--- |
| **`registry.go`** | Manages the global list of active notification providers. |
| **`provider.go`** | The base interface that every notification channel (Slack, Twilio, etc.) must implement. |
| **`typeinfo.go`** | Metadata about providers (e.g., "Does this support Voice?"). |
| **`validate.go`** | Logic for validating destination strings (e.g., checking if a phone number is valid). |
| **`validateaction.go`** | Validates custom actions for dynamic providers like Webhooks. |
| **`messagestatuser.go`** | Interface for providers that support asynchronous status updates (delivery receipts). |
| **`messagesender.go`** | Interface for the actual "Send" operation of a provider. |

## Message Formatting (`/notification/nfymsg`)

| File | Purpose |
| :--- | :--- |
| **`message.go`** | Base interface for all message content objects. |
| **`msgalert.go`** | Formats data for individual alert notifications. |
| **`msgalertbundle.go`** | Logic for summarizing multiple alerts into a single "Bundle" message. |
| **`msgalertstatus.go`** | Formats updates regarding an alert's lifecycle (Ack/Close). |
| **`msgverification.go`** | Formats the 6-digit verification code messages. |
| **`msgtest.go`** | Formats simple "Test" notifications. |
| **`sendmessage.go`** | Captures the output of a send attempt (External ID, source value). |

## Core Providers

| Sub-package | Responsibility |
| :--- | :--- |
| **`/twilio`** | Full implementation of SMS and Voice. Includes TwiML generation and status callback validation. |
| **`/slack`** | Handles Slack DMs, Channels, and User Groups using Slack's "Block Kit." |
| **`/webhook`** | Sends custom JSON payloads to external HTTP endpoints. |
| **`/email`** | SMTP-based delivery for alert emails. |