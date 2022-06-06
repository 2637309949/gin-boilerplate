package store

import (
	"strconv"
	"time"

	"github.com/robfig/go-cache"
)

//InMemoryStore represents the cache with memory persistence
type InMemoryStore struct {
	cache.Cache
}

// NewInMemoryStore returns a InMemoryStore
func NewInMemoryStore(defaultExpiration time.Duration) *InMemoryStore {
	return &InMemoryStore{*cache.New(defaultExpiration, time.Minute)}
}

// Get (see CacheStore interface)
func (c *InMemoryStore) Get(key string, ptrValue interface{}) error {
	val, found := c.Cache.Get(key)
	if !found {
		return ErrCacheMiss
	}
	valBytes, err := Serialize(val)
	if err != nil {
		return err
	}
	return Deserialize(valBytes, ptrValue)
}

// Set (see CacheStore interface)
func (c *InMemoryStore) Set(key string, value interface{}, expires time.Duration) error {
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	c.Cache.Set(key, b, expires)
	return nil
}

// Add (see CacheStore interface)
func (c *InMemoryStore) Add(key string, value interface{}, expires time.Duration) error {
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	err = c.Cache.Add(key, b, expires)
	if err == cache.ErrKeyExists {
		return ErrNotStored
	}
	return err
}

// Replace (see CacheStore interface)
func (c *InMemoryStore) Replace(key string, value interface{}, expires time.Duration) error {
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	if err := c.Cache.Replace(key, b, expires); err != nil {
		return ErrNotStored
	}
	return nil
}

// Delete (see CacheStore interface)
func (c *InMemoryStore) Delete(key string) error {
	if found := c.Cache.Delete(key); !found {
		return ErrCacheMiss
	}
	return nil
}

// Increment (see CacheStore interface)
func (c *InMemoryStore) Increment(key string, delta uint64) (uint64, error) {
	val, found := c.Cache.Get(key)
	if !found {
		return 0, ErrCacheMiss
	}
	valBytes, err := Serialize(val)
	if err != nil {
		return 0, err
	}
	currentVal, err := strconv.ParseUint(string(valBytes), 10, 64)
	if err == nil {
		sum := currentVal + delta
		b, err := Serialize(sum)
		if err != nil {
			return 0, err
		}
		c.Cache.Set(key, b, 0)
		return sum, nil
	}
	newValue, err := c.Cache.Increment(key, delta)
	if err == cache.ErrCacheMiss {
		return 0, ErrCacheMiss
	}
	return newValue, err
}

// Decrement (see CacheStore interface)
func (c *InMemoryStore) Decrement(key string, delta uint64) (uint64, error) {
	val, found := c.Cache.Get(key)
	if !found {
		return 0, ErrCacheMiss
	}
	valBytes, err := Serialize(val)
	if err != nil {
		return 0, err
	}

	currentVal, err := strconv.ParseUint(string(valBytes), 10, 64)
	if err == nil {
		sum := currentVal - delta
		if delta > currentVal {
			sum = 0
		}
		b, err := Serialize(sum)
		if err != nil {
			return 0, err
		}
		c.Cache.Set(key, b, 0)
		return sum, nil
	}
	newValue, err := c.Cache.Decrement(key, delta)
	if err == cache.ErrCacheMiss {
		return 0, ErrCacheMiss
	}
	return newValue, err
}

// Flush (see CacheStore interface)
func (c *InMemoryStore) Flush() error {
	c.Cache.Flush()
	return nil
}
