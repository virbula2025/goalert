# GoAlert GraphQL API (`/graphql2`)

This directory contains the primary API layer for GoAlert. It uses a "Schema-First" approach powered by `gqlgen`, providing a strongly-typed bridge between the frontend and the backend business logic.

## Overview

The API is structured to decouple the GraphQL schema definitions from the actual data-fetching logic. It leverages Go's type system to ensure that API requests are validated and mapped correctly to the internal domain models (like `alert.Alert` or `service.Service`).

## Core Architecture

1.  **Code Generation**: Most of the boilerplate (parsing, marshalling, and interface definitions) is generated automatically from the schema.
2.  **Resolvers**: Custom logic for fetching data is implemented in "Resolvers." If a field can be mapped directly to a Go struct field, `gqlgen` handles it; otherwise, a resolver method is used.
3.  **Application Integration**: The `graphqlapp` sub-package bridges the generated schema with the application's stores and services.

---

## File Breakdown

### Boilerplate & Generation
* **`gen.go`**: Contains the `//go:generate` directives. It coordinates the deletion of old generated files and the execution of `gqlgen` and other custom tools like `configparams` and `limitapigen` to rebuild the API.
* **`generated.go`**: The auto-generated heart of the API. It contains the executable schema logic, complexity analysis configurations, and the `ResolverRoot` interface which the application must implement.
* **`models_gen.go`**: Contains Go struct definitions for GraphQL input types and enums (e.g., `CreateAlertInput`, `AlertStatus`) generated from the schema.

### Custom Scalars & Mapping
* **`isotimestamp.go`**: Implements custom marshalling and unmarshalling for the `ISOTimestamp` scalar. This ensures that time is consistently handled in `RFC3339Nano` format across the API.
* **`mapconfig.go` & `maplimit.go`**: Generated utility files that map internal system configurations and limits into a format accessible by the GraphQL API.

### Implementation (Resolvers)
* **`resolver.go`**: Defines the base `Resolver` struct, which holds references to the application's backend "Stores" (e.g., `AlertStore`, `UserStore`).
* **`graphqlapp/alert.go`**: Implements the `AlertResolver`. It handles complex field resolution for alerts, such as converting internal status codes into GraphQL-friendly enums and calculating alert metrics.
* **`graphqlapp/service.go`**: Contains logic for resolving Service-related queries, including their associated escalation policies, integration keys, and maintenance status.

---

## How it Works

### 1. Request Flow
When a query arrives, `generated.go` parses the request. It uses the **Complexity Root** to calculate if the query is too expensive to execute. If valid, it calls the methods defined in the `ResolverRoot` implemented by the application logic.

### 2. Data Mapping
The API frequently maps internal database models to external API models. For example, in `AlertResolver`, numeric database IDs are converted into strings for GraphQL compatibility, and internal state constants are mapped to readable enums.

### 3. Validation
Input validation is performed within the resolver logic before reaching the database layer. It uses internal validation packages to ensure data integrity and returns specific error extensions for the frontend to handle.

### 4. Concurrency
The generated code handles concurrent field resolution. When fetching lists, it uses `sync.WaitGroup` to resolve multiple items in parallel where possible, significantly improving performance for complex dashboard views that aggregate data from multiple stores.