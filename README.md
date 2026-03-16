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
│   │   │       ├── 000001_create_users_table.down.sql
│   │   │       └── 000001_create_users_table.up.sql
│   │   └── sqlite/
│   │       └── sqlite.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── internal/
│       ├── config/
│       │   └── routes.go
│       ├── handlers/
│       │   ├── handler.go
│       │   └── user/
│       │       └── user.go
│       ├── models/
│       │   └── user.go
│       ├── repository/
│       │   ├── repository.go
│       │   └── user/
│       │       ├── create.go
│       │       └── user.go
│       └── service/
│           ├── service.go
│           └── user/
│               └── user.go
├── frontend/
│   ├── .editorconfig
│   ├── .gitattributes
│   ├── .gitignore
│   ├── .oxlintrc.json
│   ├── .prettierrc.json
│   ├── Dockerfile
│   ├── eslint.config.js
│   ├── index.html
│   ├── jsconfig.json
│   ├── package-lock.json
│   ├── package.json
│   ├── public/
│   │   └── favicon.ico
│   ├── README.md
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.js
│   │   ├── router/
│   │   │   └── index.js
│   │   └── stores/
│   │       └── counter.js
│   └── vite.config.js
└── README.md
```
