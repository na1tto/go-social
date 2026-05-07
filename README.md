# Go Social

Backend API for a simple social network, focused on user registration, authentication, posts, comments, followers, and a personalized feed.

The project is built with Go, PostgreSQL, Chi, JWT authentication, Swagger documentation, optional Redis caching, and account activation by email.

> Status: this project is under development. Some routes are already implemented, but there are known limitations documented at the end of this README.

## Table of Contents

- [Tech Stack](#tech-stack)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Environment Variables](#environment-variables)
- [Running Locally](#running-locally)
- [Database, Migrations and Seed](#database-migrations-and-seed)
- [Swagger Documentation](#swagger-documentation)
- [Authentication](#authentication)
- [Main Routes](#main-routes)
- [Response Format](#response-format)
- [Data Model](#data-model)
- [Useful Commands](#useful-commands)
- [Known Limitations](#known-limitations)
- [Suggested Roadmap](#suggested-roadmap)

## Tech Stack

- Go `1.25.3`
- PostgreSQL `16`
- Optional Redis cache
- Chi router
- Chi CORS middleware
- JWT with `github.com/golang-jwt/jwt/v5`
- Bcrypt via `golang.org/x/crypto`
- Request validation with `go-playground/validator`
- Logging with `go.uber.org/zap`
- Swagger/OpenAPI with `swaggo`
- Database migrations with `golang-migrate`
- Email delivery through Mailtrap, with structure also prepared for SendGrid

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

The frontend is not part of the current scope of this backend documentation. A dedicated frontend application is planned as a future project.

## Project Structure

```text
go-social/
├── cmd/
│   ├── api/
│   │   ├── api.go
│   │   ├── auth.go
│   │   ├── feed.go
│   │   ├── health.go
│   │   ├── middleware.go
│   │   ├── posts.go
│   │   └── users.go
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
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Requirements

Install the following tools:

- Go compatible with the version defined in `go.mod`.
- Docker and Docker Compose.
- GNU Make.
- `golang-migrate`.
- `swag`, only if you need to regenerate the Swagger documentation.

Install `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Install `golang-migrate`:

```bash
# Choose the installation method that matches your operating system.
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```

## Environment Variables

Create a `.env` file in the project root. This file is automatically loaded by the backend and is also used by the `Makefile`.

Example for local development:

```env
# API
ADDR=:8080
ENV=development
EXTERNAL_URL=localhost:8080
FRONTEND_URL=http://localhost:5173

# PostgreSQL / Docker Compose
DB_USER=admin
DB_PASSWORD=adminpassword
POSTGRES_DB=gosocial
DB_ADDR=postgres://admin:adminpassword@localhost/gosocial?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

# Redis
REDIS_ENABLED=false
REDIS_ADDR=localhost:6379
REDIS_PW=
REDIS_DB=0

# Health check Basic Auth
AUTH_BASIC_USER=admin
AUTH_BASIC_PASS=admin

# JWT
AUTH_TOKEN_SECRET=change-me-in-development

# Email / account activation
FROM_EMAIL=no-reply@example.com
MAILTRAP_API_KEY=
MAILTRAP_USERNAME=
MAILTRAP_PASSWORD=
SENDGRID_API_KEY=
```

Notes:

- `ADDR` must use a valid `http.Server` address format, such as `:8080`.
- User registration triggers email delivery. Configure Mailtrap to test the full registration and activation flow.
- `REDIS_ENABLED=false` makes the backend read authenticated users directly from PostgreSQL.
- `.env` is ignored by Git, so it must be created locally.
- `FRONTEND_URL` is currently used to build the account activation link. A dedicated frontend will be developed later as a separate project.

## Running Locally

Clone the repository:

```bash
git clone https://github.com/na1tto/go-social.git
cd go-social
```

Start PostgreSQL, Redis, and Redis Commander:

```bash
docker compose up -d
```

Run the migrations:

```bash
make migrate-up
```

Optionally seed the database with test data:

```bash
make seed
```

Run the API:

```bash
go run ./cmd/api
```

Test the health check with Basic Auth:

```bash
curl -u admin:admin http://localhost:8080/v1/health
```

Expected response:

```json
{
  "data": {
    "env": "development",
    "status": "ok",
    "version": "0.0.1"
  }
}
```

## Database, Migrations and Seed

Migration files are stored in:

```text
cmd/migrate/migrations
```

Create a new migration:

```bash
make migration migration_name
```

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down 1
```

Run the database seed:

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

If the Swagger UI opens but does not load the schema, check `EXTERNAL_URL`, `ADDR`, and regenerate the documentation with `make gen-docs`.

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
| `GET` | `/users/feed` | Bearer Token | Returns the paginated user feed. |

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
| `POST` | `/posts/{postId}/` | Bearer Token | Creates a comment on the post. |
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

```bash
# start local infrastructure
docker compose up -d

# stop local infrastructure
docker compose down

# apply migrations
make migrate-up

# rollback one migration
make migrate-down 1

# create a new migration
make migration migration_name

# run seed
make seed

# run API
go run ./cmd/api

# generate Swagger docs
make gen-docs
```

Docker services:

| Service | Port | Description |
|---|---:|---|
| PostgreSQL | `5432` | Main database. |
| Redis | `6379` | Optional cache. |
| Redis Commander | `8081` | Local UI for inspecting Redis. |

Redis Commander:

```text
http://127.0.0.1:8081
```

## Known Limitations

- The project is still under development and does not currently include a test suite or CI pipeline.
- The frontend is not part of the current implementation scope. A dedicated frontend application is planned for a future project.
- User registration depends on email delivery; without Mailtrap configured, the full registration and activation flow may fail.
- `since` and `until` are parsed in the feed query but are not currently applied in the SQL query.
- The Swagger schema URL may require adjustments to `ADDR` or `EXTERNAL_URL`, depending on how the API is executed.
- The `Makefile` contains a `powershell` step in `make gen-docs`, which may require adjustment on Linux or macOS.

## Suggested Roadmap

- Remove hardcoded IDs and always use the authenticated user from the JWT.
- Apply `since` and `until` filters to the feed SQL query.
- Add unit and integration tests.
- Create a `.env.example` file.
- Improve mailer bootstrap error handling.
- Add refresh token support.
- Create a separate frontend project for login, registration, feed, and post creation.
- Add a CI pipeline.
- Review and standardize Swagger annotations.
- Add rate limiting.
