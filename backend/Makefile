-include .env
export

BINARY_NAME=schedule-api
MIGRATE=go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	-path ./migrations \
	-database "$(DATABASE_URL)"

.PHONY: run build migrate migrate-down migrate-drop migrate-create seed deps lint test

run:
	go run ./cmd/api

build:
	go build -o bin/$(BINARY_NAME) ./cmd/api

deps:
	go mod tidy

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-drop:
	$(MIGRATE) drop -f

migrate-create:
	@read -p "Migration name: " name; \
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		create -ext sql -dir ./migrations -seq $$name

seed:
	go run ./cmd/seed

test:
	go test ./...

lint:
	golangci-lint run ./...
