package main

import "sync"

type LockManager struct {
	mutex sync.Mutex
	locks map[string]*sync.Mutex
}

func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*sync.Mutex),
	}
}

func (lm *LockManager) Lock(key string) {
	lm.mutex.Lock()
	lock, exists := lm.locks[key]
	if !exists {
		lock = &sync.Mutex{}
		lm.locks[key] = lock
	}
	lm.mutex.Unlock()

	lock.Lock()
}

func (lm *LockManager) Unlock(key string) {
	lm.mutex.Lock()
	lock, exists := lm.locks[key]
	if exists {
		delete(lm.locks, key) // remove from map safely
	}
	lm.mutex.Unlock()

	if exists {
		lock.Unlock() // ok to unlock after map deletion
	}
}
