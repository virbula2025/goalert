# GoAlert Application Bootstrapping (`/app`)

The `/app` directory is the core orchestration layer of GoAlert. It functions as the "Main" package's engine, responsible for taking raw configuration and turning it into a running, interconnected system of database stores, API servers, and background workers.

If you want to understand how the various pieces of GoAlert (like `/engine`, `/graphql2`, and `/notification`) talk to each other, this is the place to look.

## Architectural Role

GoAlert follows a "Dependency Injection" pattern managed manually within the `App` struct. Instead of using global variables, the `App` struct holds references to every store and manager. This ensures that the entire lifecycle (startup, running, and shutdown) is predictable and testable.



---

## Core Component: The `App` Struct

Located in `app.go`, the `App` struct is the "State" of the running program. It includes:
- **Database Connections**: References to the PostgreSQL pool.
- **Stores**: Every domain store (e.g., `AlertStore`, `UserStore`, `ScheduleStore`).
- **Servers**: The HTTP server for the UI/API and the SMTP server for email integration.
- **Managers**: The background engine managers (e.g., `EscalationManager`).

---

## The Startup Lifecycle

The application starts through a series of `init` functions. This sequence is vital because many components depend on others (e.g., the GraphQL API cannot start until the UserStore is initialized).

### 1. Initialization (`init*.go` files)
- **`initstores.go`**: The first major step. It initializes the Go database handle and instantiates every store in the system.
- **`initauth.go`**: Sets up the authentication handlers (OIDC, Basic Auth) and links them to the UserStore.
- **`initgraphql.go`**: Plugs the GraphQL schema into the HTTP router.
- **`initengine.go`**: Starts the background workers that process alerts and notifications.
- **`initsmtpsrv.go`**: Bootstraps the internal SMTP server to allow GoAlert to receive emails directly.

### 2. The Execution Loop (`runapp.go`)
Once initialized, `app.Run()` is called. This function:
1. Starts the HTTP server.
2. Begins the background engine loops.
3. Listens for OS signals (like SIGTERM) to trigger a graceful shutdown.

### 3. Graceful Shutdown (`shutdown.go`)
GoAlert is designed for high availability. When shutting down:
- It stops accepting new HTTP requests.
- It allows the background engine managers to finish their current "processing cycle."
- It closes database connections cleanly to prevent orphaned locks.

---

## Key Files & Their Responsibilities

| File | Responsibility |
| :--- | :--- |
| `app.go` | Defines the central `App` struct and dependency list. |
| `startup.go` | The entry logic for starting the application components. |
| `initstores.go` | Connects to Postgres and sets up all data access layers. |
| `inithttp.go` | Configures the web server, middleware (Gzip, Logging), and routing. |
| `config.go` | Logic for loading and watching system configuration changes. |
| `lifecycle/` | (Sub-directory) Manages the "Pause/Resume" state of the app (used during Switchovers). |
| `middleware.go`| Contains logic for Request ID tracking, recovery from panics, and auth sessions. |

---

## Notable Logic: The "Switchover" Lifecycle

One of the most advanced features found in `/app` (and the `/app/lifecycle` sub-package) is the ability to **Pause** the application.

- **`pause.go`**: Implements the logic to temporarily halt the background engine. This is used when migrating the database to a new instance (Switchover) to ensure no data is processed while the database is being swapped.



---

## How to use this for learning

1.  **Trace a Store**: Open `initstores.go` and see how `app.AlertStore` is created. Follow it into the `/alert` directory to see how it's implemented.
2.  **Trace a Request**: Open `inithttp.go` to see the standard middleware. Then look at `initgraphql.go` to see how `/api/graphql` is mapped to the resolvers.
3.  **Trace the Engine**: Open `initengine.go` to see how the background managers you studied in `/engine` are actually started and passed the database handles.

---

## Summary for Developers
If you are adding a new core feature to GoAlert:
1.  Add your new **Store** or **Manager** to the `App` struct in `app.go`.
2.  Initialize it in a new or existing `init*.go` file.
3.  Inject it into the GraphQL resolvers in `initgraphql.go`.