package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec("")
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("")
		return err
	})
}
