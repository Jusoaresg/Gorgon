package concurrency

import (
	"sync"
)

func WithWriteLock(mutex *sync.Mutex, fn func() error) error {
	mutex.Lock()
	defer mutex.Unlock()
	return fn()
}
