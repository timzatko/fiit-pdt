include .env
export

help:
	@echo "lint                      to run golangci-lint"
	@echo "lint-fix                  to fix lint errors"

lint:
	golangci-lint run ./...

lint-fix:
	goimports --local github.com/timzatko/fiit-pdt -w .

db-seed:
	go run cmd/populatedb/main.go

to-elastic:
	go run cmd/toelastic/main.go

db-migrate-up:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations up

db-migrate-up-1:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations up 1

db-migrate-up-6:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations up 6

db-migrate-down-1:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations down 1

db-migrate-down:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations down

db-migrate-drop:
	migrate -database "postgres://${PDT_DATABASE_USER}:${PDT_DATABASE_PASSWORD}@localhost:5432/${PDT_DATABASE_DB}?sslmode=disable" -path db/migrations drop

elastic-tweets-put-index:
	sh ./elastic/scripts/index.sh

elastic-tweets-put-settings:
	sh ./elastic/scripts/settings.sh

elastic-tweets-put-mapping:
	sh ./elastic/scripts/mapping.sh
