include .env
export

help:
	@echo "lint                      to run golangci-lint"
	@echo "lint-fix                  to fix lint errors"
	@echo "populate-db			     populate database with downloaded data from ./data directory"

lint:
	golangci-lint run ./...

lint-fix:
	goimports --local github.com/timzatko/fiit-pdt -w .

db-seed:
	go run cmd/populatedb/main.go

db-migrate-up:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations up

db-migrate-up-one:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations up 1

db-migrate-down-one:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations down 1

db-migrate-drop:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations drop
