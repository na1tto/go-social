# Go Social

Backend API for a simple social network, focused on user registration, authentication, posts, comments, followers and a personalized feed.

The project is built with Go, PostgreSQL, Chi, JWT authentication, Swagger documentation, optional Redis caching and account activation by email.

> Status: this project is under development.

---

## Table of Contents

* [Tech Stack](#tech-stack)
* [Development Environment](#development-environment)
* [Requirements](#requirements)
* [Environment Variables](#environment-variables)
* [Using the Dev Container with VS Code](#using-the-dev-container-with-vs-code)
* [Using the Dev Container with Zed](#using-the-dev-container-with-zed)
* [Development Workflow Inside the Container](#development-workflow-inside-the-container)
* [Database, Migrations and Seed](#database-migrations-and-seed)
* [Swagger Documentation](#swagger-documentation)
* [Useful Commands](#useful-commands)
* [Ports and Services](#ports-and-services)
* [Git Workflow](#git-workflow)
* [Troubleshooting](#troubleshooting)
* [Project Structure](#project-structure)

---

## Tech Stack

* Go `1.25.3`
* PostgreSQL `16`
* Redis `7`, optional cache
* Chi router
* JWT with `github.com/golang-jwt/jwt/v5`
* Bcrypt via `golang.org/x/crypto`
* Request validation with `go-playground/validator`
* Logging with `go.uber.org/zap`
* Swagger/OpenAPI with `swaggo`
* Database migrations with `golang-migrate`
* Live reload with `air`
* Email delivery through Mailtrap, with structure also prepared for SendGrid

---

## Development Environment

This repository contains a Dev Container setup under:

```text
.devcontainer/
├── Dockerfile
├── devcontainer.json
├── docker-compose.devcontainer.yml
└── post-create.sh
```

The Dev Container provides:

* a Go development container using the `app` service;
* a mounted workspace at `/workspace`;
* PostgreSQL via the `db` service;
* Redis via the `redis` service;
* Redis Commander via the `redis-commander` service;
* Go module/build cache volumes;
* development tools installed inside the container: `air`, `swag` and `migrate`.

Recommended rule for this project:

```text
Use the container for Go, Air, migrations, seed, Swagger and tests.
Use the host machine for Git push/pull if your SSH credentials are configured outside the container.
```

---

## Requirements

Install on the host machine:

* Docker Desktop;
* Docker Compose;
* Git;
* VS Code with the Dev Containers extension, or Zed with Dev Containers support;
* WSL 2 enabled if you are developing on Windows.

You do not need to install Go, `air`, `swag` or `migrate` on the host when using the Dev Container. These tools are installed inside the container.

---

## Environment Variables

Create a local `.env` file from `.env.example`.

### Linux/macOS/WSL

```bash
cp .env.example .env
```

### Windows PowerShell

```powershell
Copy-Item .env.example .env
```

Example `.env`:

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
REDIS_COMMANDER_PORT=8082
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

Do not commit your real `.env` file.

---

## Using the Dev Container with VS Code

### 1. Clone the repository

```bash
git clone https://github.com/na1tto/go-social.git
cd go-social
git checkout chore/devcontainer-collaborative-setup
```

### 2. Create the `.env` file

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

### 3. Open the project in VS Code

Open the repository root in VS Code. The correct folder is the one that contains:

```text
.git/
.devcontainer/
docker-compose.yml
go.mod
Makefile
```

### 4. Reopen in container

Open the Command Palette:

```text
Ctrl + Shift + P
```

Run:

```text
Dev Containers: Reopen in Container
```

VS Code will build the image, start the Compose services and connect the editor to the `app` container.

### 5. Open a terminal inside the container

After the container opens:

```text
Terminal > New Terminal
```

The terminal should be inside `/workspace`.

Validate:

```bash
pwd
ls -la
go version
air -v
swag --version
migrate -version
```

Expected workspace path:

```text
/workspace
```

### 6. Start development

```bash
make download
make migrate-up
make seed
make dev
```

The API should be available at:

```text
http://localhost:8080
```

Health check:

```bash
curl -u admin:admin http://localhost:8080/v1/health
```

---

## Using the Dev Container with Zed

### 1. Clone the repository

```bash
git clone https://github.com/na1tto/go-social.git
cd go-social
git checkout chore/devcontainer-collaborative-setup
```

### 2. Create the `.env` file

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

### 3. Open the project root in Zed

Open the repository root, not the `.devcontainer` folder.

The correct folder contains:

```text
.git/
.devcontainer/
docker-compose.yml
go.mod
Makefile
```

### 4. Open in Dev Container

When Zed detects `.devcontainer/devcontainer.json`, it may show a prompt asking whether to open the project inside the container. Choose:

```text
Open in Container
```

If the prompt does not appear, open Zed's command palette and use the remote/dev container option, usually available through:

```text
Project: Open Remote
```

Then select the option to connect/open the project in a Dev Container.

### 5. Validate the container terminal

Open a terminal in Zed and validate:

```bash
pwd
ls -la
go version
air -v
swag --version
migrate -version
```

Expected path:

```text
/workspace
```

If `/workspace` is empty, Zed probably opened the wrong folder or the bind mount was not created correctly. Close the container session, reopen the repository root and reconnect to the Dev Container.

### 6. Start development

```bash
make download
make migrate-up
make seed
make dev
```

If the API is not reachable through `localhost:8080` in Zed, start the Compose stack manually from the host and then enter the app container:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  up -d --build
```

Then open a shell in the `app` container:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  exec app bash
```

Inside the container:

```bash
cd /workspace
make dev
```

---

## Development Workflow Inside the Container

The daily workflow should be executed from `/workspace` inside the `app` container.

### First run

```bash
cd /workspace
make download
make migrate-up
make seed
make dev
```

### Normal development

Use one terminal for the API:

```bash
make dev
```

Use a second terminal for commands:

```bash
git status
go test ./...
make migrate-up
make gen-docs
```

To stop the running process, press:

```text
Ctrl + C
```

---

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

Drop and recreate the local database schema:

```bash
make reset-db
```

Run seed:

```bash
make seed
```

The seed creates test data for local development. Do not use it against production data.

---

## Swagger Documentation

Generated Swagger files are located in:

```text
docs/
```

Regenerate documentation:

```bash
make gen-docs
```

With the API running, open:

```text
http://localhost:8080/v1/swagger/index.html
```

---

## Useful Commands

### Container lifecycle from the host

Start or rebuild the development stack:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  up -d --build
```

Stop the development stack:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  down --remove-orphans
```

Stop and remove volumes, including local database data:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  down --remove-orphans --volumes
```

Enter the app container manually:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  exec app bash
```

### Commands inside the container

```bash
# show tool versions
make tools

# download Go dependencies
make download

# tidy Go modules
make tidy

# run tests
make test

# apply migrations
make migrate-up

# rollback one migration
make migrate-down

# run seed
make seed

# run API with live reload
make dev

# regenerate Swagger docs
make gen-docs
```

### Manual API run without Air

```bash
go build -buildvcs=false -o ./bin/main ./cmd/api
./bin/main
```

---

## Ports and Services

| Service         |  Inside Docker network |       Host access | Description    |
| --------------- | ---------------------: | ----------------: | -------------- |
| API             |             `app:8080` |  `localhost:8080` | Go backend API |
| PostgreSQL      |              `db:5432` | `localhost:15432` | Main database  |
| Redis           |           `redis:6379` | `localhost:16379` | Optional cache |
| Redis Commander | `redis-commander:8081` |  `localhost:8082` | Redis web UI   |

---

## Git Workflow

The workspace is mounted as a bind mount:

```text
host project folder <-> container /workspace
```

That means files edited inside `/workspace` are the same files stored on disk in the host project folder.

Recommended workflow for this project:

```text
Use Git normally on the host machine.
Use the container for development commands.
```

Example:

Inside the container:

```bash
make dev
go test ./...
```

On the host:

```bash
git status
git add .
git commit -m "chore: improve devcontainer workflow"
git push
```

If you choose to use Git inside the container and receive a safe directory error, run:

```bash
git config --global --add safe.directory /workspace
git config core.filemode false
```

The `post-create.sh` script also applies these settings during container setup.

---

## Troubleshooting

### `/workspace` is empty

The project root was probably not mounted correctly.

Check inside the container:

```bash
pwd
ls -la /workspace
```

You should see:

```text
.git
.devcontainer
go.mod
Makefile
docker-compose.yml
cmd
internal
```

If these files are missing, close the Dev Container and reopen the repository root, not the `.devcontainer` folder.

---

### Git shows many modified files that were not edited

This is usually caused by Linux/Windows file mode differences.

Run inside `/workspace`:

```bash
git config core.filemode false
git status
```

---

### `git push` fails inside the container with `Permission denied (publickey)`

The container does not have access to your host SSH key by default.

Recommended solution:

```text
Run git push from the host machine instead of inside the container.
```

Because `/workspace` is mounted from the host, commits and file changes are shared.

---

### Port `5432` is already allocated

The project maps PostgreSQL to `localhost:15432` on the host to avoid conflicts with a local PostgreSQL installation.

Check if another container is using the port:

```powershell
docker ps --format "table {{.Names}}\t{{.Ports}}"
```

Stop the stack:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  down --remove-orphans
```

If you need a clean reset:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  down --remove-orphans --volumes
```

---

### Database container is unhealthy

Check logs:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  logs --tail=100 db
```

If the local database volume was created with an old configuration, reset volumes:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  down --remove-orphans --volumes
```

Then rebuild:

```powershell
docker compose `
  -f docker-compose.yml `
  -f .devcontainer\docker-compose.devcontainer.yml `
  up -d --build
```

---

### Air starts but appears to do nothing

Run Air explicitly:

```bash
air -c .air.toml
```

If it still looks idle, test the build manually:

```bash
go build -buildvcs=false -o ./bin/main ./cmd/api
./bin/main
```

In another terminal:

```bash
curl -u admin:admin http://localhost:8080/v1/health
```

If `./bin/main` runs without output, the server may still be active. Some Go applications do not print logs when starting.

---

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
├── web/
├── .air.toml
├── .dockerignore
├── .env.example
├── .gitattributes
├── .gitignore
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Recommended Daily Flow

```bash
cd /workspace
make download
make migrate-up
make dev
```

In another terminal:

```bash
go test ./...
make gen-docs
git status
```

Run `git commit` and `git push` on the host machine if your Git credentials are configured outside the container.
