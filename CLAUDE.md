# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Washboard is an internal server control panel that integrates with PortainerEE to manage container stacks and services. It provides a web UI for updating containers, managing stack lifecycle (start/stop/update), and monitoring image update status.

## Tech Stack

- **Backend:** Go 1.22 with Gin-gonic, MongoDB, Gorilla WebSocket
- **Frontend:** Vue 3 + Vuetify 3 + TypeScript, built with Vite 5, state managed by Pinia
- **Package manager:** pnpm (frontend)
- **Deployment:** Docker multi-stage builds, nginx for frontend serving, GitHub Actions CI/CD pushing to ghcr.io

## Build & Run Commands

### Backend
```bash
cd backend
go mod download
go build -o washboard-backend
```

### Frontend
```bash
cd frontend
pnpm install
pnpm dev          # Dev server with hot reload (--host enabled)
pnpm build        # Production build (runs vue-tsc --noEmit first)
pnpm lint         # ESLint with auto-fix
```

### Docker (full stack)
```bash
docker compose -f docker-compose-dev.yml up --build
```
Frontend is exposed on port 10004. Backend port is configured via environment variables in `.env`.

## Architecture

### Backend (`backend/`)

Entry point is `app.go` which sets up Gin router, JWT auth middleware, CORS, and all API routes. Key packages:

- **`api/`** - HTTP handlers for all REST endpoints and WebSocket
- **`auth/`** - JWT authentication (7-day token, 30-day refresh, cookie-based)
- **`portainer/`** - PortainerEE API client integration, background update checker
- **`control/`** - Service control logic (autostart sync, stop-all)
- **`db/`** - MongoDB operations for stack settings persistence
- **`state/`** - Singleton app config/state (loaded from YAML config)
- **`types/`** - Shared type definitions

### API Route Structure

All routes are under `/api`:
- `/api/auth/*` - Login, logout, token refresh (unauthenticated)
- `/api/portainer/*` - Endpoints, containers, image status, stacks CRUD (authenticated)
- `/api/db/stacks/*` - Stack settings CRUD in MongoDB (authenticated)
- `/api/ws/stacks-update` - WebSocket for real-time stack update progress
- `/api/control/*` - Sync autostart, stop-all operations (authenticated)

### Frontend (`frontend/`)

Vue 3 SPA using file-based routing (`unplugin-vue-router`), auto-imported components (`unplugin-vue-components`), and Vuetify component library. Uses `unplugin-auto-import` for Vue/router APIs. Layout system via `vite-plugin-vue-layouts`.

### Startup Behavior

On launch, the backend syncs with Portainer and optionally starts stacks if `StartStacksOnLaunch` is configured. A background update checker runs continuously to detect available image updates.

## Configuration

App config is loaded via the `state` package from a YAML file. Key settings include:
- PortainerEE connection details
- JWT secret (auto-generated if not set, but sessions are lost on restart)
- CORS allowed origins
- MongoDB connection
- `StartEndpointId`, `StartStacksOnLaunch`

Environment variables are loaded from `.env` for Docker deployments.
