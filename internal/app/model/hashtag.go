package model

type Hashtag struct {
	Id    int64  `gorm:"primarykey"`
	Value string `gorm:"index:,unique"`
}
