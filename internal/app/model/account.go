package model

type Account struct {
	Id             int64
	ScreenName     string
	Name           string
	Description    string
	FollowersCount int
	FriendsCount   int
	StatusesCount  int
}
