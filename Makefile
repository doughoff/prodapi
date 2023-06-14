ifneq (,$(wildcard ./.env))
    include .env
    export
endif

db-create-migration:
	@read -p "Enter migration name:" name;\
		migrate create -ext sql -dir postgres/migrations $$name

db-migrate:
	migrate -path postgres/migrations -database ${DB_URL} up

db-rollback:
	migrate -path postgres/migrations -database ${DB_URL} down

build:
	@go build -o bin/production-api

run: build
	@./bin/production-api

container-db-migrate:
	docker run -v ./postgres/migrations:/migrations --network host migrate/migrate:v4.16.1 -path=/migrations/ -database ${DB_URL} up

build-container:
	@docker build -t prodapiv5 .

run-container: build-container
	@docker run -p 3088:3088 --network host --env-file .env prodapiv5

