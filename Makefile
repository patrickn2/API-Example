Command := $(firstword $(MAKECMDGOALS))
Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

include .env

test:
	go test ./..

run:
	go run ./cmd/api/main.go

dev:
	docker compose up

migrations-up:
	go run ./cmd/migrations/main.go -up

migrations-down:
	go run ./cmd/migrations/main.go -down

migrations-create:
	go run ./cmd/migrations/main.go -create '$(Arguments)'

%::
	@echo ....