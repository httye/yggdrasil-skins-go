// Package file 文件Token缓存实现（BlessingSkin兼容�?package file

import (
	"fmt"
	"sync"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/utils"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"
)

// TokenCache 文件Token缓存（Laravel兼容�?type TokenCache struct {
	cache *LaravelFileCache
	mu    sync.RWMutex
}

// NewTokenCache 创建文件Token缓存
func NewTokenCache(options map[string]any) (*TokenCache, error) {
	cacheDir := "storage/framework/cache"
	if dir, ok := options["cache_dir"].(string); ok && dir != "" {
		cacheDir = dir
	}

	return &TokenCache{
		cache: NewLaravelFileCache(cacheDir),
	}, nil
}

// Store 存储Token（优化版：先验证JWT，提取信息）
func (c *TokenCache) Store(token *yggdrasil.Token) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 第一步：验证JWT并提取信�?	claims, err := utils.ValidateJWT(token.AccessToken)
	if err != nil {
		return fmt.Errorf("invalid JWT token: %w", err)
	}

	// 计算TTL
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("token already expired")
	}

	// 创建简化的Token对象（只存储JWT中没有的信息�?	cacheToken := &yggdrasil.Token{
		AccessToken: token.AccessToken, // 保留完整的AccessToken用于兼容�?		ClientToken: token.ClientToken,
		ProfileID:   token.ProfileID,
		Owner:       claims.UserID, // 从JWT中获取用户ID
		CreatedAt:   token.CreatedAt,
		ExpiresAt:   token.ExpiresAt,
	}

	// 存储Token（使用用户ID+TokenID作为键）
	tokenKey := generateOptimizedTokenKey(claims.UserID, claims.TokenID)
	if err := c.cache.Store(tokenKey, cacheToken, ttl); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	// 更新用户Token列表（使用用户ID�?	userTokensKey := generateYggdrasilUserTokensKey(claims.UserID)

	// 获取现有Token列表
	var existingTokens []string
	if err := c.cache.Get(userTokensKey, &existingTokens); err != nil {
		// 如果获取失败，创建新列表
		existingTokens = []string{}
	}

	// 检查Token是否已存�?	found := false
	for _, accessToken := range existingTokens {
		if accessToken == token.AccessToken {
			found = true
			break
		}
	}

	// 如果不存在，添加到列�?	if !found {
		existingTokens = append(existingTokens, token.AccessToken)
	}

	// 存储更新后的Token列表（使用较长的TTL�?	userTokensTTL := 7 * 24 * time.Hour // 7�?	if err := c.cache.Store(userTokensKey, existingTokens, userTokensTTL); err != nil {
		return fmt.Errorf("failed to store user tokens list: %w", err)
	}

	return nil
}

// Get 获取Token（优化版：先验证JWT，按需查询缓存�?func (c *TokenCache) Get(accessToken string) (*yggdrasil.Token, error) {
	// 第一步：验证JWT（本地计算，极快�?	claims, err := utils.ValidateJWT(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT token: %w", err)
	}

	// 第二步：从缓存获取ClientToken等额外信�?	c.mu.RLock()
	defer c.mu.RUnlock()

	tokenKey := generateOptimizedTokenKey(claims.UserID, claims.TokenID)

	var token yggdrasil.Token
	if err := c.cache.Get(tokenKey, &token); err != nil {
		return nil, fmt.Errorf("token not found in cache: %w", err)
	}

	// 构建Token对象（结合JWT信息和缓存信息）
	result := &yggdrasil.Token{
		AccessToken: accessToken,
		ClientToken: token.ClientToken,
		ProfileID:   claims.ProfileID,
		Owner:       claims.UserID,
		CreatedAt:   token.CreatedAt,
		ExpiresAt:   token.ExpiresAt,
	}
	return result, nil
}

// Delete 删除Token（优化版：先验证JWT，提取用户ID和TokenID�?func (c *TokenCache) Delete(accessToken string) error {
	// 先验证JWT并提取信�?	claims, err := utils.ValidateJWT(accessToken)
	if err != nil {
		// JWT无效，但仍然尝试删除（兼容性）
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 从用户Token列表中移除（使用用户ID�?	c.removeTokenFromUserList(claims.UserID, accessToken)

	// 删除Token
	tokenKey := generateOptimizedTokenKey(claims.UserID, claims.TokenID)
	return c.cache.Delete(tokenKey)
}

// GetUserTokens 获取用户的所有Token（按用户ID查询�?func (c *TokenCache) GetUserTokens(userID string) ([]*yggdrasil.Token, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userTokensKey := generateYggdrasilUserTokensKey(userID)

	var accessTokens []string
	if err := c.cache.Get(userTokensKey, &accessTokens); err != nil {
		return []*yggdrasil.Token{}, nil
	}

	var tokens []*yggdrasil.Token
	for _, accessToken := range accessTokens {
		// Laravel缓存已经处理了过期检查，如果能获取到Token就说明没有过�?		if token, err := c.Get(accessToken); err == nil {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

// DeleteUserTokens 删除用户的所有Token（按用户ID�?func (c *TokenCache) DeleteUserTokens(userID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 获取用户的所有Token
	userTokensKey := generateYggdrasilUserTokensKey(userID)

	var accessTokens []string
	if err := c.cache.Get(userTokensKey, &accessTokens); err != nil {
		return nil // 用户没有Token
	}

	// 删除所有Token
	for _, accessToken := range accessTokens {
		tokenKey := generateYggdrasilTokenKey(accessToken)
		c.cache.Delete(tokenKey)
	}

	// 删除用户Token列表
	return c.cache.Delete(userTokensKey)
}

// GetUserTokenCount 获取用户Token数量
func (c *TokenCache) GetUserTokenCount(userID string) (int, error) {
	tokens, err := c.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}
	return len(tokens), nil
}

// CleanupExpired 清理过期Token
func (c *TokenCache) CleanupExpired() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.CleanupExpired()
}

// Close 关闭缓存连接
func (c *TokenCache) Close() error {
	// 文件缓存无需关闭操作
	return nil
}

// GetCacheType 获取缓存类型
func (c *TokenCache) GetCacheType() string {
	return "file"
}

// removeTokenFromUserList 从用户Token列表中移除指定Token
func (c *TokenCache) removeTokenFromUserList(userID, accessToken string) error {
	userTokensKey := generateYggdrasilUserTokensKey(userID)

	var accessTokens []string
	if err := c.cache.Get(userTokensKey, &accessTokens); err != nil {
		return nil // 用户没有Token列表
	}

	// 移除指定Token
	for i, token := range accessTokens {
		if token == accessToken {
			accessTokens = append(accessTokens[:i], accessTokens[i+1:]...)
			break
		}
	}

	// 如果列表为空，删除用户Token列表
	if len(accessTokens) == 0 {
		return c.cache.Delete(userTokensKey)
	}

	// 更新Token列表
	userTokensTTL := 7 * 24 * time.Hour // 7�?	return c.cache.Store(userTokensKey, accessTokens, userTokensTTL)
}

// generateOptimizedTokenKey 生成优化的Token键（用户ID+TokenID�?func generateOptimizedTokenKey(userID, tokenID string) string {
	return fmt.Sprintf("yggdrasil:token:%s:%s", userID, tokenID)
}
