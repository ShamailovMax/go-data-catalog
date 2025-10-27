# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Common commands

- Build server binary
  ```bash path=null start=null
  go build -o bin/server ./cmd/server
  ```
- Run the server (expects environment set or a .env file at repo root)
  ```bash path=null start=null
  go run ./cmd/server
  ```
- Format and vet
  ```bash path=null start=null
  go fmt ./...
  go vet ./...
  ```
- Run tests (repository-wide)
  ```bash path=null start=null
  go test ./...
  ```
- Run a single test
  ```bash path=null start=null
  # by name (regex) within a package
  go test -run '^TestName$' ./internal/handlers -v
  ```
- Module housekeeping
  ```bash path=null start=null
  go mod tidy
  ```

## Environment and configuration

The app loads configuration from a .env file at the repository root (if present) and then from process environment variables. Required variables:

- DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
- SERVER_PORT

Example .env (values are placeholders):
```dotenv path=null start=null
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
SERVER_PORT=8080
```

Notes:
- Postgres connection uses pgxpool with sslmode=disable and UTF-8 client encoding.
- Missing .env is tolerated; variables must still be present in the environment.

## High-level architecture

- Entry point: HTTP server in `cmd/server/main.go` using Gin.
  - Loads config via `internal/config`.
  - Initializes Postgres connection via `internal/repository/postgres/db.go` (pgxpool) and injects it into repositories.
  - Registers routes and starts Gin server.
- Handlers: `internal/handlers/artifacts.go`
  - Thin layer over repositories, handles JSON binding/encoding and HTTP status codes.
- Repositories: `internal/repository/postgres`
  - `db.go`: wraps `pgxpool.Pool` and provides lifecycle management.
  - `artifacts.go`: SQL for listing/creating artifacts.
- Models: `internal/models/models.go`
  - Plain data structs for Contacts, Artifacts, and ArtifactFields (JSON tags align with API payloads).
- Config: `internal/config/config.go`
  - Loads `.env` (if present) via `godotenv`, then parses required vars with `caarlos0/env`.

### HTTP API surface (current)
- GET `/health` → `{ "status": "OK" }`
- Group `/artifacts`
  - GET ``/`` → list artifacts
  - POST ``/`` → create artifact from JSON body

### Request/response behavior
- Global middleware sets `Content-Type: application/json; charset=utf-8` on responses.
- Gin defaults (logger, recovery) are enabled via `gin.Default()`.

## Development tips specific to this repo

- Ensure the database specified by env vars is reachable before starting the server; schema must already exist (this repo does not include migrations).
- If connection fails at startup, check env var values and that the DB accepts non-SSL connections.
