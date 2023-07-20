package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Config struct {
	Interval time.Duration
}

type PokeCache struct {
	cache  map[string]CacheEntry
	config Config
	mu     *sync.RWMutex
}

func NewCache(config Config) *PokeCache {
	pokeCache := PokeCache{
		cache:  map[string]CacheEntry{},
		config: config,
		mu:     &sync.RWMutex{},
	}

	go func() {
		for range time.Tick(config.Interval) {
			pokeCache.ReapLoop()
		}
	}()

	return &pokeCache
}

func (pokeCache *PokeCache) Add(key string, val []byte) {
	pokeCache.mu.Lock()
	defer pokeCache.mu.Unlock()

	pokeCache.cache[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (pokeCache *PokeCache) Get(key string) ([]byte, bool) {
	pokeCache.mu.RLock()
	defer pokeCache.mu.RUnlock()

	if val, ok := pokeCache.cache[key]; !ok {
		return []byte{}, false
	} else {
		return val.val, true
	}
}

func (pokeCache *PokeCache) ReapLoop() {
	pokeCache.mu.Lock()
	defer pokeCache.mu.Unlock()

	now := time.Now()

	for key, val := range pokeCache.cache {
		if now.Sub(val.createdAt) > pokeCache.config.Interval {
			delete(pokeCache.cache, key)
		}
	}
}
