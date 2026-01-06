# Konsultin Backend Boilerplate

This repository is the Konsultin backend boilerplate, curated by Kenly Krisaguino. It combines custom Konsultin libraries (routing, DTOs, error handling, logging) with a prewired Docker/dev setup so new services share the same conventions out of the box.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose (for local DB/profiles)
- `air` (hot reload) and `migrate` CLI are auto-installed by the Makefile when needed

## Quick Start

1. `make setup-project` to rename the module (prompts for project name) and seed `.env` if missing.
2. Review/edit `.env` (copy from `.env.example` if needed).
3. `make init` to tidy deps and install dev tooling.
4. `make dev` for hot-reload, or `make run` to start once.
5. `make up` to start the DB stack via Docker (profile picked by `DB_DRIVER`). Use `make down` to stop.

## Environment

Base config lives in `.env.example`; copy to `.env` and adjust. Key variables:

### Core Config
- `APP_ENV` — controls environment (`development`/`production`).
- `PORT` — API listen port; `DEBUG` toggles verbose error payloads.
- HTTP timeouts: `HTTP_READ_TIMEOUT_SECONDS`, `HTTP_WRITE_TIMEOUT_SECONDS`, `HTTP_IDLE_TIMEOUT_SECONDS`.
- Rate limiting: `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`.
- CORS: `CORS_ALLOW_ORIGINS`.

### JWT & Sessions
- `JWT_ISSUER` — JWT issuer name.
- `JWT_SECRET` — Secret key for signing JWT tokens.
- `USER_SESSION_LIFETIME` — Access token lifetime in seconds (default: 3600).
- `USER_SESSION_REFRESH_LIFETIME` — Refresh token lifetime in seconds (default: 2592000).

### OAuth Configuration
- `GOOGLE_CLIENT_ID` — Google OAuth Client ID for Sign in with Google.

### Database
- `DB_DRIVER` (`mysql`/`postgres`), `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME`.
- Connection pool: `DB_MAX_IDLE_CONN`, `DB_MAX_OPEN_CONN`, `DB_MAX_CONN_LIFETIME`.
- `DB_TIMEOUT_SECONDS`.

### NATS & Worker
- `NATS_URL` - NATS server URL (default: `nats://localhost:4222`).

### Docker
- `COMPOSE_PROJECT_NAME` — Docker compose project name (set during setup).

## Makefile Commands

- `make setup-project` — rename module to your project, update env defaults, install `air` if missing.
- `make init` — ensure `.env`, tidy modules, install `migrate` (MySQL/Postgres) and `air` if needed.
- `make dev` — start with hot reload via Air; writes temp files to `./tmp`.
- `make run` — run the API once with `go run ./app`.
- `make up` / `make down` — start/stop docker compose stack; profile derived from `DB_DRIVER` (postgres/mysql).
- `make bs` — alias to `make up` using the same profile logic.
- `make lint` — run `go vet ./...`.
- `make tidy` — run `go mod tidy`.
- `make db-up` / `make db-down` — run migrations up/down (one step) using `DB_*` connection info.
- `make db-script` — create a new timestamped SQL migration in `./migrations`.
- `make db-version` — move schema to a specific migration version.

## Authentication Flow

### 1. Anonymous Session (App Authentication)

First, the client must obtain an anonymous session token using Basic Auth:

```bash
POST /v1/users/anon/sessions
Authorization: Basic base64(clientId:clientSecret)
```

Response:
```json
{
  "data": {
    "session": { "token": "eyJhbGci...", "expiredAt": 1234567890 },
    "scopes": ["privilege1", "privilege2"]
  }
}
```

### 2. User Login with Password

Login using email/phone/username + password (requires anonymous session token):

```bash
POST /v1/users/sessions/login
Authorization: Bearer <anonymous_token>
Content-Type: application/json

{
  "identifier": "user@example.com",
  "password": "secret123",
  "device": {
    "deviceId": "device-uuid",
    "devicePlatformId": 1
  }
}
```

### 3. User Login with Google

Login using Google OAuth (requires anonymous session token):

```bash
POST /v1/users/sessions/google
Authorization: Bearer <anonymous_token>
Content-Type: application/json

{
  "provider": 2,
  "idToken": "google-id-token",
  "device": {
    "deviceId": "device-uuid",
    "devicePlatformId": 1
  }
}
```

### 4. Refresh Token

Refresh user session using refresh token:

```bash
PUT /v1/users/sessions
Authorization: Bearer <refresh_token>
Content-Type: application/json

{
  "refreshToken": "...",
  "device": { ... }
}
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Health check |
| `POST` | `/v1/cron/{cronType}` | Cron trigger |
| `POST` | `/v1/users/anon/sessions` | Create anonymous session (Basic Auth) |
| `PUT` | `/v1/users/sessions` | Refresh user session |
| `POST` | `/v1/users/sessions/login` | Login with password |
| `POST` | `/v1/users/sessions/google` | Login with Google OAuth |
| `POST` | `/v1/simulation` | Trigger Worker Simulation |

## Project Structure

```
├── app/                    # Application entry point
├── config/                 # Configuration loading
├── dto/                    # Data Transfer Objects
├── internal/
│   └── svc-core/
│       ├── model/          # Database models
│       ├── repository/     # Database operations
│       ├── service/        # Business logic
│       ├── sql/            # SQL prepared statements
│       └── pkg/
│           ├── httpk/      # HTTP utilities & middleware
│           └── oauth/      # OAuth providers
│               └── google/ # Google OAuth implementation
├── libs/                   # Shared libraries
│   ├── errk/               # Error handling
│   ├── logk/               # Logging
│   ├── natsk/              # NATS wrapper
│   ├── sqlk/               # Database utilities
│   └── timek/              # Time utilities
└── migrations/             # Database migrations
```

## Changes

> ### v1.2.0 - NATS Worker Integration
> - Add NATS support for background workers
> - Add `libs/natsk` wrapper
> - Implement `Repo -> Publish -> NATS -> Consume -> Worker` flow
> - Add Worker simulation endpoint

> ### v1.1.0 - Authentication System
> - Add flexible login (email/phone/username + password)
> - Add Google OAuth authentication
> - Add user credential management
> - Add anonymous session validation for login endpoints
> - Add JWT-based session management with access and refresh tokens
> - Add request binding and validation helpers (go-playground/validator)

> ### v1.0.0 - Initial Project
> - Create Project Whole Boilerplate Base