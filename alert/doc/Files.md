# File Summary for `/alert`

This document summarizes the role of each file in the alert management subsystem.

## Core Domain & Store

| File | Description |
| :--- | :--- |
| **`alert.go`** | Defines the primary `Alert` struct and basic validation/normalization logic. |
| **`store.go`** | The main database access layer. Contains methods for creating, updating, and retrieving alerts. |
| **`queries.sql`** | Raw SQL queries (managed by `sqlc`) for high-performance alert operations. |
| **`search.go`** | Implements complex filtering logic for the UI, such as searching by status, service, or user. |
| **`state.go`** | Manages the internal state machine logic for an alert's lifecycle. |
| **`status.go`** | Defines the `Status` type (Triggered, Active, Closed). |
| **`metadata.go`** | Logic for handling arbitrary JSON metadata attached to alerts. |

## Feature Extensions

| File | Description |
| :--- | :--- |
| **`dedup.go`** | Logic for preventing duplicate alerts based on a unique `dedup_key`. |
| **`summary.go`** | Utility for generating and validating alert summaries. |
| **`feedback.go`** | Handles user-provided feedback on alerts (e.g., "Was this alert helpful?"). |
| **`logentryfetcher.go`**| A specialized utility to batch-fetch log entries for multiple alerts efficiently. |
| **`source.go`** | Tracks the origin of an alert (e.g., "Email," "Generic API," "Grafana"). |
| **`nfydest.go`** | Maps alert events to notification targets. |

## Monitoring & Metrics

| File | Description |
| :--- | :--- |
| **`metrics.go`** | Defines Prometheus metrics for alert counts and state transitions. |
| **`alertmetrics/`** | (Sub-directory) Stores historical performance data like Mean Time to Resolve. |

## Logging Sub-package (`/alert/alertlog`)

| File | Description |
| :--- | :--- |
| **`entry.go`** | Defines what an individual log entry looks like. |
| **`store.go`** | Persistence layer for writing and reading alert audit logs. |
| **`type.go`** | Defines the enum for log events (e.g., `LogEventEscalated`). |
| **`subject.go`** | Logic for identifying the actor in a log (e.g., a specific User ID or "System"). |
| **`rawjson.go`** | Helper for processing raw JSON metadata within logs. |