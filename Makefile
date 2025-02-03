.PHONY: default build test run docs clean

APP_NAME=gobid

default: watch

watch:
	air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api"

postgres:
	docker compose up -d

migrations:
	go run ./cmd/terndotenv

sql:
	sqlc generate -f ./internal/store/pgstore/sqlc.yml
	