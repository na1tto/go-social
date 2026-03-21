# go-social

A social feed backend API built with Go, PostgreSQL, Chi, and Swagger.

This project exposes endpoints to:

- manage posts and comments
- follow/unfollow users
- fetch a personalized feed with pagination and filtering

## Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Repository Structure](#repository-structure)
- [Requirements](#requirements)
- [Environment Variables](#environment-variables)
- [Quick Start (Local)](#quick-start-local)
- [Database Migrations](#database-migrations)
- [Database Seeding](#database-seeding)
- [Running the API](#running-the-api)
- [Swagger / API Docs](#swagger--api-docs)
- [API Reference](#api-reference)
- [Response Format and Error Handling](#response-format-and-error-handling)
- [Database Schema (High Level)](#database-schema-high-level)
- [Useful Commands](#useful-commands)
- [Known Limitations](#known-limitations)
- [Troubleshooting](#troubleshooting)
- [Roadmap Ideas](#roadmap-ideas)

## Project Overview

`go-social` is a backend service focused on social interactions:

- users create posts
- posts can have comments and tags
- users can follow/unfollow other users
- a feed endpoint returns posts from followed users (plus own posts) with filters

The codebase is organized with:

- `cmd/api` for HTTP server and handlers
- `internal/store` for repository/data-access layer
- `cmd/migrate/migrations` for SQL migrations
- `cmd/migrate/seed` + `internal/db/seed.go` for test data generation

## Tech Stack

- Go `1.25.3`
- PostgreSQL `16` (via Docker Compose)
- Router/middleware: `github.com/go-chi/chi/v5`
- Validation: `github.com/go-playground/validator/v10`
- SQL driver: `github.com/lib/pq`
- Swagger UI/docs: `github.com/swaggo/swag`, `github.com/swaggo/http-swagger/v2`
- SQL migrations: `golang-migrate` CLI (invoked via `Makefile`)

## Architecture

Layered flow:

1. HTTP handler (`cmd/api/*.go`)
2. Validation and request parsing
3. Repository call (`internal/store/*.go`)
4. PostgreSQL query execution
5. Standard JSON envelope response

Request middleware in use:

- request id
- real IP
- logger
- panic recoverer
- timeout (60s)

## Repository Structure

```text
go-social/
  cmd/
    api/                 # HTTP server, routes, handlers, JSON/error helpers
    migrate/
      migrations/        # SQL migration files
      seed/              # Seed entrypoint
  docs/                  # Swagger generated files
  internal/
    db/                  # DB connection + seed generator
    env/                 # env helper functions
    store/               # repository layer (users/posts/comments/followers)
  docker-compose.yml
  Makefile
  go.mod
```

## Requirements

Install the following on your machine:

- Go `1.25.3+`
- Docker + Docker Compose
- GNU Make
- `migrate` CLI (golang-migrate)
- `swag` CLI (only if regenerating docs)

Example installs:

```bash
# swag
go install github.com/swaggo/swag/cmd/swag@latest

# golang-migrate (choose one method based on your OS/package manager)
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```

## Environment Variables

The project uses `.env` (already referenced by `Makefile` and autoloaded by API).

Current variables used:

```env
ADDR=:8080
ENV=development

DB_USER=admin
DB_PASSWORD=adminpassword
POSTGRES_DB=gosocial
DB_ADDR=postgres://admin:adminpassword@localhost/gosocial?sslmode=disable

DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

EXTERNAL_URL=localhost:8080
```

Notes:

- `ADDR` should include `:` when running the server (for example `:8080`).
- `EXTERNAL_URL` is used as Swagger host metadata.
- `.env` is ignored by git, so keep your local copy.

## Quick Start (Local)

1. Clone the repository and enter the project folder.
2. Ensure your `.env` exists with valid values.
3. Start PostgreSQL with Docker Compose.
4. Run migrations.
5. (Optional) seed database.
6. Run the API.

```bash
# 1) start postgres
docker compose up -d db

# 2) apply migrations
make migrate-up

# 3) optional: seed fake data
make seed

# 4) run API
go run ./cmd/api
```

Health check:

```bash
curl http://localhost:8080/v1/health
```

## Database Migrations

The `Makefile` includes migration helpers:

```bash
# create a migration file pair
make migration add_some_change

# apply all up migrations
make migrate-up

# rollback N migrations
make migrate-down 1
```

Migration files live in `cmd/migrate/migrations`.

## Database Seeding

Populate database with sample data:

```bash
make seed
```

Current seed behavior:

- creates ~100 users
- creates ~200 posts
- creates ~500 comments

## Running the API

```bash
go run ./cmd/api
```

Server defaults:

- bind address: `:8080`
- base path: `/v1`
- request timeout middleware: `60s`

## Swagger / API Docs

Generated docs are under `docs/`.

To regenerate:

```bash
make gen-docs
```

After server startup, open Swagger UI:

- `http://localhost:8080/v1/swagger/index.html`

## API Reference

Base URL:

```text
http://localhost:8080/v1
```

### Health

- `GET /health`
  - Returns service status, environment, and version.

### Posts

- `POST /posts/`
  - Create post
  - Body:

```json
{
  "title": "Post title",
  "content": "Post body",
  "tags": ["GoLang", "Productivity"]
}
```

- `GET /posts/{postId}`
  - Get a post by ID (includes comments)

- `PATCH /posts/{postId}`
  - Partial update of `title` and/or `content`

- `DELETE /posts/{postId}`
  - Delete a post

- `POST /posts/{postId}`
  - Create a comment on a post
  - Body:

```json
{
  "content": "Nice post!"
}
```

### Users

- `GET /users/{userId}`
  - Fetch user profile

- `PUT /users/{userId}/follow`
  - Follow user
  - Body:

```json
{
  "user_id": 40
}
```

- `PUT /users/{userId}/unfollow`
  - Unfollow user
  - Body:

```json
{
  "user_id": 40
}
```

### Feed

- `GET /users/feed`
  - Returns feed posts with metadata (`comments_count`, author username)
  - Supports query params:
    - `limit` (1..20)
    - `offset` (>=0)
    - `sort` (`asc` or `desc`)
    - `tags` (comma-separated, up to 3)
    - `search` (max 100 chars)
    - `since` (datetime)
    - `until` (datetime)

Example:

```bash
curl "http://localhost:8080/v1/users/feed?limit=10&offset=0&sort=desc&tags=GoLang,AI&search=testing"
```

## Response Format and Error Handling

Success responses are wrapped in:

```json
{
  "data": {}
}
```

Error responses are wrapped in:

```json
{
  "error": "message"
}
```

Common statuses used:

- `200 OK`
- `201 Created`
- `204 No Content`
- `400 Bad Request`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`

## Database Schema (High Level)

Main tables:

- `users`
  - `id`, `email` (CITEXT unique), `username` (unique), `password`, `created_at`

- `posts`
  - `id`, `title`, `content`, `user_id` (FK users), `tags[]`, `created_at`, `updated_at`, `version`

- `comments`
  - `id`, `post_id`, `user_id`, `content`, `created_at`

- `followers`
  - composite PK: `(user_id, follower_id)`
  - both columns FK to `users` with cascade delete

Extensions/indexes created by migrations:

- `CITEXT`
- `pg_trgm`
- GIN trigram indexes for text search on comments/posts title
- GIN index on `posts.tags`
- B-tree indexes for `users.username`, `posts.user_id`, `comments.post_id`

## Useful Commands

```bash
# start postgres
docker compose up -d db

# stop postgres
docker compose down

# apply migrations
make migrate-up

# rollback one migration
make migrate-down 1

# seed database
make seed

# run api
go run ./cmd/api

# regenerate swagger docs
make gen-docs
```

## Known Limitations

Current implementation has a few intentional/temporary constraints:

- No authentication yet.
- Some handlers still use hardcoded user IDs (for example creating posts/comments and feed owner).
- Password hashing is not implemented in this code path.
- `since`/`until` query params are parsed but not yet applied in feed SQL filtering.
- Swagger metadata and some handler annotations may not perfectly match runtime behavior.

## Troubleshooting

- Migration command fails with "unknown command migrate":
  - install `golang-migrate` CLI and ensure it is in your PATH.

- API cannot connect to DB:
  - check `DB_ADDR`
  - ensure `docker compose up -d db` is running
  - verify credentials in `.env` match container env (`DB_USER`, `DB_PASSWORD`, `POSTGRES_DB`)

- Swagger page opens but no schema:
  - run `make gen-docs`
  - restart API after regenerating docs

- Port already in use:
  - change `ADDR` in `.env` (for example `:8081`)

## Roadmap Ideas

- JWT authentication and authorization
- password hashing + secure user registration/login
- complete `since/until` feed filters in SQL
- cursor-based pagination
- rate limiting middleware
- unit/integration tests + CI pipeline
- consistent API error model and OpenAPI cleanup
