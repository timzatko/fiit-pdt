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
		// NOTE: insert order is IMPORTANT (since there are relations between tables, some entities must be inserted first!)
		// inset to accounts table
		err := q.insert(a)
		if err != nil {
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

func (q *Queue) insert(entity interface{}) error {
	q.sync.AccountsMutex.Lock()
	_, err := q.db.Model(&entity).OnConflict("DO NOTHING").Insert()
	q.sync.AccountsMutex.Unlock()

	return err
}