include .env

# For windows!!!

MIGRATE=migrate -path migrations -database \
	 "postgres://postgres:${DB_PASSWORD}@localhost:5432/test-task1?sslmode=disable"
#Given this just for example, better to replace all in env.

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

print.env:
	type .env

swag:
	swag init -g cmd/app/main.go -o ./docs

run:
	go run cmd/app/main.go

# Build binary
build:
	go build -o bin/app.exe cmd/app/main.go

clean:
	if exist bin rmdir /s /q bin
