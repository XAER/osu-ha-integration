package osu

import (
	"sync"
	"time"

	"github.com/XAER/osu-ha-integration/internal/domain"
)

type CacheItem struct {
	User      *domain.OsuUser
	Timestamp time.Time
}

type Cache struct {
	data     map[string]CacheItem
	duration time.Duration
	mu       sync.RWMutex
	logger   Logger
}

func NewCache(duration time.Duration, logger Logger) *Cache {
	return &Cache{
		data:     make(map[string]CacheItem),
		duration: duration,
		logger:   logger,
	}
}

func (c *Cache) Get(username string) (*domain.OsuUser, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[username]
	if !found {
		c.logger.Info("miss: " + username)
		return nil, false
	}
	if time.Since(item.Timestamp) > c.duration {
		c.logger.Info("expired: " + username)
		return nil, false
	}
	c.logger.Info("hit: " + username)
	return item.User, true
}

func (c *Cache) Set(username string, user *domain.OsuUser) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[username] = CacheItem{
		User:      user,
		Timestamp: time.Now(),
	}
}
