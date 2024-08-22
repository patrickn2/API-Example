Command := $(firstword $(MAKECMDGOALS))
Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

include .env

test:
	go test ./...

run-api:
	go run ./cmd/api/main.go

build-api:
	go build ./cmd/api/main.go

dev:
	docker compose up

mocks:
	mockery --all

migrations-up:
	go run ./cmd/migrations/main.go -up

migrations-down:
	go run ./cmd/migrations/main.go -down

migrations-create:
	go run ./cmd/migrations/main.go -create '$(Arguments)'

%::
	@echo ....