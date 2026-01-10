# GoAlert Architecture and Code Structure Guide

## 1. High-Level Architecture

GoAlert is designed as a **monolithic** application. A single compiled binary (`goalert`) serves both the backend API and the frontend assets. It is designed to be stateless and horizontally scalable, relying heavily on the database for state and coordination.

* **Backend:** Written in **Go** (Golang). It handles API requests, business logic, and background processing.
* **Frontend:** A Single Page Application (SPA) written in **React** and **TypeScript**. It communicates with the backend via GraphQL.
* **Database:** **PostgreSQL** is the single source of truth. It handles data storage and is also used for coordination between multiple GoAlert instances (using locks) to ensure alerts are processed exactly once.

---

## 2. Directory Structure Breakdown

The repository is organized into domain-specific packages. Here is the map to help you navigate the code:

### Core Application Logic
* **`cmd/`**: The entry point for the compiled binary.
    * `cmd/goalert/main.go`: The `main` function. It parses command-line flags (like `--db-url`) and initializes the `app`.
* **`app/`**: The "glue" of the system.
    * It initializes the database connection, HTTP servers, and wires together the `engine` and API handlers.
    * Look at `app/app.go` or `app/init.go` to see how the system boots up.
* **`engine/`**: **The heart of the alerting system.**
    * This package contains the background worker logic. It runs a continuous "cycle" to check for:
        * New alerts.
        * Escalation policies that need to trigger.
        * Notifications that need to be sent (SMS, Voice, Slack).
    * It uses a **Cycle Monitor** pattern to poll the database efficiently.

### Domain Packages (Business Logic)
These packages contain the structs, interface definitions, and database logic for specific features:

