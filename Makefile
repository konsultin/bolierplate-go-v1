SHELL := /bin/sh

COMPOSE_FILE := deployment/docker/docker-compose.yaml
COMPOSE := docker compose -f $(COMPOSE_FILE)

MIGRATIONS_DIR := migrations
MIGRATE ?= $(shell go env GOPATH)/bin/migrate

-include .env
export

.PHONY: init db-up db-down db-script db-version bs

init:
	@read -p "Nama project (tanpa spasi): " NAME; \
	if [ -z "$$NAME" ]; then echo "Nama project wajib diisi"; exit 1; fi; \
	MODULE="github.com/Konsultin/$$NAME"; \
	echo "Mengatur module ke $$MODULE"; \
	go mod edit -module "$$MODULE"; \
	find . -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' -o -name '*.md' -o -name '*.env' \) -not -path './.git/*' -not -path './vendor/*' -print0 | xargs -0 sed -i "s#github.com/Konsultin/project-goes-here#$$MODULE#g"; \
	echo "Menjalankan go mod tidy"; \
	go mod tidy; \
	echo "Menginstall migrate CLI"; \
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	if [ -f .env.example ]; then cp .env.example .env; else echo ".env.example tidak ditemukan, lewati copy"; fi; \
	echo "Inisiasi selesai"

db-url = \
if [ -z "$$DB_DRIVER" ]; then echo "DB_DRIVER tidak diset"; exit 1; fi; \
case "$$DB_DRIVER" in \
	mysql|mariadb) \
		echo "mysql://$${DB_USERNAME}:$${DB_PASSWORD}@tcp($${DB_HOST}:$${DB_PORT})/$${DB_NAME}?parseTime=true";; \
	postgres|postgresql|pg) \
		echo "postgres://$${DB_USERNAME}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=disable";; \
	*) \
		echo "unknown DB_DRIVER '$$DB_DRIVER' (mysql/postgres)"; exit 1;; \
esac

db-up:
	@DB_URL=$$($(db-url)); \
	echo "Migrating up with $$DB_URL"; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" up

db-down:
	@DB_URL=$$($(db-url)); \
	echo "Migrating down 1 step with $$DB_URL"; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" down 1

db-script:
	@read -p "Migration Name: " NAME; \
	if [ -z "$$NAME" ]; then echo "Migration Name must be set"; exit 1; fi; \
	mkdir -p "$(MIGRATIONS_DIR)"; \
	echo "Make Migration $$NAME in $(MIGRATIONS_DIR)"; \
	"$(MIGRATE)" create -ext sql -dir "$(MIGRATIONS_DIR)" -seq "$$NAME"

db-version:
	@DB_URL=$$($(db-url)); \
	read -p "Migration Target: " VERSION; \
	if [ -z "$$VERSION" ]; then echo "Versi wajib diisi"; exit 1; fi; \
	echo "Change Migration Version to $$VERSION with $$DB_URL"; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" goto "$$VERSION"

bs:
	@profile="mysql"; \
	if [ "$${DB_DRIVER}" = "postgres" ] || [ "$${DB_DRIVER}" = "postgresql" ] || [ "$${DB_DRIVER}" = "pg" ]; then \
		profile="postgres"; \
	fi; \
	echo "DB_DRIVER=$${DB_DRIVER:-unset} -> running profile '$$profile'"; \
	$(COMPOSE) --profile $$profile up -d
