// Package cache 定义缓存接口和相关类�?package cache

import (
	"time"

	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"
)

// TokenCache Token缓存接口
type TokenCache interface {
	// Store 存储Token
	Store(token *yggdrasil.Token) error

	// Get 获取Token
	Get(accessToken string) (*yggdrasil.Token, error)

	// Delete 删除Token
	Delete(accessToken string) error

	// GetUserTokens 获取用户的所有Token
	GetUserTokens(userID string) ([]*yggdrasil.Token, error)

	// DeleteUserTokens 删除用户的所有Token
	DeleteUserTokens(userID string) error

	// GetUserTokenCount 获取用户Token数量
	GetUserTokenCount(userID string) (int, error)

	// CleanupExpired 清理过期Token
	CleanupExpired() error

	// Close 关闭缓存连接
	Close() error

	// GetCacheType 获取缓存类型
	GetCacheType() string
}

// SessionCache Session缓存接口
type SessionCache interface {
	// Store 存储Session
	Store(serverID string, session *yggdrasil.Session) error

	// Get 获取Session
	Get(serverID string) (*yggdrasil.Session, error)

	// Delete 删除Session
	Delete(serverID string) error

	// CleanupExpired 清理过期Session
	CleanupExpired() error

	// Close 关闭缓存连接
	Close() error

	// GetCacheType 获取缓存类型
	GetCacheType() string
}

// CacheFactory 缓存工厂接口
type CacheFactory interface {
	// CreateTokenCache 创建Token缓存实例
	CreateTokenCache(cacheType string, options map[string]any) (TokenCache, error)

	// CreateSessionCache 创建Session缓存实例
	CreateSessionCache(cacheType string, options map[string]any) (SessionCache, error)

	// GetSupportedTypes 获取支持的缓存类�?	GetSupportedTypes() []string
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type    string         `json:"type"`    // 缓存类型
	Options map[string]any `json:"options"` // 缓存选项
}

// BlessingSkinCacheEntry BlessingSkin缓存条目格式
type BlessingSkinCacheEntry struct {
	Data      any       `json:"data"`       // 缓存数据
	ExpiresAt time.Time `json:"expires_at"` // 过期时间
}

// LaravelCacheEntry Laravel缓存条目格式（用于文件缓存兼容）
type LaravelCacheEntry struct {
	Value     string `json:"value"`      // 序列化的数据
	ExpiresAt int64  `json:"expires_at"` // 过期时间�?}
