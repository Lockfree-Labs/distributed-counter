package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// DefaultExpiration is the default duration for which a key should live.
const DefaultExpiration = 24 * time.Hour

// counterEntry holds the value and the expiration time for a counter.
type counterEntry struct {
	value  int64
	expiry time.Time
}

// CounterManager manages counters stored in memory and synced with Redis.
type CounterManager struct {
	mu          sync.RWMutex
	keyPrefix   string
	counters    map[string]*counterEntry
	changed     map[string]bool
	redisClient *redis.Client
	ctx         context.Context
}

// NewCounterManager creates a new instance of CounterManager, initializes counters
// from Redis, and starts the periodic dump routine.
func NewCounterManager() *CounterManager {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword,
	})

	cm := &CounterManager{
		keyPrefix:   "counter_",
		counters:    make(map[string]*counterEntry),
		changed:     make(map[string]bool),
		redisClient: rdb,
		ctx:         ctx,
	}

	// Initialize counters from Redis. For simplicity, we load all keys.
	keys, err := rdb.Keys(ctx, "counter_*").Result()
	if err != nil {
		fmt.Println("Error fetching keys:", err)
	} else {
		now := time.Now()
		for _, key := range keys {
			val, err := rdb.Get(ctx, key).Int64()
			if err != nil {
				fmt.Println("Error getting key", key, ":", err)
				continue
			}
			ttl, err := rdb.TTL(ctx, key).Result()
			if err != nil {
				fmt.Println("Error getting TTL for key", key, ":", err)
				continue
			}
			// If ttl is negative (key without expiry or expired), use default expiration.
			if ttl <= 0 {
				ttl = DefaultExpiration
			}
			cm.counters[key] = &counterEntry{
				value:  val,
				expiry: now.Add(ttl),
			}
		}
	}

	// Start the periodic dump to Redis every 5 seconds.
	go cm.periodicDump()

	return cm
}

// isExpired checks if the given counterEntry has expired.
func (cm *CounterManager) isExpired(entry *counterEntry) bool {
	return time.Now().After(entry.expiry)
}

// Increment increases the counter for the provided key by 1.
// If the key is not present or expired, it is (re)initialized with 0 and an expiry 24 hours from now.
func (cm *CounterManager) Increment(key string) int64 {
	key = fmt.Sprintf("%s%s", cm.keyPrefix, key)

	cm.mu.Lock()
	defer cm.mu.Unlock()
	now := time.Now()

	entry, exists := cm.counters[key]
	if !exists || cm.isExpired(entry) {
		entry = &counterEntry{
			value:  0,
			expiry: now.Add(DefaultExpiration),
		}
		cm.counters[key] = entry
	}
	entry.value++
	cm.changed[key] = true
	return entry.value
}

// Get returns the current value of the counter for the provided key.
// If the key has expired, it returns 0.
func (cm *CounterManager) Get(key string) int64 {
	key = fmt.Sprintf("%s%s", cm.keyPrefix, key)
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	entry, exists := cm.counters[key]
	if !exists {
		return 0
	}
	return entry.value
}

// periodicDump triggers dumping of changed keys to Redis every 5 seconds.
func (cm *CounterManager) periodicDump() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		cm.dumpChangedKeys()
	}
}

// dumpChangedKeys writes the changed keys to Redis with the correct expiration
// (the time left until the local expiry) and resets the change tracker.
// If a key has expired, it is removed from Redis.
func (cm *CounterManager) dumpChangedKeys() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for key := range cm.changed {
		entry, exists := cm.counters[key]
		if !exists || cm.isExpired(entry) {
			// Delete key from Redis if it has expired.
			if err := cm.redisClient.Del(cm.ctx, key).Err(); err != nil {
				fmt.Println("Error deleting expired key", key, "from Redis:", err)
			}
			delete(cm.counters, key)
		} else {
			// Reset expiry
			if err := cm.redisClient.Set(cm.ctx, key, entry.value, DefaultExpiration).Err(); err != nil {
				fmt.Println("Error dumping key", key, ":", err)
			}
		}
	}
	// Reset the changed map after dumping.
	cm.changed = make(map[string]bool)
}