* **`alert/`**: Defines what an "Alert" is, how it is created, closed, and deduplicated.
* **`schedule/`**: Logic for on-call schedules, shifts, and rotations.
* **`escalation/`**: Logic for escalation policies (who gets notified next if the first person doesn't answer).
* **`service/`**: Management of services (the entity that owns an integration and escalation policy).
* **`user/`**: User management, contact methods (phone numbers, emails), and notification rules.
* **`notification/`**: Handles the delivery of messages (Sender interfaces).

### Interfaces & API
* **`graphql2/`**: The implementation of the **GraphQL API**.
    * This is the primary API used by the frontend.
    * It contains the schema definitions and resolvers that map GraphQL queries to Go function calls.
* **`web/`**: The **Frontend** source code.
    * Contains the React application, TypeScript definitions, and UI components.
    * `web/src/` is where the actual frontend logic lives.

### Infrastructure & Utilities
* **`config/`**: Handles loading configuration from environment variables and flags.
* **`migrate/`**: Database migration files (SQL). These define the database schema and how it changes over time.
* **`auth/`**: Authentication providers (GitHub, OIDC, Basic Auth) and session management.
* **`integrationkey/`**: Handles incoming integrations (like Prometheus, Grafana, or generic webhooks) that trigger alerts.

---

## 3. Key Design Concepts & Implementation

To understand the code, you need to understand these three specific patterns GoAlert uses:

### A. The Engine "Cycle" (`engine` package)
Unlike systems that use complex message queues (like RabbitMQ), GoAlert keeps its architecture simple by using the database as a queue.

1.  **Polling Loop:** The `engine` runs a continuous loop.
2.  **Processing:** In every cycle, it asks the database: "Are there any alerts that have passed their escalation timeout?" or "Are there any outgoing messages stuck in the queue?"
3.  **Coordination:** If you run 5 instances of GoAlert, they use database locks (Advisory Locks) to ensure only one instance processes a specific alert cycle at a time. This prevents duplicate notifications.

### B. Data Access (`sqlc`)
GoAlert uses a tool called **sqlc** to generate Go code from SQL queries.

* You won't find many complex ORMs (like GORM).
* Instead, look for `.sql` files inside packages (e.g., `alert/queries.sql`).
* `sqlc` generates type-safe Go functions from these SQL files. This makes the database interactions very fast and explicit.

### C. The "Service" Model
The code is structured around "Services" (the Go interface pattern, not Microservices).

* For example, the `alert` package will expose an interface (Store or Repository) that the `graphql2` package consumes.
* **Flow:** Frontend (GraphQL) -> `graphql2` resolver -> `alert` package logic -> `sqlc` generated code -> PostgreSQL.

---

## 4. How to Read the Code (Recommended Path)

If you want to modify or understand a specific feature, follow this path:

1.  **Start at the Schema (`graphql2/schema.graphql`):** Find the mutation or query related to the feature (e.g., `createAlert`).
2.  **Find the Resolver:** Go to the `graphql2` package and find the Go function that handles that request.
3.  **Trace to Domain Logic:** The resolver will call a function in a domain package (e.g., `alert.Create`).
4.  **Check the Database:** The domain package will call a database method, often found in a `store.go` file or generated `db.go` file within that package.

---

## 5. Summary of Tech Stack

* **Language:** Go (Backend), TypeScript/React (Frontend)
* **API:** GraphQL
* **Data:** PostgreSQL (requires `pgcrypto` extension)
* **ORM:** None (uses `sqlc` for raw SQL to Go generation)
* **Build:** Docker & Make


# How GoAlert Handles High Alert Volumes (Scaling)

When a "ton of alerts" hits GoAlert (e.g., an alert storm where a monitoring system fires 1,000 requests in a minute), the system is designed to absorb the spike primarily through **deduplication** and **concurrency control**, rather than simply trying to send 1,000 SMS messages.

Here is how the system handles high volumes and how it scales to meet demand:

### 1. First Line of Defense: Intelligent Deduplication
The most critical mechanism for handling high alert volume is reducing "1,000 alerts" down to "1 actionable incident."

* **Ingestion:** When alerts hit the API, GoAlert checks the **`dedup_key`** (deduplication key).
* **Logic:**
    * If an alert with this key is already **Open** (Triggered), GoAlert **suppresses** the new alert. It logs it as a "duplicate" but does **not** create a new database record or trigger a new notification.
    * This transforms a "storm" of requests into a single database row and a single notification sequence.
    * *Code Location:* This logic is handled in the `alert` package, specifically during the alert creation flow (`alert.CreateOrUpdate`).

### 2. The Bottleneck: The "Engine" Processing Loop
If the alerts are *unique* (e.g., 1,000 *different* servers failing at once), GoAlert relies on its **Engine** to process them.

* **The Cycle:** The `engine` package runs a continuous loop (a "cycle"). It doesn't process alerts instantly in real-time interrupt style; instead, it polls the database for "work to be done" (e.g., "Find all alerts that need notifications sent right now").
* **Batching:** The engine processes records in batches. It reads a chunk of pending tasks from the database, processes them (sends SMS, updates escalation), and then looks for the next batch. This prevents the application memory from exploding during a spike.

### 3. Horizontal Scaling (Adding More Instances)
GoAlert is designed to be **horizontally scalable**. You can run multiple instances (replicas) of the `goalert` binary behind a Load Balancer to handle higher loads.

* **Stateless App:** The GoAlert binary itself is stateless. You can run 1, 10, or 50 instances.
* **Coordination (Locking):** You might wonder, *"If I have 10 instances, will they send 10 duplicate SMS messages?"*
    * **No.** They coordinate using **Database Locks** (specifically PostgreSQL Advisory Locks or `SELECT ... FOR UPDATE SKIP LOCKED` patterns).
    * When an instance's Engine Cycle picks up an alert to process, it "locks" that alert. Other instances see it is locked and skip it, moving on to the next one.
    * This allows multiple instances to churn through a massive queue of alerts in parallel without stepping on each other's toes.

### 4. Rate Limiting (Protection)
To prevent getting blocked by SMS/Voice providers (like Twilio or Mailgun) and to avoid spamming users:

* **Notification Rate Limits:** GoAlert implements logic to throttle outgoing notifications. If a user has a rule to "Notify immediately," but the system is overwhelmed, the `limit` package ensures it stays within safe API quotas for the providers.

### Summary: What happens during an Alert Storm?
1.  **Incoming:** 1,000 HTTP requests hit your Load Balancer.
2.  **Distribution:** The Load Balancer spreads these requests across your 5 GoAlert instances.
3.  **Deduplication:** The instances talk to the DB. 990 of the requests match an existing open alert and are ignored (suppressed).
4.  **Processing:** The remaining 10 unique alerts are inserted.
5.  **Notification:** The Engines on the 5 instances race to pick up these 10 alerts. They lock them, process the escalation policy, and send the notifications.

### The Ultimate Limit
Since GoAlert relies on the database for coordination and state, **PostgreSQL is the ultimate scalability limit**. To scale further, you would need to vertically scale your Postgres database (more CPU/RAM) or optimize your database disk I/O.