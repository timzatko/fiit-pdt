help:
	@echo "lint                      to run golangci-lint"
	@echo "lint-fix                  to fix lint errors"
	@echo "populate-db			     populate database with downloaded data from ./data directory"

lint:
	golangci-lint run ./...

lint-fix:
	goimports --local github.com/timzatko/fiit-pdt -w .

populate-db:
	go run cmd/populatedb/populatedb.go
