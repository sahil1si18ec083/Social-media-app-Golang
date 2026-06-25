# Social Media API

A Go backend for a social media application with authentication, email activation, JWT-based protected routes, Redis caching, rate limiting, and Swagger documentation.

## Features

- User registration with activation flow
- Login with JWT token generation
- Protected user and post routes
- Follow and unfollow users
- Create, update, delete, and fetch posts
- Create comments on posts
- Redis-based user caching
- Fixed-window rate limiting
- Swagger API docs
- PostgreSQL persistence

## Tech Stack

- Go
- Chi router
- PostgreSQL
- Redis
- Docker Compose
- Swagger

## Project Structure

```text
cmd/api/                 HTTP server, handlers, middleware
cmd/migrate/migrations/  SQL migrations
internal/auth/           JWT auth logic
internal/db/             DB connection setup
internal/env/            Environment variable helpers
internal/mailer/         Mail providers and templates
internal/ratelimiter/    Rate limiter logic
internal/store/          Database store layer
internal/store/cache/    Redis cache layer
docs/                    Swagger files
```

## Prerequisites

- Go installed
- Docker and Docker Compose installed
- `migrate` CLI installed for running migrations
- `air` installed if you want live reload during development

## Environment Variables

Create a `.env` file in the project root.

Example:

```env
ADDR=:8080
FRONTEND_URL=http://localhost:5173

DB_ADDR=postgres://admin:adminpassword@localhost:5433/socialnetwork?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

MIGRATIONS_PATH=cmd/migrate/migrations

JWT_SECRET=change-this-secret

ENV=development

MAILTRAP_HOST=live.smtp.mailtrap.io
MAILTRAP_PORT=587
MAILTRAP_USERNAME=your-mailtrap-username
MAILTRAP_PASSWORD=your-mailtrap-password
FROM_EMAIL_MT=your-email@example.com

SENDGRID_API_KEY=your-sendgrid-key
FROM_EMAIL_SG=your-email@example.com

REDIS_ADDR=localhost:6379
REDIS_PW=
REDIS_DB=0
REDIS_ENABLED=true

RequestsPerTimeFrame=7
RateLimiterEnabled=true
```

## Run Infrastructure

Start PostgreSQL, Redis, and Redis Commander:

```bash
docker compose up -d
```

Services:

- PostgreSQL: `localhost:5433`
- Redis: `localhost:6379`
- Redis Commander: `http://localhost:8081`

## Run Migrations

Create a new migration:

```bash
make migrate-create name=add_roles_table
```

Run migrations up:

```bash
make migrate-up
```

Rollback the latest migration:

```bash
make migrate-down
```

Force a migration version:

```bash
make migrate-force version=20260429211653
```

## Run the API

Normal run:

```bash
go run ./cmd/api
```

Using `air`:

```bash
air
```

Debug mode with `air`:

```bash
air -d
```

## Swagger

Swagger UI:

```text
http://localhost:8080/v1/swagger/index.html
```

Swagger JSON:

```text
http://localhost:8080/v1/swagger/doc.json
```

Regenerate Swagger docs:

```bash
swag init -g cmd/api/main.go
```

## Common Development Commands

Enter PostgreSQL inside Docker:

```bash
docker exec -it postgres-db psql -U admin -d socialnetwork
```

General PostgreSQL pattern:

```bash
docker exec -it <container-name> psql -U <db-user> -d <db-name>
```

Enter Redis CLI inside Docker:

```bash
docker exec -it redis-cache redis-cli
```

Check Redis keys:

```bash
KEYS *
```

Read a Redis value:

```bash
GET user:92
```

## Main API Routes

Authentication:

- `POST /v1/authentication/user`
- `POST /v1/authentication/login`

Users:

- `PUT /v1/users/activate/{token}`
- `GET /v1/users/feed`
- `GET /v1/users/{userID}`
- `PUT /v1/users/{userID}/follow`
- `PUT /v1/users/{userID}/unfollow`

Posts:

- `POST /v1/posts`
- `GET /v1/posts/{postID}`
- `PATCH /v1/posts/{postID}`
- `DELETE /v1/posts/{postID}`
- `POST /v1/posts/{postID}/comments`

System:

- `GET /v1/health`

## Notes

- Most application routes use Bearer token authentication.
- User activation must happen before login succeeds.
- Redis is optional and controlled by `REDIS_ENABLED`.

