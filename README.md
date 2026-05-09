# Social Network

A real-time, single-page social network application featuring private messaging, community groups, event management, and live notifications.

## Backend Architecture

The backend is built with Go and strictly adheres to **Clean Architecture** principles to separate concerns, making the codebase highly maintainable and testable:

*   **`cmd/`**: Contains the application entry point (`server/main.go`), responsible for wiring up dependencies and starting the server.
*   **`db/`**: Manages the SQLite database, including connection logic and all `.sql` migration files.
*   **`internal/`**: The core application logic, cleanly divided into layers:
    *   **Models (Domain Layer)**: Defines entities, interfaces, and core domain errors. This is the heart of the system.
    *   **Services (Application Layer)**: Implements business logic and orchestrates data flow between models and repositories.
    *   **Repositories (Infrastructure Layer)**: Handles persistence logic and SQL execution.
    *   **Handlers (Presentation Layer)**: Manages HTTP request/response cycles and interacts with Services.

## Project Structure

```
social-network/
├── .gitignore
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── db/
│   │   ├── migrations/
│   │   └── sqlite/
│   ├── internal/
│   │   ├── config/           <- Routing and app configuration
│   │   ├── handlers/         <- HTTP Presentation layer
│   │   ├── models/           <- Domain Layer (Entities & Interfaces)
│   │   ├── repository/       <- Infrastructure Layer (Persistence)
│   │   ├── service/          <- Application Layer (Business Logic)
│   │   ├── utils/            <- Common helpers (uploader, auth, etc.)
│   │   └── web/              <- Web utilities (JSON parsing, error mapping)
│   ├── Dockerfile
│   └── go.mod
├── frontend/                 <- Vue.js 3 SPA (Vite)
└── README.md
```

## Architectural Standards

The project maintains high code quality through established patterns:
- **Unified ID Type**: All entity identifiers are standardized as `int64` across all layers for consistency and scalability.
- **Dependency Inversion**: High-level modules (Services) do not depend on low-level modules (Repositories); both depend on abstractions (Interfaces in Models).
- **Clean Communication**: Real-time features are decoupled from business logic using a WebSocket Hub and a centralized Notification system.

## Commands

Start the application (Full Stack):
```bash
docker compose up --build -d
```

Stop and remove containers:
```bash
docker compose down
```

## Core Features

### Real-Time Communication
- **WebSocket Hub**: Manages live connections and online status tracking.
- **Private Messaging**: 1:1 chat with real-time delivery and image sharing support.
- **Live Notifications**: Instant alerts for group invites, join requests, and follow activity.

### Community Groups
- **Group Management**: Creation, membership lifecycle, and dedicated group chat rooms.
- **Join Workflow**: Support for both direct invitations and moderated join requests (creator approval).
- **Events**: Event scheduling with RSVP tracking ("going", "not_going").

### User Experience
- **Profile Management**: Customizable profiles with sitewide avatar rendering.
- **Responsive UI**: A modern Vue 3 interface with smooth transitions and premium toast notifications.
- **Security**: HttpOnly/SameSite session management and encrypted password storage (bcrypt).
