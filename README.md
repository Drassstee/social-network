# Social Network // Neural_Node_Link

A premium, real-time social network featuring a synthwave-inspired SPA and a robust Go backend. Built with Clean Architecture principles across the entire stack.

## Architecture

### Backend (Go + SQLite)
The backend adheres to **Clean Architecture** to ensure separation of concerns and testability:
- **Domain Layer (`models/`)**: Pure business entities and interface definitions.
- **Application Layer (`service/`)**: Orchestrates business logic and handles cross-cutting concerns.
- **Infrastructure Layer (`repository/`)**: Direct data persistence logic.
- **Presentation Layer (`handlers/`)**: HTTP and WebSocket entry points.

### Frontend (Vue 3 + Pinia + Vite)
The frontend is refactored for high performance and maintainability:
- **Store-Driven Orchestration**: All state and API logic is centralized in **Pinia** stores (`postStore`, `chatStore`, `uiStore`), decoupling the view layer from business logic.
- **Component-Driven Design**: Reusable atomic components (`UserAvatar`, `PostCard`, `SkeletonLoader`) ensure UI consistency and DRY code.
- **Responsive Navigation**: Features a hybrid sidebar/hamburger-menu system optimized for both desktop and mobile users.

## Project Structure

```
social-network/
├── backend/
│   ├── cmd/server/       <- Entry point
│   ├── internal/
│   │   ├── handlers/     <- Presentation
│   │   ├── models/       <- Domain Interfaces
│   │   ├── repository/   <- Persistence
│   │   └── service/      <- Business Logic
├── frontend/
│   ├── src/
│   │   ├── components/   <- Atomic UI units
│   │   ├── stores/       <- Pinia State (Source of Truth)
│   │   ├── views/        <- Layout containers
│   │   └── api/          <- Fetch wrappers
└── docker-compose.yml
```

## Core Features

### Real-Time Hub
- **Unified WebSocket Protocol**: Handles private messages, group chats, and global notifications.
- **Live Sync**: Instant sitewide status updates and unread message counters.
- **Global Toasts**: Real-time notification system that alerts users to events even when away from the chat view.

### Syndicate (Groups)
- **Moderated Entry**: Creator-approval workflow for join requests and invitations.
- **Group Intelligence**: Dedicated chat rooms and event scheduling with RSVP tracking.

### User Experience
- **Retro-Future Aesthetic**: A custom-built CSS design system with glassmorphism and neon accents.
- **Skeleton Loaders**: Premium "shimmer" loading states to prevent layout shifts.
- **Secure Auth**: Cookie-based session management with `SameSite=Strict` and bcrypt encryption.

## Getting Started

### Start the Containers
```bash
docker compose up --build -d
```

### Stop the Containers
```bash
docker compose down
```

### Development
1. **Backend**: `cd backend && go run cmd/server/main.go`
2. **Frontend**: `cd frontend && npm install && npm run dev`

## Architectural Standards
- **Dependency Inversion**: Services depend on abstractions, never on concrete repository implementations.
- **Store Orchestration**: Views never call APIs directly; they only interact with Pinia stores.
- **Unified Type Safety**: Standardized `int64` identifiers and snake_case JSON mapping across the full stack.
