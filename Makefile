# Create database
.PHONY: createdb 
createdb:
	psql -h localhost -U postgres \
        -c "CREATE DATABASE filmoteka ENCODING='UTF8'";

# Drop database
.PHONY: dropdb 
dropdb:
	psql -h localhost -U postgres \
        -c "DROP DATABASE IF EXISTS filmoteka";

migrateup:
	migrate -database "postgres://postgres:postgres@localhost:5432/filmoteka?sslmode=disable" -path db/migrations up

migratedown:
	migrate -database "postgres://postgres:postgres@localhost:5432/filmoteka?sslmode=disable" -path db/migrations down
	
.PHONY: build
build:
	go build -o bin/filmoteka -v ./cmd/filmoteka-rest-api

.PHONY: run 
run:
	bin/filmoteka

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build