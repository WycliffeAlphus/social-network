```md
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/              # Load env vars and configuration
│   ├── handler/             # API handlers 
│   ├── middleware/          # JWT, CORS, logging, etc.
│   ├── model/               # Domain models 
│   ├── repository/          # DB access 
│   ├── service/             # Business logic 
│   ├── routes/              # API routes
│   ├── utils/               # Helpers 
│   └── auth/                # JWT handling, token middleware, login logic
├── migrations/              # SQL migrations (use golang-migrate)
├── database/
│   └── connection.go        # DB initialization (PostgreSQL, SQLite, etc.)
├── go.mod
├── go.sum
└── README.md
```