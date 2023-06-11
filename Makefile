ifneq (,$(wildcard ./.env))
    include .env
    export
endif

db-migration:
	@read -p "Enter migration name:" name;\
		migrate create -ext sql -dir db/migrations $$name

db-migrate:
	migrate -source ./db/migrations -database ${DB_URL} up

db-rollback:
	migrate -source ./db/migrations -database ${DB_URL} down

build:
	@go build -o bin/production-api

run: build
	@./bin/production-api
