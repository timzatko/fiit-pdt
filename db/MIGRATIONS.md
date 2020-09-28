# Migrations

If you would like to create migration, use the following command. Note that you need to have [golang-migrate](https://github.com/golang-migrate) installed.

```bash
migrate create -ext sql -dir db/migrations create_users_table
```

A tutorial on migrations is available [here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md).