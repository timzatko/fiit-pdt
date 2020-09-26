# fiit-pdt

[![golangci-lint](https://github.com/timzatko/fiit-pdt/workflows/golangci-lint/badge.svg)](https://github.com/timzatko/fiit-pdt/actions?query=workflow:golangci-lint+branch:master)

# Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go 1.14](https://golang.org/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Setup

1. Run PostgreSQL via docker-compose. Database will be running on `localhost:5432` and the adminer will be available at `http://localhost:8080`.

```bash
docker-compose up
```
2. Run migrations by navigating to the migrations directory (`cd migrations`) and executing `make init`.
