# GoAlert Core Application Files (`/app`)

This document summarizes the purpose and functionality of each source file within the `/app` directory. These files are responsible for the "wiring" of the GoAlert system.

## Application Structure & Lifecycle

| File | Description |
| :--- | :--- |
| **`app.go`** | Defines the central `App` struct. This is the "God Object" that holds references to all database stores, configuration, background managers, and server listeners. |
| **`runapp.go`** | The main execution entry point. It manages the high-level loop: starting servers, beginning the engine background work, and waiting for shutdown signals. |
| **`startup.go`** | Orchestrates the sequential startup process, ensuring that the database is connected before stores are initialized, and stores are initialized before the API starts. |
| **`shutdown.go`** | Handles the graceful shutdown of the system. It ensures active requests finish, engine workers stop cleanly, and database connections are closed. |
| **`pause.go`** | Implements the "Pause" functionality. This allows the application to stop all background processing (like sending alerts) without fully shutting down—critical for maintenance or switchovers. |

## Initialization Modules

These files follow a naming convention (`init*.go`) and are responsible for setting up specific subsystems.

| File | Subsystem Responsibility |
| :--- | :--- |
| **`initstores.go`** | Initializes all "Store" objects (e.g., AlertStore, UserStore, ScheduleStore). It injects the database connection pool into each store. |
| **`initengine.go`** | Bootstraps the background engine managers (Escalation, Notification, Heartbeat) and starts their processing loops. |
| **`initgraphql.go`** | Sets up the GraphQL API layer, linking the schema and resolvers to the application's stores. |
| **`initauth.go`** | Configures authentication providers (OIDC, GitHub, Basic Auth) and the session management system. |
| **`inithttp.go`** | Configures the primary HTTP server, defining the main router and static file serving for the React frontend. |
| **`initsmtpsrv.go`** | Starts the internal SMTP server used to receive emails that trigger alerts via Email Integration Keys. |
| **`initslack.go`** | Initializes the Slack bot client and the interactive Slack notification handlers. |
| **`inittwilio.go`** | Sets up the Twilio client used for SMS and Voice notification delivery. |
| **`initriver.go`** | Initializes the **River** job queue (if used) for background task processing. |
| **`initsysapi.go`** | Sets up the internal System API used for cross-instance communication. |

## Middleware & Networking

| File | Description |
| :--- | :--- |
| **`middleware.go`** | Contains the standard HTTP middleware stack: logging, panic recovery, and authentication session checking. |
| **`middlewaregzip.go`** | Implements Gzip compression for HTTP responses to save bandwidth on large GraphQL payloads. |
| **`secureheaders.go`** | Injects security-related HTTP headers (HSTS, X-Frame-Options) into every response. |
| **`multilistener.go`** | A utility that allows the app to listen on multiple ports or protocols simultaneously (e.g., HTTP and HTTPS). |
| **`tlsconfig.go`** | Manages SSL/TLS configurations for the web server and SMTP server. |

## Configuration & Observability

| File | Description |
| :--- | :--- |
| **`config.go`** | Handles the loading of system-wide settings and monitors the database for configuration changes in real-time. |
| **`getsetconfig.go`** | Provides the API internal methods for reading and updating settings through the UI or CLI. |
| **`metrics.go`** | Defines the Prometheus metrics for the application (e.g., request latency, active engine workers). |
| **`prometheus.go`** | Sets up the `/metrics` endpoint for scraping by monitoring tools. |
| **`healthcheck.go`** | Implements the `/health` endpoints used by load balancers and Kubernetes to check if the app is ready. |
| **`pprof.go`** | Enables Go's profiling tools for debugging performance bottlenecks in production. |

## Internal Utilities

| File | Description |
| :--- | :--- |
| **`context.go`** | Helpers for managing Go contexts, ensuring that timeouts and request IDs are propagated through the call stack. |
| **`trigger.go`** | Provides a mechanism to manually trigger engine cycles, often used in testing or via the Admin UI. |
| **`defaults.go`** | Defines the default values for system configuration if none are provided. |
| **`listenevents.go`** | Logic for listening to database events (e.g., Postgres NOTIFY) to wake up engine workers. |
| **`limitconcurrencybyauthsource.go`** | A specialized middleware to prevent a single auth provider (like a buggy OIDC setup) from overwhelming the system. |