package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/timzatko/fiit-pdt/internal/app/model"
	"github.com/timzatko/fiit-pdt/internal/reader"
	"github.com/timzatko/fiit-pdt/internal/synchronizer"
	"github.com/timzatko/fiit-pdt/internal/timer"
)

type Synchronizer struct {
	AccountsMutex      sync.Mutex
	HashtagsMutex      sync.Mutex
	CountriesMutex     sync.Mutex
	TweetsMutex        sync.Mutex
	TweetHashtagsMutex sync.Mutex
	TweetMentionsMutex sync.Mutex
}

type Queue struct {
	// when using higher values than 5000 (eg. 10000)
	// the tweets won't be imported
	rts      [500]*model.RawTweet
	size     int
	db       *gorm.DB
	sync     *synchronizer.Synchronizer
	mu       *Synchronizer
	counter  int
	logLevel int
}

func NewQueue(db *gorm.DB, synchronizer *synchronizer.Synchronizer, logLevel int) reader.Queue {
	return &Queue{
		db:       db,
		sync:     synchronizer,
		counter:  0,
		logLevel: logLevel,
		mu:       &Synchronizer{},
	}
}

func (q *Queue) IsFull() bool {
	return len(q.rts) == q.size
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue) Flush() {
	if err := q.sync.Acquire(); err != nil {
		log.Panicf("failed to acquire semaphore: %v", err)
	}

	q.counter += 1
	rts, size := q.clear()
	go q.process(&rts, q.counter, size)
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

	q.rts[q.size] = rt
	q.size += 1
}

