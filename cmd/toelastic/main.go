package main

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/timzatko/fiit-pdt/cmd/toelastic/utils"
	"github.com/timzatko/fiit-pdt/internal/reader"
	"github.com/timzatko/fiit-pdt/internal/synchronizer"
)

type Constructor struct {
	sync    *synchronizer.Synchronizer
	http    *http.Client
	maxSize int
}

func (c Constructor) NewQueue(logLevel int) reader.Queue {
	return utils.NewQueue(c.sync, c.http, c.maxSize, logLevel)
}

func main() {
	logLevel := 0
	maxSize := 2500

	ctx := context.TODO()
	sync := synchronizer.NewSynchronizer(ctx, 4)
	cs := Constructor{sync: &sync, http: http.DefaultClient, maxSize: maxSize}
	fmt.Printf("batch size is %d", maxSize)

	// set timeout to 10 minutes
	http.DefaultClient.Timeout = 10 * time.Minute

	// get files in the data directory
	dataDir := path.Join("data")
	fileNames := reader.GetFiles(dataDir)

	// start reading the files
	fmt.Print("reading files...")
	fmt.Print(fileNames)

	reader.ReadFiles(cs, dataDir, fileNames, logLevel)
	sync.Wait()
}
