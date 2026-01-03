# File Summary for `/web`

This document summarizes the files found in the web subsystem, covering both the Go delivery layer and the React source tree.

## Go Backend (Root of `/web`)

| File | Purpose |
| :--- | :--- |
| **`handler.go`** | The main HTTP handler that routes requests for static assets and the SPA index. |
| **`index.go`** | Renders the primary `index.html` template and injects dynamic server-side configuration. |
| **`index.html`** | The HTML entry point for the React application. |
| **`nocache.go`** | Middleware to prevent browsers from caching sensitive or rapidly changing UI files. |
| **`etaghandler.go`** | Logic for ETag-based caching to improve performance for static JS/CSS bundles. |
| **`explore.go`** | Logic for the `/explore` endpoint, providing a GraphQL sandbox for developers. |

## React Frontend Logic (`/web/src/app`)

| Directory/File | Responsibility |
| :--- | :--- |
| **`index.tsx`** | The React entry point; sets up the Apollo/Urql client and Material UI theme providers. |
| **`apollo.js` / `urql.ts`** | Configuration for the GraphQL clients and caching policies. |
| **`/alerts`** | Components for viewing, acknowledging, and resolving alerts. |
| **`/schedules`** | Complex calendar components and logic for managing on-call rotations. |
| **`/users`** | Management for user profiles, notification rules, and contact methods. |
| **`/admin`** | Admin-only views for system configuration, SMS/Voice limits, and service metrics. |
| **`/util`** | Shared UI utilities: Date/Time formatting, custom hooks, and common React components. |
| **`/lists`** | Reusable logic for paginated and searchable list views. |
| **`/forms`** | Standardized form containers and validation helpers. |

## Testing & Tooling

| Directory/File | Responsibility |
| :--- | :--- |
| **`/cypress`** | The complete E2E test suite. Includes commands for SQL injection and Slack simulation. |
| **`/styles`** | Global CSS and Material UI theme overrides. |
| **`/public`** | Static assets like the GoAlert logo and favicons. |
| **`esbuild.cypress.js`**| Build script for compiling Cypress tests for CI/CD environments. |