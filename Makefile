MIGRATIONS_PATH ?= ./cmd/migrate/migrations
DB_ADDR ?= postgres://admin:adminpassword@localhost:5432/gosocial?sslmode=disable
MIGRATION_NAME ?=

.PHONY: dev
dev:
	@air -c .air.toml

.PHONY: tools
tools:
	@go version
	@air -v
	@swag --version
	@migrate -version

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: download
download:
	@go mod download

.PHONY: test
test:
	@go test ./...

.PHONY: wait-db
wait-db:
	@until pg_isready -d "$(DB_ADDR)" > /dev/null 2>&1; do \
		echo "waiting for postgres..."; \
		sleep 1; \
	done

.PHONY: migration
migration:
	@test -n "$(name)" || (echo "Uso: make migration name=create_users"; exit 1)
	@migrate create -seq -ext sql -dir "$(MIGRATIONS_PATH)" "$(name)"

.PHONY: migrate-up
migrate-up:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" down 1

.PHONY: migrate-drop
migrate-drop:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" drop -f

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./cmd/api/main.go -d ./cmd,./internal --parseDependency --parseInternal
	@swag fmt
	@sed -i 's|swag/v2|swag|g' docs/docs.go

.PHONY: gen-docs-windows
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal --parseDependency --parseInternal && swag fmt
	@powershell -Command "(Get-Content docs/docs.go) -replace 'swag/v2', 'swag' | Set-Content docs/docs.go"

.PHONY: setup
setup: download wait-db migrate-up gen-docs
	@echo "setup completed."

.PHONY: reset-db
reset-db: migrate-drop migrate-up seed
	@echo "databse reseted."
