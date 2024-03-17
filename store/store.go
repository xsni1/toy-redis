package store

import "sync"

type Config struct {
	Appendonly bool
}

// map could be sharded to minimize time spent waiting for locks
// or not? it could be sharded if we were to use simple map with hand-made locks
type Store struct {
	M sync.Map
}

func NewStore(config Config) *Store {
	return &Store{}
}
