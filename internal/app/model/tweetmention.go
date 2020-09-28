package model

type TweetMention struct {
	Id        int64 `gorm:"primarykey"`
	AccountId int64
	TweetId   string
}
