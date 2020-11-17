package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/timzatko/fiit-pdt/internal/app/model"
	"github.com/timzatko/fiit-pdt/internal/reader"
	"github.com/timzatko/fiit-pdt/internal/synchronizer"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Author struct {
	Id             int64  `json:"id"`
	ScreenName     string `json:"screen_name"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	StatusesCount  int    `json:"statuses_count"`
}

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Mention struct {
	Id         int64  `json:"id"`
	ScreenName string `json:"screen_name"`
	Name       string `json:"name"`
}

type Tweet struct {
	Id            string    `json:"id"`
	Content       string    `json:"content"`
	Location      *Location `json:"location"`
	RetweetCount  int64     `json:"retweet_count"`
	FavoriteCount int64     `json:"favorite_count"`
	HappenedAt    int64     `json:"happened_at"`
	Author        Author    `json:"author"`
	Country       *Country  `json:"country"`
	Hashtags      []string  `json:"hashtags"`
	Mentions      []Mention `json:"mentions"`
	ParentId      *string   `json:"parent_id"`
}

type Queue struct {
	rts      []model.RawTweet
	http     *http.Client
	maxSize  int
	sync     *synchronizer.Synchronizer
	counter  int
	logLevel int
}

type Response struct {
	Took   int           `json:"took"`
	Errors bool          `json:"errors"`
	Items  []interface{} `json:"items"`
}

func NewQueue(sync *synchronizer.Synchronizer, http *http.Client, maxSize int, logLevel int) reader.Queue {
	return &Queue{
		sync:     sync,
		counter:  0,
		maxSize:  maxSize,
		logLevel: logLevel,
		http:     http,
	}
}

func (q *Queue) IsFull() bool {
	return q.maxSize == len(q.rts)
}

func (q *Queue) IsEmpty() bool {
	return len(q.rts) == 0
}

func (q *Queue) Flush() {
	if err := q.sync.Acquire(); err != nil {
		log.Panicf("failed to acquire semaphore: %v", err)
	}

	q.counter += 1
	rts := q.clear()
	go q.process(rts, q.counter)
}

// Enqueue an entity to its own queue
func (q *Queue) Enqueue(rt *model.RawTweet) {
	if rt.RetweetedStatus != nil {
		// enqueue retweeted tweet
		q.Enqueue(rt.RetweetedStatus)
	}

	if q.IsFull() {
		q.Flush()
	}

	q.rts = append(q.rts, *rt)
}

func (q *Queue) process(rts []model.RawTweet, batchId int) {
	defer q.sync.Release()

	var cmd []string

	for _, rt := range rts {
		var loc *Location = nil

		if rt.Geo != nil {
			loc = &Location{
				Lat: rt.Geo.Coordinates[0],
				Lon: rt.Geo.Coordinates[0],
			}
		}

		var country *Country = nil
		if rt.Place != nil && len(rt.Place.CountryCode) > 0 {
			country = &Country{
				Code: rt.Place.CountryCode,
				Name: rt.Place.Country,
			}
		}

		var parentId *string = nil
		if rt.RetweetedStatus != nil {
			parentId = &rt.RetweetedStatus.IDStr
		}

		author := Author{
			Id:             rt.User.ID,
			ScreenName:     rt.User.ScreenName,
			Name:           rt.User.Name,
			Description:    rt.User.Description,
			FollowersCount: rt.User.FavouritesCount,
			FriendsCount:   rt.User.FriendsCount,
			StatusesCount:  rt.User.StatusesCount,
		}

		var hashtags []string
		for _, h := range rt.Entities.Hashtags {
			hashtags = append(hashtags, h.Text)

		}

		var mentions []Mention
		for _, m := range rt.Entities.UserMentions {
			mentions = append(mentions, Mention{
				Id:         m.ID,
				ScreenName: m.ScreenName,
				Name:       m.Name,
			})
		}

		tw := Tweet{
			Id:            rt.IDStr,
			Content:       rt.FullText,
			RetweetCount:  rt.RetweetCount,
			FavoriteCount: rt.FavoriteCount,
			HappenedAt:    toTime(rt).Unix(),
			Location:      loc,
			Country:       country,
			Author:        author,
			ParentId:      parentId,
			Hashtags:      hashtags,
			Mentions:      mentions,
		}

		b, err := json.Marshal(tw)
		if err != nil {
			fmt.Printf("failed to unmarshal tweet: %v\n", err)
			continue
		}

		cmd = append(cmd, `{ "index": { } }`, string(b))
	}

	var body string
	for i := 0; i < len(cmd); i++ {
		body += cmd[i] + "\n"
	}

	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:9200/tweets/_bulk", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Printf("failed to create request for batch #%d: %v\n", batchId, err)
	}
	req.Header.Set("Content-Type", "application/json")

	var r *http.Response
	r, err = q.http.Do(req)
	if err != nil {
		fmt.Printf("failed to do request for batch #%d: %v\n", batchId, err)
	}

	var rBodyBytes []byte
	rBodyBytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("failed to read all: %v\n", err)
	}

	if r.StatusCode == 200 {
		var rBody Response
		err = json.Unmarshal(rBodyBytes, &rBody)
		if err != nil {
			fmt.Printf("failed to unmarshall: %v\n", err)
		}

		if rBody.Errors == true {
			fmt.Printf("bulk failed: %v\n", rBody.Items)
		}
	} else {
		fmt.Printf("bulk failed with status=%d: %v\n", r.StatusCode, rBodyBytes)
	}

	_ = r.Body.Close()

	if q.logLevel <= 1 {
		fmt.Printf("batch #%d request sent...\n", batchId)
	}
}

// clear the queue and return the old queue
func (q *Queue) clear() []model.RawTweet {
	rts := q.rts
	q.rts = []model.RawTweet{}
	return rts
}

func toTime(rt model.RawTweet) time.Time {
	t, err := time.Parse(time.RubyDate, rt.CreatedAt)
	if err != nil {
		log.Panicf("error while parsing time: %s", err)
	}
	return t
}
