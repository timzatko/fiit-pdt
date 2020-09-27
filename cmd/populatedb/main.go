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

	"github.com/go-pg/pg"

	"github.com/timzatko/fiit-pdt/internal/database"
)

// Synchronizer is used to sync inserting entities of different types
// between different files. So at one time, only to one entity is being written.
type Synchronizer struct {
	AccountsMutex sync.Mutex
}

func main() {
	// connect to the database
	db := database.Connect()
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

func readFiles(db *pg.DB, dataDir string, files []string) {
	var wg sync.WaitGroup
	var synchronizer Synchronizer

	for _, file := range files {
		wg.Add(1)
		// read file and insert to database in goroutine
		go readFile(db, &wg, &synchronizer, dataDir, file)
	}

	// wait until all files are read
	wg.Wait()
}

func readFile(db *pg.DB, wg *sync.WaitGroup, synchronizer *Synchronizer, dataDir string, fileName string) {
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
	q := NewQueue(db, synchronizer)

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
		q.add(rawTweet)

		// if queue is full, flush it and insert entities to the database
		if q.isFull() {
			q.send()
		}
	}

	// if queue is not empty, flush it and insert remaining entities to the database
	if !q.isEmpty() {
		q.send()
	}

	// wait until all entities in queue are inserted to database
	q.wg.Wait()

	if err := s.Err(); err != nil {
		log.Panic(err)
	}

	log.Printf("done %s...", fileName)
}

func parseJson(j []byte) (RawTweet, error) {
	var rt RawTweet
	err := json.Unmarshal(j, &rt)
	return rt, err
}
