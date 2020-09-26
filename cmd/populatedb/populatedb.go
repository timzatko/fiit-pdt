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

	"github.com/timzatko/fiit-pdt/internal/database"
)

func main() {
	// connect to the database
	db := database.Connect()
	defer database.Close(db)

	// get files in the data directory
	dataDir := path.Join("data")
	files := getFiles(dataDir)

	// start reading the files
	fmt.Print("reading files...")
	fmt.Print(files)

	readFiles(dataDir, files)
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
		if re.Match([]byte(fileName)) {
			fileNames = append(fileNames, fileName)
		}
	}
	return fileNames
}

func readFiles(dataDir string, files []string) {
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go readFile(&wg, dataDir, file)
	}

	wg.Wait()
}

func readFile(wg *sync.WaitGroup, dataDir string, fileName string) {
	defer wg.Done()
	var jsonWg sync.WaitGroup

	log.Printf("reading file %s...", fileName)

	var err error
	file, err := os.Open(path.Join(dataDir, fileName))

	if err != nil {
		log.Fatal(err)
	}

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

	defer func() {
		err := gz.Close()

		if err != nil {
			log.Panicf("error while closing the file (%s): %s", fileName, err)
		}
	}()

	s := bufio.NewScanner(gz)
	for s.Scan() {
		wg.Add(1)
		go parseJson([]byte(s.Text()))
	}

	if err := s.Err(); err != nil {
		log.Panic(err)
	}

	jsonWg.Wait()
	log.Printf("done %s...", fileName)
}

func parseJson(j []byte) {
	var t Tweet
	err := json.Unmarshal(j, &t)

	if err != nil {
		log.Panicf("unable to unmarshal json: %s", err)
	}

	// log.Print(t)
}

// Tweet interface generated with https://mholt.github.io/json-to-go/
type Tweet struct {
	CreatedAt        string `json:"created_at"`
	ID               int64  `json:"id"`
	IDStr            string `json:"id_str"`
	FullText         string `json:"full_text"`
	Truncated        bool   `json:"truncated"`
	DisplayTextRange []int  `json:"display_text_range"`
	Entities         struct {
		Hashtags     []interface{} `json:"hashtags"`
		Symbols      []interface{} `json:"symbols"`
		UserMentions []interface{} `json:"user_mentions"`
		Urls         []struct {
			URL         string `json:"url"`
			ExpandedURL string `json:"expanded_url"`
			DisplayURL  string `json:"display_url"`
			Indices     []int  `json:"indices"`
		} `json:"urls"`
	} `json:"entities"`
	Source               string      `json:"source"`
	InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
	User                 User        `json:"user"`
	Geo                  interface{} `json:"geo"`
	Coordinates          interface{} `json:"coordinates"`
	Place                interface{} `json:"place"`
	Contributors         interface{} `json:"contributors"`
	IsQuoteStatus        bool        `json:"is_quote_status"`
	RetweetCount         int         `json:"retweet_count"`
	FavoriteCount        int         `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	Retweeted            bool        `json:"retweeted"`
	PossiblySensitive    bool        `json:"possibly_sensitive"`
	Lang                 string      `json:"lang"`
}

// User interface generated with https://mholt.github.io/json-to-go/
type User struct {
	ID          int64  `json:"id"`
	IDStr       string `json:"id_str"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Entities    struct {
		URL struct {
			Urls []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"url"`
		Description struct {
			Urls []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	Protected                      bool        `json:"protected"`
	FollowersCount                 int         `json:"followers_count"`
	FriendsCount                   int         `json:"friends_count"`
	ListedCount                    int         `json:"listed_count"`
	CreatedAt                      string      `json:"created_at"`
	FavouritesCount                int         `json:"favourites_count"`
	UtcOffset                      interface{} `json:"utc_offset"`
	TimeZone                       interface{} `json:"time_zone"`
	GeoEnabled                     bool        `json:"geo_enabled"`
	Verified                       bool        `json:"verified"`
	StatusesCount                  int         `json:"statuses_count"`
	Lang                           interface{} `json:"lang"`
	ContributorsEnabled            bool        `json:"contributors_enabled"`
	IsTranslator                   bool        `json:"is_translator"`
	IsTranslationEnabled           bool        `json:"is_translation_enabled"`
	ProfileBackgroundColor         string      `json:"profile_background_color"`
	ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool        `json:"profile_background_tile"`
	ProfileImageURL                string      `json:"profile_image_url"`
	ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
	ProfileBannerURL               string      `json:"profile_banner_url"`
	ProfileImageExtensionsAltText  interface{} `json:"profile_image_extensions_alt_text"`
	ProfileBannerExtensionsAltText interface{} `json:"profile_banner_extensions_alt_text"`
	ProfileLinkColor               string      `json:"profile_link_color"`
	ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string      `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	HasExtendedProfile             bool        `json:"has_extended_profile"`
	DefaultProfile                 bool        `json:"default_profile"`
	DefaultProfileImage            bool        `json:"default_profile_image"`
	Following                      bool        `json:"following"`
	FollowRequestSent              bool        `json:"follow_request_sent"`
	Notifications                  bool        `json:"notifications"`
	TranslatorType                 string      `json:"translator_type"`
}
