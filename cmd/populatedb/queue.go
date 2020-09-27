package main

import (
	"log"
	"sync"

	"github.com/go-pg/pg"

	"github.com/timzatko/fiit-pdt/internal/app/model"
)

type Queue struct {
	Accounts []model.Account
	wg       sync.WaitGroup
	db       *pg.DB
	sync     *Synchronizer
}

func NewQueue(db *pg.DB, synchronizer *Synchronizer) Queue {
	return Queue{
		Accounts: []model.Account{},
		db:       db,
		sync:     synchronizer,
	}
}

func (q *Queue) isFull() bool {
	return len(q.Accounts) >= 1000
}

func (q *Queue) isEmpty() bool {
	return len(q.Accounts) == 0
}

func (q *Queue) send() {
	if !q.isFull() {
		return
	}

	q.wg.Add(1)

	a := q.Accounts
	q.Accounts = []model.Account{}

	go func() {
		// TODO: common mutex for all files running in parallel
		q.sync.AccountsMutex.Lock()
		_, err := q.db.Model(&a).OnConflict("DO NOTHING").Insert()
		q.sync.AccountsMutex.Unlock()

		if err != nil {
			q.sync.AccountsMutex.Unlock()
			log.Panicf("error while inserting: %s", err)
		}

		// TODO: insert to hashtags
		// TODO: insert to countries
		// TODO: insert to tweets
		// TODO: insert to tweet_hashtags
		// TODO: insert to tweet_mentions

		defer q.wg.Done()
	}()
}

func (q *Queue) add(rt RawTweet) {
	acc := model.Account{
		Id:             rt.User.ID,
		ScreenName:     rt.User.ScreenName,
		Name:           rt.User.Name,
		Description:    rt.User.Description,
		FollowersCount: rt.User.FollowersCount,
		FriendsCount:   rt.User.FriendsCount,
		StatusesCount:  rt.User.StatusesCount,
	}

	q.Accounts = append(q.Accounts, acc)
}
