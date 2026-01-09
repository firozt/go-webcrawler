/*
	Thread safe in memory cache
	Implemented using a hashmap to a cacheObj
	CacheObj holds the data and metadata of the object, including optionals
	such as expiration date
*/

package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// ==================== structs ==================== //

type CacheObject[T any] struct {
	ExpireAt *int64
	Value    T
}

type CacheStore[T any] struct {
	data  map[string]*CacheObject[T]
	size  uint64
	mutex sync.Mutex
}

// ==================== methods ==================== //

func MakeCacheStore[T any]() *CacheStore[T] {
	return &CacheStore[T]{
		data: make(map[string]*CacheObject[T]),
		size: 0,
	}
}

func (store *CacheStore[T]) Get(key string) (T, bool) {
	var zero T
	store.mutex.Lock()
	defer store.mutex.Unlock()

	cacheObj, ok := store.data[key]

	if !ok {
		return zero, false // cache miss
	}

	now := time.Now().Unix()
	// has expiration and is expired
	if cacheObj.ExpireAt != nil && (*cacheObj.ExpireAt-now) <= 0 {
		delete(store.data, key) // delete from cache
		return zero, false
	}

	return cacheObj.Value, true
}

// set function with no ttl
func (store *CacheStore[T]) Set(key string, value T) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, seen := store.data[key]; !seen {
		store.incrementSize()
	}
	store.data[key] = &CacheObject[T]{
		Value:    value,
		ExpireAt: nil, // nil has no expiration
	}
}

// set function with a ttl
func (store *CacheStore[T]) SetExpires(key string, value T, duration int64) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	expiration := getUnixExpiration(duration)
	store.data[key] = &CacheObject[T]{
		Value:    value,
		ExpireAt: &expiration,
	}
}

// given duration (seconds) works out what unix time it expires at
func getUnixExpiration(duration int64) int64 {
	now := time.Now().Unix()
	expected := now + duration
	return expected
}

// size getter and incrementers
func (store *CacheStore[T]) Size() uint64 {
	return store.size
}

func (store *CacheStore[T]) incrementSize() {
	atomic.AddUint64(&store.size, 1)
}
