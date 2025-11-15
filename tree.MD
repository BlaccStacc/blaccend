.
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── router.go
│   │   ├── auth_handlers.go
│   │   ├── user_handlers.go
│   │   ├── twofa_handlers.go
│   │   ├── health_handlers.go
│   │   └── middleware.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── db/
│   │   ├── db.go
│   │   └── init.sql            # (or keep in /migrations, see below)
│   │
│   ├── mail/
│   │   └── mailer.go
│   │
│   ├── security/
│   │   ├── password.go
│   │   ├── token.go
│   │   └── jwt.go
│   │
│   └── models/
│       └── user.go
│
├── migrations/
│   ├── 001_init.sql
│   ├── 002_add_auth_fields.sql
│   └── 003_add_2fa_fields.sql
│
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── deb_setup.sh
├── arch_setup.sh
├── run.sh
├── README.md
└── cheat.MD
