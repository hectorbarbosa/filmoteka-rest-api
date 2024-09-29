# Create local database
.PHONY: createdb 
createdb:
	psql -h localhost -U postgres \
        -c "CREATE DATABASE filmoteka ENCODING='UTF8'";

# Drop local database
.PHONY: dropdb 
dropdb:
	psql -h localhost -U postgres \
        -c "DROP DATABASE IF EXISTS filmoteka";

# Create Docker container `filmoteka` tables
migrateup:
	docker compose run api migrate -path /api/migrations/ -database postgres://user:password@postgres:5432/filmoteka?sslmode=disable up

# Drop Docker container `filmoteka` tables
migratedown:
	docker compose run api migrate -path /api/migrations/ -database postgres://user:password@postgres:5432/filmoteka?sslmode=disable down 
	
.PHONY: build
build:
	go build -o bin/filmoteka -v ./cmd/filmoteka-rest-api

.PHONY: run 
run:
	bin/filmoteka -env ./local.env

.PHONY: swagger 
swagger:
	swag init -d ./cmd/filmoteka-rest-api,./internal/restapi,./internal/app/models,./internal,internal/restapi/models

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build