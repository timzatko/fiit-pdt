# fiit-pdt

[![golangci-lint](https://github.com/timzatko/fiit-pdt/workflows/golangci-lint/badge.svg)](https://github.com/timzatko/fiit-pdt/actions?query=workflow:golangci-lint+branch:master)

# Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go 1.14](https://golang.org/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) - a CLI tool for working with migrations 

## Week 1 - 7: PostgreSQL & Postgis

### Setup

1. Navigate to the root of the repository and load environment variables by running - `source .env` [[reference]](https://gist.github.com/mihow/9c7f559807069a03e302605691f85572#gistcomment-2721971)
2. Run PostgreSQL via docker-compose. Database will be running on `localhost:5432` and the adminer will be available at `http://localhost:8080`. You can also run it in detached mode by providing `-d` flag (useful when doing bulk import). 
```bash
docker-compose up
```
3. Run migrations by running - `make db-migrate-up`. I recommend running just the first 6 migrations with required tables by running - `make db-migrate-up-6`, then running `make db-seed` and then `make db-migrate-up`. The migration will be much faster because of no constraints.  

#### Data

You need to download data to [data](./data) directory, download link is in [data/README.txt](./data/README.txt).
The smaller subset of data is directly in the repository in [data/examples](./data/examples).
To import data run - `make db-seed`.

## Week 8 - 12: ElasticSearch

### Setup

1. Navigate to [elastic](./elastic) folder - `cd elastic`
2. Run the elastic search docker image - `docker-compose up`
3. Create tweets index and setup mapping - `make elastic-tweets-put-index && elastic-tweets-put-mapping`
4. Import data - `make to-elastic` (requires tweets in [data](./data) directory)
