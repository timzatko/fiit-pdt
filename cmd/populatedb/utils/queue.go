package utils

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/timzatko/fiit-pdt/internal/app/model"
)

type Queue struct {
	rts  []*RawTweet
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
	return len(q.rts) >= 100
}

func (q *Queue) IsEmpty() bool {
	return len(q.rts) == 0
}

func (q *Queue) Flush() {
	if !q.IsFull() {
		return
	}

	q.Wg.Add(1)

	rts := q.clear()
	go q.process(rts)
}

// Enqueue an entity to its own queue
func (q *Queue) Enqueue(rt *RawTweet) {
	if rt.RetweetedStatus != nil {
		// enqueue retweeted tweet
		q.Enqueue(rt.RetweetedStatus)
	}

	q.rts = append(q.rts, rt)
}

func (q *Queue) process(rts []*RawTweet) {
	defer q.Wg.Done()

	// ACCOUNTS
	var accs []model.Account
	for _, rt := range rts {
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
	for _, rt := range rts {
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
	for _, rt := range rts {
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
	// lock also countries, since we are selecting country_id
	q.sync.CountriesMutex.Lock()

	var queries []string
	var vars []interface{}

	for _, rt := range rts {
		var parentId *string = nil
		if rt.RetweetedStatus != nil {
			parentId = &rt.RetweetedStatus.IDStr
		}

		var countryId = "NULL"
		if rt.Place != nil && len(rt.Place.CountryCode) > 0 {
			countryId = fmt.Sprintf("(SELECT id FROM countries WHERE name='%s' LIMIT 1)", rt.Place.Country)
		}

		var location = "NULL"
		if rt.Geo != nil {
			location = fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), 4326)", rt.Geo.Coordinates[0], rt.Geo.Coordinates[1])
		}

		happenedAt := toTime(rt).Unix()

		// SQL Injection, NOT COOL!
		q := fmt.Sprintf("(?, ?, %s, ?, ?, to_timestamp(?), ?, %s, ?)", location, countryId)

		// This is quite hacky, but this is the best way to use sub queries in insert query
		queries = append(queries, q)
		vars = append(vars, []interface{}{rt.IDStr, rt.FullText, rt.RetweetCount, rt.FavoriteCount, happenedAt, rt.User.IDStr, parentId}...)
	}

	query := fmt.Sprintf("INSERT INTO tweets (id, content, location, retweet_count, favorite_count, happened_at, author_id, country_id, parent_id) VALUES %s ON CONFLICT DO NOTHING;", strings.Join(queries, ", "))
	q.db.Exec(query, vars...)

	q.sync.CountriesMutex.Unlock()
	q.sync.TweetsMutex.Unlock()

	// TWEET HASHTAGS
	q.sync.TweetHashtagsMutex.Lock()

	queries = []string{}
	vars = []interface{}{}

	for _, rt := range rts {
		for _, h := range rt.Entities.Hashtags {
			queries = append(queries, "((SELECT id FROM hashtags WHERE value=? LIMIT 1), ?)")
			vars = append(vars, []interface{}{h.Text, rt.IDStr}...)
		}
	}

	query = fmt.Sprintf("INSERT INTO tweet_hashtags (hashtag_id, tweet_id) VALUES %s ON CONFLICT DO NOTHING;", strings.Join(queries, ", "))
	// insert to tweet hashtags table
	q.db.Exec(query, vars...)

	q.sync.TweetHashtagsMutex.Unlock()

	// TWEET MENTIONS

	var tms []model.TweetMention
	for _, rt := range rts {
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
func (q *Queue) clear() []*RawTweet {
	rts := q.rts
	q.rts = []*RawTweet{}

	return rts
}

func toTime(rt *RawTweet) time.Time {
	t, err := time.Parse(time.RubyDate, rt.CreatedAt)
	if err != nil {
		log.Panicf("error while parsing time: %s", err)
	}
	return t
}
