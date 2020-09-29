package utils

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/timzatko/fiit-pdt/internal/app/model"
)

type Queue struct {
	rts  [2500]*RawTweet
	size int
	Wg   sync.WaitGroup
	db   *gorm.DB
	sync *Synchronizer
}

func NewQueue(db *gorm.DB, synchronizer *Synchronizer) Queue {
	return Queue{
		db:   db,
		sync: synchronizer,
	}
}

func (q *Queue) IsFull() bool {
	return len(q.rts) == q.size
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue) Flush() {
	if !q.IsFull() {
		return
	}

	q.Wg.Add(1)

	rts, size := q.clear()
	go q.process(&rts, size)
}

// Enqueue an entity to its own queue
func (q *Queue) Enqueue(rt *RawTweet) {
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

func (q *Queue) process(rts *[2500]*RawTweet, size int) {
	defer q.Wg.Done()

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
		q.insert(&accs, &q.sync.AccountsMutex)
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
		q.insert(&hts, &q.sync.HashtagsMutex)
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
		q.insert(&countries, &q.sync.CountriesMutex)
	}

	// TWEETS
	q.sync.TweetsMutex.Lock()
	q.sync.CountriesMutex.Lock() // lock also countries, since we are selecting country_id

	var tweets []map[string]interface{}

	for i := 0; i < size; i++ {
		rt := rts[i]

		var loc interface{} = nil
		if rt.Geo != nil {
			loc = clause.Expr{SQL: "ST_SetSRID(ST_MakePoint(%f, %f), 4326)", Vars: []interface{}{rt.Geo.Coordinates[0], rt.Geo.Coordinates[1]}}
		}

		var cid interface{} = nil
		if rt.Place != nil && len(rt.Place.CountryCode) > 0 {
			cid = clause.Expr{SQL: "(SELECT id FROM countries WHERE name=? LIMIT 1)", Vars: []interface{}{rt.Place.Country}}
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

	// insert to tweets hashtags table
	q.db.Model(model.Tweet{}).Clauses(clause.OnConflict{DoNothing: true}).Create(tweets)
	q.sync.CountriesMutex.Unlock()
	q.sync.TweetsMutex.Unlock()

	// TWEET HASHTAGS
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
		q.sync.TweetHashtagsMutex.Lock()
		q.db.Model(model.TweetHashtag{}).Clauses(clause.OnConflict{DoNothing: true}).Create(ths)
		q.sync.TweetHashtagsMutex.Unlock()
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
		q.insert(&tms, &q.sync.TweetMentionsMutex)
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
func (q *Queue) clear() ([2500]*RawTweet, int) {
	rts := q.rts
	size := q.size

	q.rts = [2500]*RawTweet{}
	q.size = 0

	return rts, size
}

func toTime(rt *RawTweet) time.Time {
	t, err := time.Parse(time.RubyDate, rt.CreatedAt)
	if err != nil {
		log.Panicf("error while parsing time: %s", err)
	}
	return t
}
