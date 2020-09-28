package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sync"

	"gorm.io/gorm"

	"github.com/timzatko/fiit-pdt/cmd/populatedb/utils"
	"github.com/timzatko/fiit-pdt/internal/database"
)

func main() {
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

	readFiles(db, dataDir, fileNames)
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

func readFiles(db *gorm.DB, dataDir string, files []string) {
	var wg sync.WaitGroup
	var synchronizer utils.Synchronizer

	for _, file := range files {
		wg.Add(1)
		// read file and insert to database in goroutine
		go readFile(db, &wg, &synchronizer, dataDir, file)
		// TODO: remove this to go through all files
		// break
	}

	// wait until all files are read
	wg.Wait()
}

func readFile(db *gorm.DB, wg *sync.WaitGroup, synchronizer *utils.Synchronizer, dataDir string, fileName string) {
	defer wg.Done()

	log.Printf("reading file %s...", fileName)

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
	q := utils.NewQueue(db, synchronizer)

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

	// wait until all entities in queue are inserted to database
	q.Wg.Wait()

	if err := s.Err(); err != nil {
		log.Panic(err)
	}

	log.Printf("done %s...", fileName)
}

func parseJson(j []byte) (utils.RawTweet, error) {
	var rt utils.RawTweet
	err := json.Unmarshal(j, &rt)
	return rt, err
}
