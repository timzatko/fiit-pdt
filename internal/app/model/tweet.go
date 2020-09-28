package model

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Tweet struct {
	Id            string `gorm:"primarykey"`
	Content       string
	Location      *Location
	RetweetCount  int64
	FavoriteCount int64
	HappenedAt    Timestamp
	AuthorId      string
	CountryId     *int64
	ParentId      *string
	Country       Country
	Author        Account
	Hashtags      []Hashtag `gorm:"many2many:tweet_hashtags;"`
	Mentions      []Account `gorm:"many2many:tweet_mentions;"`
}

// Location
type Location struct {
	X, Y float64
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
	// TODO: Scan a value into struct from database driver
	log.Println(v)

	return nil
}

func (loc *Location) GormDataType() string {
	return "geometry"
}

func (loc *Location) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	if loc != nil {
		return clause.Expr{
			SQL:  "ST_SetSRID(ST_MakePoint(?, ?), 4326)",
			Vars: []interface{}{loc.X, loc.Y},
		}
	}

	// if loc is nil than the location is provided, hence in the postgres it should be NULL
	return clause.Expr{SQL: "NULL"}
}

// Timestamp
type Timestamp struct {
	T time.Time
}

// Scan implements the sql.Scanner interface
func (t *Timestamp) Scan(v interface{}) error {
	// TODO: Scan a value into struct from database driver
	log.Println(v)

	return nil
}

func (t Timestamp) GormDataType() string {
	return "geometry"
}

func (t Timestamp) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "to_timestamp(?)",
		Vars: []interface{}{t.T.Unix()},
	}
}
