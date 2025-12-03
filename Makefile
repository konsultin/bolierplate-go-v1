SHELL := /bin/sh

COMPOSE_FILE := deployment/docker/docker-compose.yaml
COMPOSE := docker compose -f $(COMPOSE_FILE)

MIGRATIONS_DIR := migrations
MIGRATE ?= $(shell go env GOPATH)/bin/migrate

-include .env
export

.PHONY: setup-project init db-up db-down db-script db-version bs

setup-project:
	@read -p "Project name (no spaces): " NAME; \
	if [ -z "$$NAME" ]; then echo "Project name is required"; exit 1; fi; \
	MODULE="github.com/Konsultin/$$NAME"; \
	echo "Setting module to $$MODULE"; \
	go mod edit -module "$$MODULE"; \
	find . -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' -o -name '*.md' -o -name '*.env' \) -not -path './.git/*' -not -path './vendor/*' -print0 | xargs -0 sed -i "s#github.com/Konsultin/project-goes-here#$$MODULE#g"; \
	if [ -f .env.example ] && [ ! -f .env ]; then \
		echo "Creating .env from .env.example"; \
		cp .env.example .env; \
	fi; \
	if [ -f .env ]; then \
		if grep -q '^COMPOSE_PROJECT_NAME=' .env; then \
			sed -i "s/^COMPOSE_PROJECT_NAME=.*/COMPOSE_PROJECT_NAME=$$NAME/" .env; \
		else \
			echo "COMPOSE_PROJECT_NAME=$$NAME" >> .env; \
		fi; \
		if grep -q '^LOG_NAMESPACE=' .env; then \
			sed -i "s/^LOG_NAMESPACE=.*/LOG_NAMESPACE=$$NAME/" .env; \
		else \
			echo "LOG_NAMESPACE=$$NAME" >> .env; \
		fi; \
		echo "Set COMPOSE_PROJECT_NAME=$$NAME for container names"; \
		echo "Set LOG_NAMESPACE=$$NAME"; \
	else \
		echo ".env not found; skipped updating COMPOSE_PROJECT_NAME/LOG_NAMESPACE"; \
	fi; \
	echo "Project setup completed. Run 'make init' to finish initialization"; \
	printf '\033[0;31m%s\033[0m\n' "Reminder: setup environment first (e.g. update .env) before running services"

init:
	@if grep -q '^module github.com/Konsultin/project-goes-here' go.mod; then \
		echo "Project name is still the default. Run 'make setup-project' first."; \
		exit 1; \
	fi; \
	if [ -f .env.example ] && [ ! -f .env ]; then \
		echo "Creating .env from .env.example"; \
		cp .env.example .env; \
	fi; \
	echo "Running go mod tidy"; \
	go mod tidy; \
	echo "Installing migrate CLI"; \
	go install -tags "postgres,mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	echo "Initialization completed"; \
	printf '\033[0;31m%s\033[0m\n' "Reminder: setup environment first (e.g. update .env) before running services"

db-url = \
if [ -z "$$DB_DRIVER" ]; then echo "DB_DRIVER is not set"; exit 1; fi; \
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
	echo "Migrating up..."; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" up

db-down:
	@DB_URL=$$($(db-url)); \
	echo "Migrating down 1 step..."; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" down 1

db-script:
	@read -p "Migration Name: " NAME; \
	if [ -z "$$NAME" ]; then echo "Migration Name must be set"; exit 1; fi; \
	mkdir -p "$(MIGRATIONS_DIR)"; \
	echo "Make Migration $$NAME in $(MIGRATIONS_DIR)"; \
	"$(MIGRATE)" create -ext sql -dir "$(MIGRATIONS_DIR)" -seq "$$NAME"

db-version:
	@DB_URL=$$($(db-url)); \
	read -p "Target Migration Version: " VERSION; \
	if [ -z "$$VERSION" ]; then echo "Version is required"; exit 1; fi; \
	echo "Change Migration version to $$VERSION"; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DB_URL" goto "$$VERSION"

bs:
	@profile="mysql"; \
	if [ "$${DB_DRIVER}" = "postgres" ] || [ "$${DB_DRIVER}" = "postgresql" ] || [ "$${DB_DRIVER}" = "pg" ]; then \
		profile="postgres"; \
	fi; \
	echo "DB_DRIVER=$${DB_DRIVER:-unset} -> running profile '$$profile'"; \
	$(COMPOSE) --profile $$profile up -d
