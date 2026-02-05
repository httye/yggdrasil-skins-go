// Package cache 用户信息缓存
package cache

import (
	"sync"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/middleware"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"
)

// UserCacheItem 用户缓存�?type UserCacheItem struct {
	User      *yggdrasil.User
	ExpiresAt time.Time
}

// IsExpired 检查是否过�?func (item *UserCacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// UserCache 用户信息缓存
type UserCache struct {
	cache    sync.Map
	duration time.Duration
}

// NewUserCache 创建用户缓存
func NewUserCache(duration time.Duration) *UserCache {
	cache := &UserCache{
		duration: duration,
	}

	// 启动清理协程
	go cache.cleanup()

	return cache
}

// Get 获取用户信息
func (uc *UserCache) Get(key string) (*yggdrasil.User, bool) {
	if value, ok := uc.cache.Load(key); ok {
		item := value.(*UserCacheItem)
		if !item.IsExpired() {
			middleware.GlobalCacheMonitor.RecordHit()
			return item.User, true
		}
		// 过期则删�?		uc.cache.Delete(key)
	}

	middleware.GlobalCacheMonitor.RecordMiss()
	return nil, false
}

// Set 设置用户信息
func (uc *UserCache) Set(key string, user *yggdrasil.User) {
	item := &UserCacheItem{
		User:      user,
		ExpiresAt: time.Now().Add(uc.duration),
	}
	uc.cache.Store(key, item)
}

// Delete 删除用户信息
func (uc *UserCache) Delete(key string) {
	uc.cache.Delete(key)
}

// cleanup 定期清理过期缓存
func (uc *UserCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		uc.cache.Range(func(key, value interface{}) bool {
			item := value.(*UserCacheItem)
			if item.IsExpired() {
				uc.cache.Delete(key)
			}
			return true
		})
	}
}

// GetStats 获取缓存统计
func (uc *UserCache) GetStats() map[string]interface{} {
	count := 0
	expired := 0

	uc.cache.Range(func(key, value interface{}) bool {
		count++
		item := value.(*UserCacheItem)
		if item.IsExpired() {
			expired++
		}
		return true
	})

	return map[string]interface{}{
		"total_items":            count,
		"expired_items":          expired,
		"valid_items":            count - expired,
		"cache_duration_minutes": uc.duration.Minutes(),
	}
}

// 全局用户缓存实例（默�?分钟缓存，可通过配置修改�?var GlobalUserCache = NewUserCache(5 * time.Minute)

// InitUserCache 根据配置初始化用户缓�?func InitUserCache(duration time.Duration) {
	if duration > 0 {
		GlobalUserCache = NewUserCache(duration)
	}
}

// CachedUserLookup 带缓存的用户查询装饰�?func CachedUserLookup(key string, lookupFunc func() (*yggdrasil.User, error)) (*yggdrasil.User, error) {
	// 尝试从缓存获�?	if user, found := GlobalUserCache.Get(key); found {
		return user, nil
	}

	// 缓存未命中，执行查询
	user, err := lookupFunc()
	if err != nil {
		return nil, err
	}

	// 存储到缓�?	GlobalUserCache.Set(key, user)
	return user, nil
}
