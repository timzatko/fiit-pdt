package model

import "time"

type Tweet struct {
	Id            string
	Content       string
	Location      interface{}
	RetweetCount  int
	FavoriteCount int
	HappenedAt    time.Time
	AuthorId      int64
	Author        *Account `pg:"rel:has-one"`
	CountryId     int
	ParentId      string
}
