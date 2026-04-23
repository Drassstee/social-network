Structure in backend: database(db folder), server(cmd/server/main.go), app(internal folder)
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
│   │   │       ├── 000002_create_posts_table.up.sql
│   │   │       ├── 000003_create_messages_table.up.sql
│   │   │       ├── 000004_create_groups_tables.up.sql
│   │   │       ├── 000005_create_group_events_table.up.sql
│   │   │       ├── 000006_create_group_messages_table.up.sql
│   │   │       └── 000007_create_notifications_table.up.sql
│   │   └── sqlite/
│   │       └── sqlite.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── internal/
│   │   ├── config/
│   │   │   └── routes.go
│   │   ├── handlers/
│   │   │   ├── chat/
│   │   │   │   └── chat_handler.go
│   │   │   ├── group/
│   │   │   │   └── group_handler.go
│   │   │   ├── notifications/
│   │   │   │   └── notifications_handler.go
│   │   │   └── handler.go
│   │   ├── models/
│   │   │   ├── chat.go
│   │   │   ├── group.go
│   │   │   ├── notification.go
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── chat/
│   │   │   │   └── chat_repository.go
│   │   │   ├── group/
│   │   │   │   └── group_repository.go
│   │   │   ├── notifications/
│   │   │   │   └── notifications_repository.go
│   │   │   └── repository.go
│   │   ├── service/
│   │   │   ├── chat/
│   │   │   │   └── hub.go
│   │   │   ├── group/
│   │   │   │   └── group_service.go
│   │   │   ├── notifications/
│   │   │   │   └── notifications.go
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
```bash
 docker compose up --build -d
```
Stop and remove containers:
```bash
 docker compose down
```
stop and remove containers and database:
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
- **Clean Architecture**: Decoupled signaling interface allowing for future expansion (e.g., push notifications, email).
