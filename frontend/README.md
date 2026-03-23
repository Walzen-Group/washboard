# Washboard Frontend

The frontend for **Washboard**, a Docker stack management dashboard by Walzen-Group. It provides a web interface for managing and updating container stacks via Portainer.

## Tech Stack

- **Vue 3** with **TypeScript**
- **Vuetify 3** (Material Design UI)
- **Pinia** for state management
- **Vite** for builds and HMR
- **Axios** for API communication
- **WebSockets** for real-time updates

## Features

- **Docker Manager** — View, start, stop, restart, and reorder stacks with drag-and-drop
- **Docker Update Manager** — Queue and apply image updates across stacks
- **Dashboard** — Overview of available updates and port allocation
- **Authentication** — Login with automatic token refresh
- **Dark Mode** — Automatic system theme detection

## Getting Started

### Install dependencies

```bash
pnpm install
```

### Start the dev server

```bash
pnpm dev
```

The app will be available at [http://localhost:3000](http://localhost:3000). In development, the API proxy points to `http://localhost:8080`.

### Build for production

```bash
pnpm build
```

### Lint

```bash
pnpm lint
```

## Project Structure

```
src/
├── api/          # API helpers (Axios, WebSocket)
├── components/   # Reusable Vue components
├── layouts/      # Page layouts (default, login)
├── pages/        # File-based routing (index, login, docker-manager, etc.)
├── plugins/      # Vuetify, router, Pinia setup
├── stores/       # Pinia state stores
└── types/        # TypeScript interfaces and enums
```

Routes are auto-generated from the `pages/` directory via `unplugin-vue-router`.