func (q *Queue) process(rts *[500]*model.RawTweet, batchId int, size int) {
	defer q.sync.Release()

	if q.logLevel > 1 {
		log.Printf("processing batch #%d with %d tweets...", batchId, size)
		defer timer.Duration(timer.Track(fmt.Sprintf("batch %d processed!", batchId)))
	}

	// ACCOUNTS
	var accs []model.Account

	for i := 0; i < size; i++ {
		rt := rts[i]

		accs = append(accs, model.Account{
			Id:             rt.User.ID,
			ScreenName:     rt.User.ScreenName,
			Name:           rt.User.Name,
			Description:    rt.User.Description,
			FollowersCount: rt.User.FollowersCount,
			FriendsCount:   rt.User.FriendsCount,
			StatusesCount:  rt.User.StatusesCount,
		})

		// add also users from user mentions
		for _, um := range rt.Entities.UserMentions {
			accs = append(accs, model.Account{
				Id:         um.ID,
				Name:       um.Name,
				ScreenName: um.ScreenName,
			})
		}
	}
	// insert to accounts table
	if len(accs) > 0 {
		if q.logLevel > 1 {
			log.Printf("batch #%d inserting %d accounts...", batchId, len(accs))
		}
		res := q.insert(&accs, &q.mu.AccountsMutex)
		if res.Error != nil {
			log.Panicf("error: batch #%d unable to insert accounts: %s", batchId, res.Error)
		}
	}

	// HASHTAGS
	var hts []model.Hashtag
	for i := 0; i < size; i++ {
		rt := rts[i]

		for _, rawHashtag := range rt.Entities.Hashtags {
			ht := model.Hashtag{
				Value: rawHashtag.Text,
			}
			hts = append(hts, ht)
		}
	}

	// insert to hashtags table
	if len(hts) > 0 {
		if q.logLevel > 1 {
			log.Printf("batch #%d inserting %d hashtags...", batchId, len(hts))
		}
		res := q.insert(&hts, &q.mu.HashtagsMutex)
		if res.Error != nil {
			log.Panicf("error: batch #%d unable to insert hashtags: %s", batchId, res.Error)
		}
	}

	// COUNTRIES
	var countries []model.Country
	for i := 0; i < size; i++ {
		rt := rts[i]

		if rt.Place != nil && len(rt.Place.CountryCode) > 0 {
			c := model.Country{
				Code: rt.Place.CountryCode,
				Name: rt.Place.Country,
			}

			countries = append(countries, c)
		}
	}

	// insert to countries database
	if len(countries) > 0 {
		if q.logLevel > 1 {
			log.Printf("batch #%d inserting %d countries...", batchId, len(countries))
		}
		res := q.insert(&countries, &q.mu.CountriesMutex)
		if res.Error != nil {
			log.Panicf("error: batch #%d unable to insert countries: %s", batchId, res.Error)
		}
	}

	// TWEETS
	q.mu.TweetsMutex.Lock()
	var tweets []map[string]interface{}

	for i := 0; i < size; i++ {
		rt := rts[i]

		var loc interface{} = nil
		if rt.Geo != nil {
			loc = clause.Expr{SQL: "ST_SetSRID(ST_MakePoint(?, ?), 4326)", Vars: []interface{}{rt.Geo.Coordinates[0], rt.Geo.Coordinates[1]}}
		}

		var cid interface{} = nil
		if rt.Place != nil && len(rt.Place.CountryCode) > 0 {
			cid = clause.Expr{SQL: "(SELECT id FROM countries WHERE code=? LIMIT 1)", Vars: []interface{}{rt.Place.CountryCode}}
		}

		var pid interface{} = nil
		if rt.RetweetedStatus != nil {
			pid = rt.RetweetedStatus.IDStr
		}

		tweets = append(tweets, map[string]interface{}{
			"Id":            rt.IDStr,
			"Content":       rt.FullText,
			"RetweetCount":  rt.RetweetCount,
			"FavoriteCount": rt.FavoriteCount,
			"AuthorId":      rt.User.IDStr,
			"HappenedAt":    clause.Expr{SQL: "to_timestamp(?)", Vars: []interface{}{toTime(rt).Unix()}},
			"CountryId":     cid,
			"Location":      loc,
			"ParentId":      pid,
		})
	}

	// insert to tweets table
	if q.logLevel > 1 {
		log.Printf("batch #%d inserting %d tweets...", batchId, len(tweets))
	}
	res := q.db.Model(model.Tweet{}).Clauses(clause.OnConflict{DoNothing: true}).Create(tweets)
	q.mu.TweetsMutex.Unlock()
	if res.Error != nil {
		log.Panicf("error: batch #%d unable to insert tweets: %s", batchId, res.Error)
	}

	//// TWEET HASHTAGS
	var ths []map[string]interface{}

	for i := 0; i < size; i++ {
		rt := rts[i]

		for _, h := range rt.Entities.Hashtags {
			ths = append(ths, map[string]interface{}{
				"HashtagId": clause.Expr{SQL: "(SELECT id FROM hashtags WHERE value=? LIMIT 1)", Vars: []interface{}{h.Text}},
				"TweetId":   rt.IDStr,
			})
		}
	}

	if len(ths) > 0 {
		q.mu.TweetHashtagsMutex.Lock()
		if q.logLevel > 1 {
			log.Printf("batch #%d inserting %d tweet hashtags...", batchId, len(ths))
		}
		res := q.db.Model(model.TweetHashtag{}).Clauses(clause.OnConflict{DoNothing: true}).Create(ths)
		q.mu.TweetHashtagsMutex.Unlock()
		if res.Error != nil {
			log.Panicf("error: batch #%d unable to insert tweet hashtags: %s", batchId, res.Error)
		}
	}

	// TWEET MENTIONS
	var tms []model.TweetMention
	for i := 0; i < size; i++ {
		rt := rts[i]

		for _, um := range rt.Entities.UserMentions {
			tm := model.TweetMention{
				TweetId:   rt.IDStr,
				AccountId: um.ID,
			}
			tms = append(tms, tm)
		}
	}

	// insert to tweet mentions table
	if len(tms) > 0 {
		if q.logLevel > 1 {
			log.Printf("batch #%d inserting %d tweet mentions...", batchId, len(tms))
		}
		res := q.insert(&tms, &q.mu.TweetMentionsMutex)
		if res.Error != nil {
			log.Panicf("error: batch #%d unable to insert tweet mentions: %s", batchId, res.Error)
		}
	}
}

// insert the entities to the database in batch command
func (q *Queue) insert(entities interface{}, mutex *sync.Mutex) *gorm.DB {
	mutex.Lock()
	res := q.db.Clauses(clause.OnConflict{DoNothing: true}).Create(entities)
	mutex.Unlock()

	return res
}

// clear the queue and return the old queue
func (q *Queue) clear() ([500]*model.RawTweet, int) {
	rts := q.rts
	size := q.size

	q.rts = [500]*model.RawTweet{}
	q.size = 0

	return rts, size
}

func toTime(rt *model.RawTweet) time.Time {
	t, err := time.Parse(time.RubyDate, rt.CreatedAt)
	if err != nil {
		log.Panicf("error while parsing time: %s", err)
	}
	return t
}
