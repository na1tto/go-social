```text
  ____          ____             _       _
 / ___| ___    / ___|  ___   ___(_) __ _| |
| |  _ / _ \   \___ \ / _ \ / __| |/ _` | |
| |_| | (_) |   ___) | (_) | (__| | (_| | |
 \____|\___/   |____/ \___/ \___|_|\__,_|_|
```

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-metrics-E6522C?logo=prometheus&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-API-85EA2D?logo=swagger&logoColor=black)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-web-3178C6?logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-web-646CFF?logo=vite&logoColor=white)

**[Português](README.md) · [English](README.en.md)**

</div>

---

# Go Social

A social network API written in Go, built to develop backend engineering skills through authentication, persistence, caching, and observability in a single project.

The core already supports user registration, posts, comments, following users, and a personalized feed. The current focus is making these flows more reliable and understanding application behavior under load.

> **Under development.** An initial frontend lives in `web/`, with account confirmation and a placeholder home page. Login, feed, and publishing through the interface are still part of the project's evolution.

<details>
<summary><kbd>Implemented features · click to expand</kbd></summary>

- Registration with bcrypt, email activation, and JWT authentication.
- Ownership and role permissions: `user`, `moderator`, and `admin`.
- Create, read, update, and delete posts; create comments.
- Following users and a feed with pagination, sorting, text search, and tag filtering.
- Optional Redis caching for authenticated users.
- In-memory rate limiting, structured logs, and Prometheus metrics.
- SQL migrations, mock-based tests, CI, and a containerized development environment.

</details>

## Contents

- [01 · Quick start](#01--quick-start)
- [02 · Technologies and architecture](#02--technologies-and-architecture)
- [03 · Configuration](#03--configuration)
- [04 · Using the API](#04--using-the-api)
- [05 · Observability and decisions](#05--observability-and-decisions)
- [06 · Limitations and next steps](#06--limitations-and-next-steps)
- [07 · Development](#07--development)

---

## 01 · Quick start

The recommended workflow is the **Dev Container**. It requires Git, Docker with Compose, and a compatible editor such as VS Code with Dev Containers or Zed. On Windows, use Docker Desktop with WSL 2.

From the project root:

```sh
cp .env.example .env
```

In PowerShell, use `Copy-Item .env.example .env`. Open the project in the Dev Container and run in its terminal:

```sh
make migrate-up
make dev
```

The API is available at **http://localhost:8080**. Check it with `curl http://localhost:8080/v1/health`; this route is currently public. The container provides Go, Air, Swag, and Migrate; migrations are run manually.

To test registration and activation, configure Mailtrap as described in [configuration](#03--configuration). [DEVCONTAINER.md](DEVCONTAINER.md) covers editor setup, tools, and troubleshooting.

<details>
<summary><kbd>Running on the host machine · alternative</kbd></summary>

Install Go compatible with [go.mod](go.mod), Docker Compose, and the Migrate CLI with PostgreSQL support. Make and Air are needed for the development shortcuts.

In `.env`, replace Docker network addresses with the host ports:

```dotenv
DB_ADDR=postgres://admin:adminpassword@localhost:15432/gosocial?sslmode=disable
REDIS_ADDR=localhost:16379
```

Then:

```sh
docker compose up -d
make migrate-up
go run ./cmd/api
```

The API loads `.env` automatically; the Makefile also reads it. The root Compose file starts the infrastructure, but does not start the API.

</details>

---

## 02 · Technologies and architecture

| Layer | Tools and responsibility |
|---|---|
| API | Go · Chi · go-playground/validator |
| Persistence | PostgreSQL 16 · SQL with `database/sql` and `lib/pq` · golang-migrate |
| Authentication | JWT · bcrypt · roles and ownership checks |
| Cache and request limits | Redis 7 for users · in-memory fixed-window rate limiting |
| Email | Active Mailtrap SMTP sandbox · SendGrid implementation available, unused at bootstrap |
| Observability | Zap · Prometheus · expvar |
| Initial interface | React 19 · TypeScript · Vite · React Router |
| Development | Docker Compose · Dev Containers · Air · Swag · GitHub Actions |

The API is a single Go service organized into layers. Middleware handles authentication and shared request concerns; handlers validate input and coordinate operations; `internal/store` manages database access. There is no separate business service layer yet.

```text
HTTP client / web (React)
          |
          v
  Chi + middleware
          |
          v
  Handlers (cmd/api) ------> Mailtrap (activation)
          |
          +---------------> Redis (user cache)
          |
          v
  Store (internal/store)
          |
          v
      PostgreSQL

  API ----> Zap logs
  Prometheus ----> API /metrics
```

The data model connects users and roles, activation invitations, posts, comments, and followers. Posts have a version for concurrent update control; the database uses `CITEXT`, `pg_trgm`, and indexes for search and tags.

| Directory | What to find |
|---|---|
| `cmd/api/` | Routes, handlers, middleware, and HTTP responses |
| `internal/store/` | SQL queries, interfaces, and cache |
| `internal/auth/`, `internal/mailer/` | JWT and email delivery |
| `internal/observability/`, `internal/ratelimiter/` | Metrics and request limits |
| `cmd/migrate/`, `internal/db/` | Migrations, connection, and seed |
| `web/` | Initial account confirmation interface |
| `docs/` | Generated Swagger and decisions in `docs/adr/` |
| `.devcontainer/`, `.github/workflows/` | Development environment and CI |

---

## 03 · Configuration

Start from [.env.example](.env.example). It targets the Dev Container network; `.env` is ignored by Git.

| Variable | Current use |
|---|---|
| `ADDR` / `ENV` | HTTP address (`:8080`) and environment |
| `DB_ADDR` | PostgreSQL: `db:5432` in the container or `localhost:15432` on the host |
| `DB_USER`, `DB_PASSWORD`, `POSTGRES_DB` | Credentials and database created by Compose |
| `REDIS_ENABLED` / `REDIS_ADDR` | Cache: `redis:6379` in the container or `localhost:16379` on the host |
| `AUTH_TOKEN_SECRET` | JWT secret; add your own value to `.env` |
| `AUTH_BASIC_USER` / `AUTH_BASIC_PASS` | Protect `/v1/debug/vars`; local default is `admin:admin` |
| `FROM_EMAIL` / `MAILTRAP_*` | Sender and email configuration |
| `FRONTEND_URL` | Activation link base; defaults to `http://localhost:5173` |
| `RATE_LIMITER_ENABLED` / `RATELIMITER_REQUESTS_COUNT` | Local limit; defaults to 20 requests per 3-second window |
| `EXTERNAL_URL` | Sets the Swagger host; integration still needs review |

**Email:** fill in `FROM_EMAIL`, `MAILTRAP_API_KEY`, `MAILTRAP_USERNAME`, and `MAILTRAP_PASSWORD`. The constructor requires an API key, but delivery currently uses the SMTP sandbox username and password. If sending fails, registration attempts to undo user creation.

**CORS:** the example contains `CORS_ALLOWED_ORIGIN`, but the code reads `FRONTEND_ORIGIN` and currently accepts any origin through a callback. The method list also omits `PATCH`. Fixing this is on the roadmap.

**Local services:** API `:8080` · PostgreSQL `:15432` · Redis `:16379` · Redis Commander `:8082` · Prometheus `:9090`.

<details>
<summary><kbd>Activation interface · optional</kbd></summary>

With Node compatible with [web/.nvmrc](web/.nvmrc) installed, run in another terminal:

```sh
cd web
npm ci
npm run dev
```

The interface calls `http://localhost:8080/v1` by default. To change the API, set `VITE_API_URL` in the Vite environment. Keep `FRONTEND_URL` aligned with the port shown by the frontend; email links open `/confirm/:token`.

</details>

---

## 04 · Using the API

Local base: `http://localhost:8080/v1`. The flow is **registration → activation → JWT issuance → authenticated routes**.

| Method | Route relative to `/v1` | Access / purpose |
|---|---|---|
| `POST` | `/authentication/user` | Public · accepts `username`, `email`, and `password` |
| `PUT` | `/users/activate/{token}` | Public · activates the account |
| `POST` | `/authentication/token` | Public · accepts `email` and `password`, returns JWT |
| `GET` | `/users/{userId}/` | JWT · retrieves a user |
| `PUT` | `/users/{userId}/follow` or `/users/{userId}/unfollow` | JWT · follows or unfollows |
| `GET` | `/users/feed` | JWT · paginated feed |
| `POST` | `/posts/` | JWT · creates a post with `title`, `content`, and `tags` |
| `GET` / `POST` | `/posts/{postId}/` | JWT · reads a post with comments / comments with `content` |
| `PATCH` / `DELETE` | `/posts/{postId}/` | JWT + ownership or role · updates / deletes |

Send `Authorization: Bearer <token>` on protected routes. Authors can edit and delete their own posts; moderators can edit other users' posts, and admins can also delete them.

The feed accepts `limit` (1–20), `offset` (≥ 0), `sort` (`asc`/`desc`), `tags` (up to 3, comma-separated), and `search` (up to 100 characters). `since` and `until` are parsed, but do not yet filter the SQL query.

JSON responses use `{"data": ...}` for success and `{"error": "..."}` for failures; deletions may return `204` without a body.

**Reference:** [Swagger YAML](docs/swagger.yaml) and [Swagger JSON](docs/swagger.json). The UI is at `/v1/swagger/index.html`, but the schema URL assembled by the application still needs correction; use the files if it does not load.

---

## 05 · Observability and decisions

Logs include `request_id`, method, route, status, duration, and context state. Metrics for `/v1` routes track request volume, duration, and in-flight requests.

| Endpoint | Current access |
|---|---|
| `GET /v1/health` | Public · status, environment, and version |
| `GET /metrics` | Public · Prometheus collection |
| `GET /v1/debug/vars` | Basic Auth · runtime and connection pool |

The [local Prometheus configuration](internal/observability/prometheus/prometheus.yml) scrapes `host.docker.internal:8080` every 15 seconds. Its UI is at **http://localhost:9090**.

### ADR 0001 · HTTP cancellations

[ADR 0001 — Handling HTTP request cancellations](docs/adr/0001-http-request-cancellation.md) (in Portuguese), accepted on **2026-10-03**, records a decision prompted by load testing: recognized cancellations should become neither internal errors nor artificial successes.

The decision calls for Info-level logging without a JSON response, correlation through `request_id`, and preservation of any status already written. When no status has been written and the context is canceled, logs use `0` and metrics use `status="canceled"`. That `0` exists only in observability; it is not an HTTP code sent to the client. Deadline expiration is outside the ADR's scope.

**Pending alignment:** the condition in `errors.go` currently checks only `errors.Is(err, context.Canceled)`; it must also require cancellation of the HTTP context, as documented in the decision.

---

## 06 · Limitations and next steps

The project already has mock-based tests, CI, rate limiting, and a runtime Dockerfile. The next stage is extending and consolidating these capabilities. This table connects verified code limitations with proposed improvements, without delivery commitments.

| Current limitation | Proposed next step |
|---|---|
| Registration depends on synchronous Mailtrap sandbox delivery; mailer initialization errors are ignored | Validate bootstrap and improve delivery, retries, and registration recovery |
| Activation tokens are also returned in JSON; no refresh tokens | Review the activation contract and extend the authentication lifecycle |
| Permissive CORS, inconsistent configuration, and missing `PATCH` | Unify configuration, restrict origins, and cover supported methods |
| Feed date filters do not reach SQL | Implement `since`/`until` and validate with integration tests |
| Enabled caching propagates Redis failures | Define fallback and invalidation behavior |
| Rate limiter is in memory per process and retains old entries | Define cleanup and a strategy for multiple instances |
| Swagger has schema URL and generation command issues | Align routes, metadata, and generation across environments |
| Cancellation handling still differs from the ADR | Complete the condition and expand cancellation and timeout tests |
| Frontend only covers account confirmation | Add registration, login, feed, and publishing |
| Dockerfile uses `cmd/api/*.go`, including test files in the build | Fix and validate the image; expand integration tests and deployment validation |

Before publishing an instance, replace example secrets and define access to `/metrics` and diagnostic routes. The current configuration targets development.

---

## 07 · Development

Inside the Dev Container, with `.env` created:

```sh
make test                            # Go tests
go vet ./...                         # static analysis
make migrate-up                      # apply migrations
make migration name=change_name      # create a SQL migration
make dev                             # API with live reload
```

The [audit workflow](.github/workflows/audit.yaml) already runs dependency verification, build, `go vet`, Staticcheck, and tests with `-race`. The [CHANGELOG](CHANGELOG.md) records releases; [ADRs](docs/adr/) explain architectural decisions.

<details>
<summary><kbd>Demo database and documentation generation</kbd></summary>

`make seed` generates 100 users, 200 posts, and 500 comments with the demo password `123123`. The seed executable reads `DB_ADDR` from the environment: outside the Dev Container, export it before running, since the seed does not load `.env` itself.

`make migrate-down` rolls back one migration. `make reset-db` drops the schema, recreates it, and seeds the database; use only with a disposable database.

The Makefile provides `make gen-docs` and `make gen-docs-win`. The first still combines the search directory and main file path inconsistently. A direct alternative from the project root is:

```sh
swag init -g ./api/main.go -d cmd,internal --parseDependency --parseInternal
swag fmt
```

Review generated files before including them in a commit.

</details>

---

**Documentation reviewed:** 2026-10-03 · **Languages:** [Português](README.md) / [English](README.en.md)
