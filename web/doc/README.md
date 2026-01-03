# GoAlert Web UI (`/web`)

The `/web` directory contains the source code for the GoAlert user interface. It is split into a Go-based backend (responsible for asset serving and HTML templating) and a React-based Single Page Application (SPA).

## Architecture

GoAlert uses a modern frontend architecture built on React, TypeScript, and GraphQL.



### 1. The Go Web Handler
The Go files in the root of `/web` manage the HTTP delivery of the frontend.
- **Embedded Assets**: The compiled React bundle is embedded into the Go binary using the `embed` package.
- **Index Template**: `index.go` and `index.html` handle the initial page load, injecting configuration variables (like the application name and theme) into the global window object.
- **Cache Management**: `nocache.go` and `etaghandler.go` ensure that users always have the latest version of the UI while optimizing browser performance.

### 2. The React Frontend (`/src/app`)
The frontend is a TypeScript application that uses **urql** for GraphQL communication and **Material UI** for the design system.
- **Routing**: Client-side routing is handled via `react-router`.
- **State Management**: Most state is managed through the GraphQL cache, with local UI state handled via React Hooks.
- **Themes**: Supports dynamic light/dark mode switching via `theme/`.

### 3. Testing Infrastructure
GoAlert has a heavy investment in frontend reliability:
- **Cypress (`/src/cypress`)**: End-to-end (E2E) tests that simulate real user interactions, from logging in to creating complex on-call schedules.
- **Storybook**: Used for isolated component development and visual testing.

## UI Organization

The application is modularized by feature:
- **`/alerts`**: Alert lists, details, and creation dialogs.
- **`/schedules`**: Calendar views, rotation management, and temporary overrides.
- **`/users`**: Profile management, contact methods, and notification rules.
- **`/admin`**: System-wide configuration, limits, and message logs.



## Development Workflow

1. **Frontend Dev**: Run the development server in the `web/src` directory (typically via `npm start`).
2. **GraphQL Codegen**: When changing the `.graphql` schema, the frontend types are regenerated to maintain type safety.
3. **Build**: The `esbuild` or `webpack` process bundles the TSX into a static JS/CSS set, which is then picked up by the Go `embed` directive.