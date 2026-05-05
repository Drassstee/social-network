# Social Network

A real-time, single-page social network application featuring private messaging, community groups, event management, and live notifications.

## Backend Architecture

The backend is built with Go and strictly adheres to **Clean Architecture** principles to separate concerns, making the codebase highly maintainable and testable:

*   **`cmd/`**: Contains the application entry point (`server/main.go`), responsible for wiring up dependencies and starting the server.
*   **`db/`**: Manages the SQLite database, including connection logic and all `.sql` migration files.
*   **`internal/`**: The core application logic, cleanly divided into three distinct layers:
    *   **Handlers (HTTP Layer)**: Parses incoming requests, validates basic input, and returns JSON responses.
    *   **Services (Business Logic)**: Contains the core rules of the application (e.g., verifying permissions, orchestrating data).
    *   **Repositories (Data Access)**: The *only* layer allowed to execute SQL queries and interact directly with the database.

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
│   │   │   └── sqlite/
│   │   │       ├── 000001_create_users_table.up.sql
│   │   │       ├── 000002_create_groups_table.up.sql
│   │   │       ├── 000003_create_posts_table.up.sql
│   │   │       ├── 000004_create_sessions_table.up.sql
│   │   │       ├── 000005_create_followers_table.up.sql
│   │   │       ├── 000006_create_messages_table.up.sql
│   │   │       ├── 000007_create_group_members_table.up.sql
│   │   │       ├── 000008_create_group_events_table.up.sql
│   │   │       ├── 000009_create_group_messages_table.up.sql
│   │   │       └── 000010_create_notifications_table.up.sql
│   │   └── sqlite/
│   │       └── sqlite.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── internal/
│   │   ├── config/
│   │   │   └── routes.go
│   │   ├── handlers/
│   │   │   ├── chat/
│   │   │   ├── group/
│   │   │   ├── notifications/
│   │   │   ├── post/
│   │   │   ├── user/
│   │   │   └── handler.go
│   │   ├── models/
│   │   │   ├── chat.go
│   │   │   ├── group.go
│   │   │   ├── notification.go
│   │   │   ├── post.go
│   │   │   ├── profile.go
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── chat/
│   │   │   ├── group/
│   │   │   ├── notifications/
│   │   │   ├── post/
│   │   │   ├── user/
│   │   │   └── repository.go
│   │   ├── service/
│   │   │   ├── chat/
│   │   │   ├── group/
│   │   │   ├── notifications/
│   │   │   ├── post/
│   │   │   ├── user/
│   │   │   └── service.go
│   │   ├── sessions/
│   │   │   └── sessions.go
│   │   ├── utils/
│   │   │   └── utils.go
│   │   └── web/
│   │       └── web.go
├── frontend/ (Vue.js SPA)
└── README.md
```
## Commands

Start the application in the background:
```bash
docker compose up --build -d
```

Stop and remove containers:
```bash
docker compose down
```

Stop and remove containers and the database volume:
```bash
docker compose down -v
```
## Features

### Real-Time Chat
Full-featured private messaging system powered by WebSockets.
- **WebSocket Hub**: Manages real-time connections, online status tracking, and message broadcasting.
- **Online Status**: Real-time visibility of online users with optimized bulk fetching.
- **Image Sharing**: Support for image uploads in chat (stored locally and served via `/api/v1/uploads/`).
- **Performance**: High-performance broadcasting using in-memory group member tracking and username caching to minimize database load.

### Groups
Robust group management for community interaction.
- **Group Lifecycle**: Create groups with titles and descriptions. Creators are automatically granted administrative roles.
- **Membership Management**: 
    - **Invitations**: Members can invite other users to join.
    - **Join Requests**: Users can request to join groups, requiring approval by the group creator.
- **Events**: Create group events with RSVP support ("going", "not_going").
- **Group Chat**: Each group has its own dedicated real-time chat room.
- **Transactional Consistency**: All multi-step membership operations are protected by atomic database transactions (`WithTx`).

### Notifications
Real-time alerting system for cross-module interactions.
- **Real-Time Signaling**: Instant browser alerts for new group invitations, join requests, and approval outcomes.
- **Persistent Feed**: Notifications are stored in the database, ensuring users see them even after refreshing.
- **Metadata Integration**: Automatically captures actor usernames and target titles to provide rich, informative alerts (e.g., "John Doe invited you to 'The Coding Club'").
- **Extensible Design**: Decoupled signaling interface allowing for future expansion (e.g., push notifications, email).

### User Experience & Identity
- **Sitewide Avatar Rendering**: Persistent avatar uploads during registration that are displayed dynamically across the Activity Feed, Profile View, Comments, and Navbar.
- **Premium Feedback System**: Custom, globally styled toast notifications with slide-in animations (success/error variants) that replace native, disruptive browser `alert()` popups.
- **High-Quality Media**: Configured to support up to 10MB image uploads seamlessly via Nginx.

### Security & Performance
- **Session Security**: Session tokens are secured using `HttpOnly` and `SameSite: Strict` cookies, providing native protection against XSS and CSRF attacks without the need for complex middleware.
- **Performance Optimized**: Built with a lightweight Go backend, highly indexed SQLite database, and a fast, optimized Vue 3 SPA powered by Vite.
