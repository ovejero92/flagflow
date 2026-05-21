.PHONY: run build test migrate-up migrate-down docker-up docker-down sqlc

DATABASE_URL ?= postgres://flaguser:flagpass@localhost:5432/feature_flags?sslmode=disable
MIGRATIONS_PATH ?= file://migrations

run:
	go run ./cmd/server

build:
	go build -o bin/feature-flag-service ./cmd/server

test:
	go test ./...

migrate-up:
	docker run --rm -v "$(CURDIR)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(DATABASE_URL)" up

migrate-down:
	docker run --rm -v "$(CURDIR)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(DATABASE_URL)" down 1

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

sqlc:
	sqlc generate
