package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"gorm.io/gorm"

	"github.com/timzatko/fiit-pdt/cmd/populatedb/utils"
	"github.com/timzatko/fiit-pdt/internal/database"
	"github.com/timzatko/fiit-pdt/internal/timer"
)

func main() {
	logLevel := 0

	// connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("error while openning database connection: %s", err)
	}
	// close the database after everything is done
	defer database.Close(db)

	// get files in the data directory
	dataDir := path.Join("data")
	fileNames := getFiles(dataDir)

	// start reading the files
	fmt.Print("reading files...")
	fmt.Print(fileNames)

	readFiles(db, dataDir, fileNames, logLevel)
}

func getFiles(dataDir string) []string {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Panicf("failed to read the data directory: %s", err)
	}

	re := regexp.MustCompile(`^.+\.jsonl\.gz$`)
	var fileNames []string
	for _, file := range files {
		fileName := file.Name()
		// files must be of type .jsonl.gz
		if re.Match([]byte(fileName)) {
			fileNames = append(fileNames, fileName)
		}
	}
	return fileNames
}

func readFiles(db *gorm.DB, dataDir string, files []string, logLevel int) {
	ctx := context.TODO()
	synchronizer := utils.NewSynchronizer(ctx, 16)

	for i, file := range files {
		readFile(db, &synchronizer, dataDir, file, i+1, len(files), logLevel)
	}

	synchronizer.Wait()
}

func readFile(db *gorm.DB, synchronizer *utils.Synchronizer, dataDir string, fileName string, i int, c int, logLevel int) {
	log.Printf("reading file %s (%d/%d)...", fileName, i, c)
	defer timer.Duration(timer.Track(fmt.Sprintf("done reading %s (%d/%d)...", fileName, i, c)))

	var err error
	file, err := os.Open(path.Join(dataDir, fileName))

	if err != nil {
		log.Fatal(err)
	}

	// handle error while closing file
	defer func() {
		err := file.Close()

		if err != nil {
			log.Panicf("error while closing the file (%s): %s", fileName, err)
		}
	}()

	gz, err := gzip.NewReader(file)
	if err != nil {
		log.Panicf("error while reading the file (%s): %s", fileName, err)
	}

	// handle error while closing the gzip reader
	defer func() {
		err := gz.Close()

		if err != nil {
			log.Panicf("error while closing the gzip file reader (%s): %s", fileName, err)
		}
	}()

	// create a new queue for entities
	q := utils.NewQueue(db, synchronizer, logLevel)

	s := bufio.NewScanner(gz)
	for s.Scan() {
		j := s.Text()
		// parse the one line which represents one tweet
		rawTweet, err := parseJson([]byte(j))

		if err != nil {
			log.Printf("unable to unmarshal %s", j)
			continue
		}

		// add raw tweet to queue, so it will be inserted later in bulk
		q.Enqueue(&rawTweet)

		// if queue is full, flush it and insert entities to the database
		if q.IsFull() {
			q.Flush()
		}
	}

	// if queue is not empty, flush it and insert remaining entities to the database
	if !q.IsEmpty() {
		q.Flush()
	}

	if err := s.Err(); err != nil {
		log.Panic(err)
	}
}

func parseJson(j []byte) (utils.RawTweet, error) {
	var rt utils.RawTweet
	err := json.Unmarshal(j, &rt)
	return rt, err
}
