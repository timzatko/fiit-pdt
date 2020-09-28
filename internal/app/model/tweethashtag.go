package model

type TweetHashtag struct {
	Id        int64 `gorm:"primarykey"`
	HashtagId int64
	TweetId   string
}
