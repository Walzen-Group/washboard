# Washboard Backend

Go-based REST API backend for Washboard, an internal server control panel for managing Docker containers and stacks via Portainer.

## Tech Stack

- **Go 1.22** with [Gin](https://github.com/gin-gonic/gin) web framework
- **MongoDB** for persistence
- **JWT** authentication (`gin-jwt`)
- **WebSocket** support for real-time updates (`gorilla/websocket`)
- **In-memory caching** with fallback cache for resilience (`go-cache`)

## Project Structure

```
backend/
├── app.go                 # Entry point, router setup, JWT config
├── api/                   # HTTP handlers
│   ├── db.go              # Stack settings CRUD endpoints
│   ├── docker-update-manager.go  # Container/image status endpoints
│   ├── stack-manager.go   # Stack control endpoints (start/stop/action)
│   └── websocket.go       # WebSocket handler for real-time updates
├── portainer/             # Portainer API integration layer
├── db/                    # MongoDB operations
├── state/                 # App configuration (YAML + env var overrides)
├── auth/                  # JWT authentication
├── control/               # Business logic (auto-start sync, stop-all)
├── types/                 # Data structures & constants
├── helper/                # Utility functions
└── werrors/               # Custom error types
```

## Configuration

Configuration is loaded from `secrets.yaml` in the working directory. Environment variables override YAML values.

| Variable | Description | Required |
|----------|-------------|----------|
| `PORTAINER_SECRET` | Portainer API key | Yes |
| `PORTAINER_URL` | Portainer base URL (e.g. `http://portainer:9000/api`) | Yes |
| `DB_URL` | MongoDB connection string | Yes |
| `USER` | Login username | Yes |
| `PASSWORD` | Login password (bcrypt hashed) | Yes |
| `JWT_SECRET` | JWT signing key (auto-generated if not set) | No |
| `CACHE_DURATION_MINUTES` | Image status cache TTL (default: `1`) | No |
| `CORS` | Comma-separated allowed origins | No |
| `START_STACKS_ON_LAUNCH` | Auto-start stacks on app launch (default: `false`) | No |
| `START_ENDPOINT_ID` | Default Portainer endpoint ID (default: `1`) | No |

## Running Locally

```bash
# Install dependencies
go mod download

# Run the application
go run .

# Or use Air for live reload
air
```

The server starts on port **8080**.

## Running with Docker

```bash
docker-compose -f docker-compose-dev.yml up
```

The backend is built as a multi-stage Docker image (`golang:1.22-alpine` -> `alpine:latest`) and runs as a non-root user.

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/login` | Login (returns JWT) |
| POST | `/api/auth/logout` | Logout |
| POST | `/api/auth/refresh_token` | Refresh JWT token |

### Portainer (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/portainer/endpoint` | Get endpoint ID by name |
| GET | `/api/portainer/stacks` | Get all stacks and containers |
| GET | `/api/portainer/containers` | Get containers for a stack |
| GET | `/api/portainer/image-status` | Get image update status |
| POST | `/api/portainer/update-container` | Pull new image for container |
| POST | `/api/portainer/stacks/:id/start` | Start a stack |
| POST | `/api/portainer/stacks/:id/stop` | Stop a stack |
| PUT | `/api/portainer/stacks/:id/update` | Update stack configuration |
| POST | `/api/portainer/containers/:containerId/:action` | Container action (start/stop/restart/kill/pause/resume) |

### Stack Settings (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/db/stacks` | Create stack settings |
| GET | `/api/db/stacks` | Get all stack settings |
| GET | `/api/db/stacks/:name` | Get stack settings by name |
| PUT | `/api/db/stacks/:name` | Update stack settings |
| DELETE | `/api/db/stacks/:name` | Delete stack settings |
| POST | `/api/db/sync` | Sync Portainer stacks with database |

### Control (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/control/sync-autostart` | Start stacks marked with autoStart |
| POST | `/api/control/stop-all` | Stop all stacks (by priority order) |

### WebSocket (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/ws/stacks-update` | Real-time stack update status stream |

## Key Features

- **Image update detection** with configurable caching and background refresh (every 24h)
- **Priority-based orchestration** for startup and shutdown sequences
- **Auto-start sync** to restore stacks after restarts
- **Self-preservation** — skips stopping stacks containing washboard images
- **Fallback cache** that persists across Portainer API failures
- **Structured logging** with file rotation (10 MB max per file)

## Database

- **Database:** `washb` (MongoDB)
- **Collections:** `stack_settings`, `group_settings`, `accounts`

The `stack_settings` collection stores stack metadata with fields: `stackName`, `stackId`, `priority`, and `autoStart`.
