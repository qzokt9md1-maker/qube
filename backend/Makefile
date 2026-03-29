.PHONY: dev build run test migrate gqlgen

dev:
	go run ./cmd/server

build:
	go build -o bin/qube ./cmd/server

run: build
	./bin/qube

test:
	go test ./... -v

gqlgen:
	go run github.com/99designs/gqlgen generate

migrate-up:
	docker exec -i qube-postgres psql -U qube -d qube < migrations/001_initial.sql

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-reset:
	docker compose down -v && docker compose up -d
