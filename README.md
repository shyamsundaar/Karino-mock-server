# Karino Mock Server

> A small Go-based mock server for Karino — provides example endpoints and data for development and testing.

## Features
- Simple REST mock implemented in Go
- Swagger/OpenAPI docs included (`docs/swagger.yaml`, `docs/swagger.json`)
- Docker Compose for easy local environment

## Prerequisites
- Go 1.20+ installed
- Docker & Docker Compose (for containerized run)

## Quick Start

1. Start services with Docker Compose:

```bash
make dev
# or
docker-compose up -d
```

2. Run the server locally (development):

```bash
# install live-reload tool if needed
make install-modules

# then start with air (if installed)
make start-server

# or run directly with Go
go run main.go
```

4. Stop Docker Compose when finished:

```bash
make dev-down
```

## Running the generator
There is a generator command under `cmd/generate`. To run it:

```bash
go run ./cmd/generate
# or
go run ./cmd/generate/main.go
```

## API Documentation
Open the included Swagger files to view the API surface:

- `docs/swagger.yaml`
- `docs/swagger.json`

You can load `docs/swagger.yaml` into Swagger UI or similar tools.

## Project Layout
- `main.go` — application entrypoint
- `controllers/` — HTTP route handlers (e.g. `farmers.controller.go`)
- `models/` — data models (e.g. `farmers.model.go`)
- `initializers/` — initialization code (DB, env)
- `cmd/` — subcommands and tools (codegen under `cmd/generate`)
- `docs/` — swagger definitions
- `docker-compose.yml` — local stack
- `app.env` — example environment variables

## Notes
- The project is intentionally lightweight and intended as a mock server for local development and testing. Adjust environment variables in `app.env` as needed.
- If you prefer not to use Docker, ensure any required services (DB) are running locally and configured via environment variables.

## Next steps
- Run the server and open the Swagger file in Swagger UI to explore endpoints.
- If you'd like, I can add example curl commands or a Postman collection.

---
Generated on 2025-12-29
# Debugger

```sh
go install github.com/go-delve/delve/cmd/dlv@latest
```