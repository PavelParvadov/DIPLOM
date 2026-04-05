# HappyHouse Architecture

HappyHouse is split into two independent applications:

- `backend/` contains the Go API built in clean architecture style
- `frontend/` contains the React + TypeScript SPA

## Backend layers

- `internal/domain` describes domain entities and repository contracts
- `internal/usecase` contains business rules
- `internal/repository/postgres` implements PostgreSQL data access
- `internal/transport/http` contains routing, middleware, handlers and DTOs

## Data flow

1. HTTP handler accepts the request and validates transport input.
2. Use case checks house membership, roles and business rules.
3. Repository reads or writes data in PostgreSQL.
4. Handler returns a normalized JSON response.

## Frontend structure

- `src/app` for router, layout and providers
- `src/pages` for screens
- `src/features` for user scenarios and forms
- `src/entities` for typed domain entities
- `src/shared` for API client, UI primitives and utilities
