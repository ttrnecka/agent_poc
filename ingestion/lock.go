package main

import "sync"

type LockManager struct {
	mutex sync.Mutex
	locks map[string]*Lock
}

type Lock struct {
	counter int
	mutex   sync.Mutex
}

func NewLockManager() *LockManager {
	return &LockManager{
		mutex: sync.Mutex{},
		locks: make(map[string]*Lock),
	}
}

func (lm *LockManager) Lock(key string) {
	lm.mutex.Lock()
	lock, exists := lm.locks[key]
	if !exists {
		lock = &Lock{
			counter: 0,
			mutex:   sync.Mutex{},
		}
		lm.locks[key] = lock
	}
	lock.counter++
	lm.mutex.Unlock()
	lock.mutex.Lock()
}

func (lm *LockManager) Unlock(key string) {
	lm.mutex.Lock()
	lock, exists := lm.locks[key]
	if exists {
		lock.counter--
		if lock.counter == 0 {
			delete(lm.locks, key) // remove from map safely
		}
	}
	lm.mutex.Unlock()
	if exists {
		lock.mutex.Unlock() // ok to unlock after map deletion
	}
}
