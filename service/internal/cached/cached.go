package cached

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// CacheManager представляет собой менеджер кэша.
type CacheManager struct {
	Cache *cache.Cache
	Mu    sync.Mutex
}

// NewCache создает новый CacheManager.
func NewCache(defaultExpiration, cleanupInterval time.Duration) *CacheManager {
	return &CacheManager{
		Cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set устанавливает значение в кэше с заданным ключом и сроком действия.
func (cm *CacheManager) Set(key string, value interface{}, expiration time.Duration) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Set(key, value, expiration)
}

// Get получает значение из кэша по ключу.
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	return cm.Cache.Get(key)
}

// Delete удаляет значение из кэша по ключу.
func (cm *CacheManager) Delete(key string) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Delete(key)
}

// Flush очищает весь кэш.
func (cm *CacheManager) Flush() {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Cache.Flush()
}
