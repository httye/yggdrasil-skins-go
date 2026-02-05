// Package memory 内存Session缓存实现
package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"
)

// SessionCache 内存Session缓存
type SessionCache struct {
	sessions map[string]*sessionEntry // serverID -> sessionEntry
	mu       sync.RWMutex
}

// sessionEntry Session缓存条目
type sessionEntry struct {
	Session   *yggdrasil.Session
	ExpiresAt time.Time
}

// NewSessionCache 创建内存Session缓存
func NewSessionCache(options map[string]any) (*SessionCache, error) {
	return &SessionCache{
		sessions: make(map[string]*sessionEntry),
	}, nil
}

// Store 存储Session
func (c *SessionCache) Store(serverID string, session *yggdrasil.Session) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Session固定过期时间�?20�?	expiresAt := time.Now().Add(120 * time.Second)

	c.sessions[serverID] = &sessionEntry{
		Session:   session,
		ExpiresAt: expiresAt,
	}

	return nil
}

// Get 获取Session
func (c *SessionCache) Get(serverID string) (*yggdrasil.Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.sessions[serverID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// 检查是否过�?	if time.Now().After(entry.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	// 返回副本
	sessionCopy := *entry.Session
	return &sessionCopy, nil
}

// Delete 删除Session
func (c *SessionCache) Delete(serverID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.sessions, serverID)
	return nil
}

// CleanupExpired 清理过期Session
func (c *SessionCache) CleanupExpired() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for serverID, entry := range c.sessions {
		if now.After(entry.ExpiresAt) {
			delete(c.sessions, serverID)
		}
	}

	return nil
}

// Close 关闭缓存连接
func (c *SessionCache) Close() error {
	// 内存缓存无需关闭操作
	return nil
}

// GetCacheType 获取缓存类型
func (c *SessionCache) GetCacheType() string {
	return "memory"
}
