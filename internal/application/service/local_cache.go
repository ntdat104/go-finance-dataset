package service

import (
	"sync"
	"time"
)

type LocalCacheSvc interface {
	Set(key string, value any, ttl time.Duration)
	Get(key string) (any, bool)
	Del(key string)
	Has(key string) bool
}

type cacheItem struct {
	value      any
	expireTime time.Time
}

type localCacheSvc struct {
	store sync.Map
}

func NewLocalCacheSvc() LocalCacheSvc {
	s := &localCacheSvc{
		store: sync.Map{},
	}
	// Start cleanup ticker
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.cleanUp()
		}
	}()
	return s
}

func (c *localCacheSvc) Set(key string, value any, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	c.store.Store(key, cacheItem{
		value:      value,
		expireTime: expire,
	})
}

func (c *localCacheSvc) GetExpireTime(key string) (*time.Time, bool) {
	val, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(cacheItem)
	return &item.expireTime, true
}

func (c *localCacheSvc) Get(key string) (any, bool) {
	val, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(cacheItem)
	if time.Now().After(item.expireTime) {
		c.store.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *localCacheSvc) Del(key string) {
	c.store.Delete(key)
}

func (c *localCacheSvc) Has(key string) bool {
	_, exists := c.Get(key)
	return exists
}

func (c *localCacheSvc) cleanUp() {
	now := time.Now()
	c.store.Range(func(key, val any) bool {
		item := val.(cacheItem)
		if now.After(item.expireTime) {
			c.store.Delete(key)
		}
		return true
	})
}
