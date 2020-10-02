# fiit-pdt

[![golangci-lint](https://github.com/timzatko/fiit-pdt/workflows/golangci-lint/badge.svg)](https://github.com/timzatko/fiit-pdt/actions?query=workflow:golangci-lint+branch:master)

# Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go 1.14](https://golang.org/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) - a CLI tool for working with migrations 

## Setup

1. Navigate to the root of the repository and load environment variables by running - `source .env` [[reference]](https://gist.github.com/mihow/9c7f559807069a03e302605691f85572#gistcomment-2721971)
2. Run PostgreSQL via docker-compose. Database will be running on `localhost:5432` and the adminer will be available at `http://localhost:8080`. You can also run it in detached mode by providing `-d` flag (useful when doing bulk import). 
```bash
docker-compose up
```
3. Run migrations by navigating to the migrations' directory (`cd migrations`) and executing `make init`.

### Data

You need to download data to [data](./data) directory, download link is in [data/README.txt](./data/README.txt).
The smaller subset of data is directly in the repository in [data/examples](./data/examples).