package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	User     string
	Password string
	Db       string
	Port     string `default:"5432"`
}

func Connect() (*gorm.DB, error) {
	var conn Connection
	if err := envconfig.Process("PDT_DATABASE", &conn); err != nil {
		log.Fatal(err.Error())
	}

	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 60 * time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Bratislava", conn.User, conn.Password, conn.Db, conn.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})
}

func Close(db *gorm.DB) {
	sqlDb, err := db.DB()

	if err != nil {
		log.Fatalf("error while closing the database: %s", err)
	}

	err = sqlDb.Close()

	if err != nil {
		log.Fatalf("error while closing the database: %s", err)
	}
}
