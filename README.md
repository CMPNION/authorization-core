# Auth Service (Go)

Production-ready authentication service built with Go, following DDD and Clean Architecture principles.  
It supports user registration, login, and a protected `GetMe` endpoint using JWT, with Postgres for persistence and Redis for token storage. The whole stack runs via Docker Compose.

---

## GitHub Short Description (EN)

A production-ready Go authentication service with JWT, Postgres, Redis, Docker Compose, and Clean Architecture.

---

## Features

- User registration and login
- JWT access tokens
- Protected `GET /me` endpoint
- Postgres storage for users
- Redis storage for issued tokens
- Clean Architecture / DDD-friendly structure
- Docker Compose for local development

---

## Architecture (DDD + Clean)

- **Domain**: core entities and errors
- **Usecases**: application logic (auth service)
- **Repositories**: infrastructure adapters (Postgres, Redis, JWT, bcrypt)
- **Delivery**: HTTP handlers and middleware
- **App wiring**: dependency assembly and server setup

---

## Tech Stack

- Go
- Postgres
- Redis
- JWT
- Docker / Docker Compose

---

## Project Structure

- `cmd/app` — entry point
- `internal/app` — dependency wiring and server setup
- `internal/domain` — entities and errors
- `internal/usecase` — business logic
- `internal/repository` — persistence and adapters
- `internal/delivery/httpdelivery` — HTTP handlers, middleware, routing
- `deployments/postgres` — init SQL

---

## Getting Started

### 1) Configure environment

Create `.env` in the project root. Example:

    APP_PORT=4000

    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=authdb
    DB_SSLMODE=disable

    REDIS_ADDR=redis:6379
    REDIS_PASSWORD=
    REDIS_DB=0

    JWT_SECRET=change-me
    TOKEN_TTL=15m

### 2) Run with Docker Compose

    docker compose up -d --build

Only the application port is exposed outside the Docker network. Postgres and Redis are internal.

---

## API

### POST /signup

Request body:

    {
      "login": "user",
      "password": "pass"
    }

Response:

    { "token": "..." }

---

### POST /login

Request body:

    {
      "login": "user",
      "password": "pass"
    }

Response:

    { "token": "..." }

---

### GET /me

Headers:

    Authorization: Bearer <token>

Response:

    {
      "id": 1,
      "login": "user"
    }

---

## Notes

- JWT is signed with `JWT_SECRET`.
- Tokens are stored in Redis to allow future revoke/blacklist logic.
- Passwords are hashed with bcrypt.
- Database schema is initialized via `deployments/postgres/init.sql`.

---

## Development

If you run locally without Docker, ensure Postgres and Redis are accessible and update `.env` accordingly.

---

## Roadmap Ideas

- Refresh tokens
- Token revocation and session management
- Rate limiting
- Observability (metrics, tracing)
- Integration and contract tests