//in-memory map

package main

import "sync"

var (
	users = make(map[string]User)
	mu    sync.Mutex
)
