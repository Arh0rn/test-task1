include .env

MIGRATE=migrate -path migrations -database \
	 "postgres://postgres:${DB_PASSWORD}@localhost:5432/test-task1?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

print.env:
	type .env