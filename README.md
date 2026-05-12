# Go Social

Backend API for a simple social network, focused on user registration, authentication, posts, comments, followers, and a personalized feed.

The project is built with Go, PostgreSQL, Chi, JWT authentication, Swagger documentation, optional Redis caching, database migrations, live reload with Air, and account activation by email.

> Status: this project is under development. Some features are already implemented, but there are known limitations documented at the end of this README.

This repository currently focuses on the backend API. A dedicated frontend application is planned as a separate future project.

## Table of Contents

- [Tech Stack](#tech-stack)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Development Options](#development-options)
- [Development with Dev Container](#development-with-dev-container)
- [Environment Variables](#environment-variables)
- [Running Without the Dev Container](#running-without-the-dev-container)
- [Database, Migrations and Seed](#database-migrations-and-seed)
- [Swagger Documentation](#swagger-documentation)
- [Authentication](#authentication)
- [Main Routes](#main-routes)
- [Response Format](#response-format)
- [Data Model](#data-model)
- [Useful Commands](#useful-commands)
- [Ports and Services](#ports-and-services)
- [Known Limitations](#known-limitations)
- [Suggested Roadmap](#suggested-roadmap)

## Tech Stack

- Go `1.25.x`
  - `go.mod` declares Go `1.25.3`
  - the Dev Container image uses `golang:1.25.10-bookworm`
- PostgreSQL `16`
- Redis `7`
- Chi router
- Chi CORS middleware
- JWT with `github.com/golang-jwt/jwt/v5`
- Bcrypt via `golang.org/x/crypto`
- Request validation with `go-playground/validator`
- Logging with `go.uber.org/zap`
- Swagger/OpenAPI with `swaggo`
- Database migrations with `golang-migrate`
- Live reload with `air`
- Email delivery through Mailtrap, with structure also prepared for SendGrid
- Docker Compose
- Dev Container support for reproducible development environments

## Features

- User registration with hashed passwords.
- JWT token generation.
- Account activation using an email token.
- Health check protected with Basic Auth.
- Create, read, update, and delete posts.
- Create comments on posts.
- Follow and unfollow users.
- Personalized feed with pagination, sorting, tag filtering, and text search.
- Basic role-based permission model: `user`, `moderator`, and `admin`.
- Optional Redis cache for authenticated user lookups.
- Swagger UI for API inspection.
- Dev Container setup for development across different machines.
- Live reload workflow through Air.

## Architecture

The backend follows a layered structure:

```text
HTTP request
   ↓
cmd/api handlers
   ↓
middlewares, validation and request parsing
   ↓
internal/store
   ↓
PostgreSQL
```

Main components:

- `cmd/api`: HTTP server, routes, handlers, middlewares, authentication, and JSON responses.
- `internal/store`: data access layer.
- `internal/store/cache`: Redis cache layer.
- `internal/db`: database connection and seed logic.
- `internal/auth`: JWT generation and validation.
- `internal/mailer`: email delivery.
- `cmd/migrate/migrations`: SQL migration files.
- `docs`: generated Swagger documentation.
- `.devcontainer`: containerized development environment.
- `DEVCONTAINER.md`: complete guide for using the project with Dev Containers, VS Code, and Zed.

## Project Structure

```text
go-social/
├── .devcontainer/
│   ├── Dockerfile
│   ├── devcontainer.json
│   ├── docker-compose.devcontainer.yml
│   └── post-create.sh
├── cmd/
│   ├── api/
│   │   ├── api.go
│   │   ├── api_test.go
│   │   ├── auth.go
│   │   ├── errors.go
│   │   ├── feed.go
│   │   ├── health.go
│   │   ├── json.go
│   │   ├── main.go
│   │   ├── middleware.go
│   │   ├── posts.go
│   │   ├── test_utils.go
│   │   └── users.go
│   │   ├── users_test.go.go
│   └── migrate/
│       ├── migrations/
│       └── seed/
├── docs/
├── internal/
│   ├── auth/
│   ├── db/
│   ├── env/
│   ├── mailer/
│   └── store/
├── .air.toml
├── .env.example
├── DEVCONTAINER.md
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Development Options

There are two supported development workflows.

### Recommended workflow

Use the Dev Container.

This gives every developer the same environment with Go, Air, Swag, Migrate, PostgreSQL client tools, Redis access, module caches, and build caches already configured.

### Alternative workflow

Run everything directly on the host machine.

This requires installing Go, `air`, `swag`, `migrate`, Make, Docker, and all required tooling manually.

## Development with Dev Container

This repository includes a Dev Container setup for a reproducible development environment across different machines and editors.

The Dev Container configuration is located under:

```text
.devcontainer/
├── Dockerfile
├── devcontainer.json
├── docker-compose.devcontainer.yml
└── post-create.sh
```

It provides the backend development environment with Go, Air, Swag, Migrate, PostgreSQL, Redis, Redis Commander, Go module cache, and Go build cache.

For the complete setup and usage guide, including instructions for **Zed Editor**, **VS Code**, environment variables, migrations, seed, Swagger generation, Git workflow, ports, and troubleshooting, see:

 [DEVCONTAINER.md]( DEVCONTAINER.md)

## Environment Variables

Create a `.env` file from `.env.example`:

```bash
cp .env.example .env
```

Current `.env.example` is optimized for the Dev Container workflow:

```env
ADDR=:8080
ENV=development
EXTERNAL_URL=http://localhost:8080

DB_USER=admin
DB_PASSWORD=adminpassword
POSTGRES_DB=gosocial
POSTGRES_PORT=15432

DB_ADDR=postgres://admin:adminpassword@db:5432/gosocial?sslmode=disable

FROM_EMAIL=no-reply@example.com

SENDGRID_API_KEY=

MAILTRAP_API_KEY=
MAILTRAP_USERNAME=
MAILTRAP_PASSWORD=

REDIS_ENABLED=true
REDIS_ADDR=redis:6379
REDIS_PORT=16379

RATELIMITER_REQUESTS_COUNT=20
RATE_LIMITER_ENABLED=true
```

Important distinction:

```text
Inside the Dev Container:
PostgreSQL -> db:5432
Redis      -> redis:6379

From the host machine:
PostgreSQL -> localhost:15432
Redis      -> localhost:16379
Redis UI   -> http://localhost:8082
API        -> http://localhost:8080
```

Optional variables supported by the application include:

```env
AUTH_BASIC_USER=admin
AUTH_BASIC_PASS=admin
AUTH_TOKEN_SECRET=change-me-in-development
FRONTEND_URL=http://localhost:5173
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m
REDIS_PW=
REDIS_DB=0
```

Notes:

- `.env` is ignored by Git and must not be committed.
- `FRONTEND_URL` is currently used to build account activation links. The frontend itself is planned as a future separate project.
- User registration triggers email delivery. Configure Mailtrap to test the full registration and activation flow.
- If email credentials are not configured, the registration flow may fail when trying to send the activation email.

## Running Without the Dev Container

This workflow is available, but the Dev Container is recommended.

### 1. Install required tools

Install on the host machine:

- Go compatible with this project.
- Docker and Docker Compose.
- GNU Make.
- `air`.
- `swag`.
- `golang-migrate`.

Install tools manually:

```bash
go install github.com/air-verse/air@v1.65.1
go install github.com/swaggo/swag/cmd/swag@v1.16.6
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0
```

### 2. Start infrastructure

```bash
docker compose up -d
```

### 3. Configure host database URL

When running commands directly from the host, use the host-mapped PostgreSQL port:

```bash
export DB_ADDR="postgres://admin:adminpassword@localhost:15432/gosocial?sslmode=disable"
```

On Windows PowerShell:

```powershell
$env:DB_ADDR = "postgres://admin:adminpassword@localhost:15432/gosocial?sslmode=disable"
```

### 4. Run setup and API

```bash
make migrate-up
make seed
go run ./cmd/api
```

For live reload:

```bash
make dev
```

## Database, Migrations and Seed

Migration files are stored in:

```text
cmd/migrate/migrations
```

Create a new migration:

```bash
make migration name=create_users
```

Apply migrations:

```bash
make migrate-up
```

Rollback one migration:

```bash
make migrate-down
```

Drop the schema:

```bash
make migrate-drop
```

Reset the local database:

```bash
make reset-db
```

Run seed:

```bash
make seed
```

The current seed generates approximately:

- 100 users.
- 200 posts.
- 500 comments.

The default password for seeded users is:

```text
123123
```

## Swagger Documentation

Generated Swagger files are located in:

```text
docs/
```

Regenerate the documentation:

```bash
make gen-docs
```

With the API running, open:

```text
http://localhost:8080/v1/swagger/index.html
```

If the Swagger UI opens but does not load the schema, check `EXTERNAL_URL`, `ADDR`, and the generated files under `docs/`.

## Authentication

The project currently uses three main authentication and authorization flows.

### 1. Basic Auth for the health check

`GET /v1/health` requires Basic Auth.

```bash
curl -u admin:admin http://localhost:8080/v1/health
```

Credentials are configured through:

```env
AUTH_BASIC_USER=admin
AUTH_BASIC_PASS=admin
```

If these variables are not set, the application defaults to:

```text
admin:admin
```

### 2. User registration and account activation

Create a user:

```http
POST /v1/authentication/user
Content-Type: application/json
```

```json
{
  "username": "luiz",
  "email": "luiz@example.com",
  "password": "123123"
}
```

The backend:

1. hashes the password with bcrypt;
2. creates the user;
3. generates an activation token;
4. stores the token hash in the database;
5. sends an activation link by email;
6. returns the user and the plain activation token in the response.

Activate the user:

```http
PUT /v1/users/activate/{token}
```

### 3. JWT Bearer Token

Generate a token:

```http
POST /v1/authentication/token
Content-Type: application/json
```

```json
{
  "email": "luiz@example.com",
  "password": "123123"
}
```

Use the token in protected routes:

```http
Authorization: Bearer <token>
```

## Main Routes

Base URL:

```text
http://localhost:8080/v1
```

### Health

| Method | Route | Authentication | Description |
|---|---|---|---|
| `GET` | `/health` | Basic Auth | Returns API status, environment, and version. |

### Authentication

| Method | Route | Authentication | Description |
|---|---|---|---|
| `POST` | `/authentication/user` | Public | Registers a user and starts the activation flow. |
| `POST` | `/authentication/token` | Public | Generates a JWT for an active user. |

### Users

| Method | Route | Authentication | Description |
|---|---|---|---|
| `PUT` | `/users/activate/{token}` | Public | Activates a user using an activation token. |
| `GET` | `/users/{userId}` | Bearer Token | Fetches an active user profile. |
| `PUT` | `/users/{userId}/follow` | Bearer Token | Follows a user. |
| `PUT` | `/users/{userId}/unfollow` | Bearer Token | Unfollows a user. |
| `GET` | `/users/feed` | Bearer Token | Returns the paginated feed for the authenticated user. |

Feed query parameters:

| Parameter | Rule | Description |
|---|---|---|
| `limit` | `1..20` | Maximum number of posts. |
| `offset` | `>= 0` | Pagination offset. |
| `sort` | `asc` or `desc` | Sorts by creation date. |
| `tags` | up to 3 comma-separated tags | Filters by tags. |
| `search` | up to 100 characters | Searches in title or content. |
| `since` | `YYYY-MM-DD HH:MM:SS` | Parsed, but not currently applied in the SQL query. |
| `until` | `YYYY-MM-DD HH:MM:SS` | Parsed, but not currently applied in the SQL query. |

Example:

```bash
curl "http://localhost:8080/v1/users/feed?limit=10&offset=0&sort=desc&tags=GoLang,AI&search=testing" \
  -H "Authorization: Bearer <token>"
```

### Posts

| Method | Route | Authentication | Description |
|---|---|---|---|
| `POST` | `/posts/` | Bearer Token | Creates a post. |
| `GET` | `/posts/{postId}/` | Bearer Token | Fetches a post by ID, including comments. |
| `POST` | `/posts/{postId}/` | Bearer Token | Creates a comment on the post as the authenticated user. |
| `PATCH` | `/posts/{postId}/` | Bearer Token + ownership/role | Updates the title and/or content. |
| `DELETE` | `/posts/{postId}/` | Bearer Token + ownership/role | Deletes the post. |

Create a post:

```json
{
  "title": "First post",
  "content": "Post content",
  "tags": ["GoLang", "API"]
}
```

Update a post:

```json
{
  "title": "Updated title",
  "content": "Updated content"
}
```

Create a comment:

```json
{
  "content": "Example comment"
}
```

## Response Format

Successful responses are wrapped in `data`:

```json
{
  "data": {}
}
```

Error responses are wrapped in `error`:

```json
{
  "error": "error message"
}
```

HTTP status codes used by the project:

- `200 OK`
- `201 Created`
- `204 No Content`
- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`

## Data Model

Main tables:

### `users`

- `id`
- `email`
- `username`
- `password`
- `created_at`
- `is_active`
- `role_id`

### `roles`

- `id`
- `name`
- `level`
- `description`

Initial roles:

| Role | Level | General permission |
|---|---:|---|
| `user` | 1 | Can create posts and comments. |
| `moderator` | 2 | Can update posts from other users. |
| `admin` | 3 | Can update and delete posts from other users. |

### `user_invitations`

- `token`
- `user_id`
- `expiry`

### `posts`

- `id`
- `title`
- `content`
- `user_id`
- `tags`
- `created_at`
- `updated_at`
- `version`

### `comments`

- `id`
- `post_id`
- `user_id`
- `content`
- `created_at`

### `followers`

- `user_id`
- `follower_id`
- `created_at`

Relevant indexes and extensions:

- `CITEXT` for case-insensitive emails.
- `pg_trgm` for text search.
- GIN indexes for text search and tag filtering.
- Relationship indexes for users, posts, and comments.

## Useful Commands

### Development commands

```bash
# show tool versions
make tools

# download Go dependencies
make download

# tidy Go modules
make tidy

# run tests
make test

# run API with live reload
make dev

# run initial setup
make setup
```

### Database commands

```bash
# apply migrations
make migrate-up

# rollback one migration
make migrate-down

# drop schema
make migrate-drop

# run seed
make seed

# reset local database
make reset-db

# create a migration
make migration name=create_users
```

### Swagger commands

```bash
make gen-docs
```

### Docker commands

```bash
# start infrastructure
docker compose up -d

# stop infrastructure
docker compose down

# stop infrastructure and remove volumes
docker compose down --volumes
```

### Manual API run without Air

```bash
go build -buildvcs=false -o ./bin/main ./cmd/api
./bin/main
```

## Ports and Services

| Service | Inside Docker network | Host access | Description |
|---|---:|---:|---|
| API | `app:8080` | `localhost:8080` | Go backend API |
| PostgreSQL | `db:5432` | `localhost:15432` | Main database |
| Redis | `redis:6379` | `localhost:16379` | Optional cache |
| Redis Commander | `redis-commander:8081` | `localhost:8082` | Redis web UI |

## Known Limitations

- The project is still under development.
- The repository currently focuses on the backend API. A dedicated frontend application is planned as a future separate project.
- User registration depends on email delivery; without Mailtrap configured, the full registration and activation flow may fail.
- `since` and `until` are parsed in the feed query but are not currently applied in the SQL query.
- The Swagger schema URL may require adjustments to `ADDR` or `EXTERNAL_URL`, depending on how the API is executed.
- The current `Makefile` contains a duplicated `gen-docs` target name, including a Windows-oriented recipe. This should be reviewed to avoid ambiguity between Linux/macOS and Windows documentation generation workflows.
- When running directly on the host, `.env.example` cannot be used unchanged for database access because it targets Docker service names such as `db` and `redis`. Use `localhost:15432` and `localhost:16379` from the host.

## Suggested Roadmap

- Apply `since` and `until` filters to the feed SQL query.
- Review and split Linux/macOS and Windows Swagger generation commands in the `Makefile`.
- Add unit and integration tests.
- Improve mailer bootstrap error handling.
- Add refresh token support.
- Create a separate frontend project for login, registration, feed, and post creation.
- Add a CI pipeline.
- Review and standardize Swagger annotations.
- Add rate limiting.
- Add production-oriented Docker image separate from the Dev Container.
