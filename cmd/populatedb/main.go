package main

import (
	"context"
	"fmt"
	"log"
	"path"

	"gorm.io/gorm"

	"github.com/timzatko/fiit-pdt/cmd/populatedb/utils"
	"github.com/timzatko/fiit-pdt/internal/database"
	"github.com/timzatko/fiit-pdt/internal/reader"
	"github.com/timzatko/fiit-pdt/internal/synchronizer"
)

type Constructor struct {
	db   *gorm.DB
	sync *synchronizer.Synchronizer
}

func (c Constructor) NewQueue(logLevel int) reader.Queue {
	return utils.NewQueue(c.db, c.sync, logLevel)
}

func main() {
	logLevel := 0

	// connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("error while openning database connection: %s", err)
	}
	// close the database after everything is done
	defer database.Close(db)

	ctx := context.TODO()
	sync := synchronizer.NewSynchronizer(ctx, 4)
	cs := Constructor{db: db, sync: &sync}

	// get files in the data directory
	dataDir := path.Join("data")
	fileNames := reader.GetFiles(dataDir)

	// start reading the files
	fmt.Print("reading files...")
	fmt.Print(fileNames)

	reader.ReadFiles(cs, dataDir, fileNames, logLevel)
}
