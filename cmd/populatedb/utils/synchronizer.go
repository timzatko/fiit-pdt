package utils

import "sync"

// Synchronizer is used to sync inserting entities of different types
// between different files. So at one time, only to one entity is being written.
type Synchronizer struct {
	AccountsMutex      sync.Mutex
	HashtagsMutex      sync.Mutex
	CountriesMutex     sync.Mutex
	TweetsMutex        sync.Mutex
	TweetHashtagsMutex sync.Mutex
	TweetMentionsMutex sync.Mutex
}
