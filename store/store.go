package store

import "sync"

// map could be sharded to minimize time spent waiting for locks
// or not? it could be sharded if we were to use simple map with hand-made locks
var store = sync.Map{}
