package id

import "sync"

var current = uint64(0)
var mu = sync.Mutex{}

func Next() uint64 {
	mu.Lock()
	defer mu.Unlock()

	v := current

	current++

	return v
}
