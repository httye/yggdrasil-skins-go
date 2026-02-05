// Package memory 内存缓存实现
package memory

import (
	"fmt"
	"sync"

	"github.com/httye/yggdrasil-skins-go/src/utils"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"
)

// TokenCache 内存Token缓存（优化版：支持JWT优先验证�?type TokenCache struct {
	tokens     map[string]*yggdrasil.Token // "userID:tokenID" -> Token（简化版�?	userTokens map[string][]string         // userID -> []tokenID
	mu         sync.RWMutex
}

// NewTokenCache 创建内存Token缓存
func NewTokenCache(options map[string]any) (*TokenCache, error) {
	return &TokenCache{
		tokens:     make(map[string]*yggdrasil.Token),
		userTokens: make(map[string][]string),
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

	// 创建简化的Token对象（只存储JWT中没有的信息�?	cacheToken := &yggdrasil.Token{
		AccessToken: token.AccessToken, // 保留完整的AccessToken用于兼容�?		ClientToken: token.ClientToken,
		ProfileID:   claims.ProfileID, // 从JWT中获取ProfileID
		Owner:       claims.UserID, // 从JWT中获取用户ID
		CreatedAt:   token.CreatedAt,
		ExpiresAt:   token.ExpiresAt,
	}

	// 存储Token（使用用户ID:TokenID作为键）
	tokenKey := fmt.Sprintf("%s:%s", claims.UserID, claims.TokenID)
	c.tokens[tokenKey] = cacheToken

	// 更新用户Token列表（使用用户ID�?	userTokens := c.userTokens[claims.UserID]

	// 检查是否已存在
	found := false
	for _, tokenID := range userTokens {
		if tokenID == claims.TokenID {
			found = true
			break
		}
	}

	if !found {
		c.userTokens[claims.UserID] = append(userTokens, claims.TokenID)
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

	tokenKey := fmt.Sprintf("%s:%s", claims.UserID, claims.TokenID)
	token, exists := c.tokens[tokenKey]
	if !exists {
		return nil, fmt.Errorf("token not found in cache")
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

// Delete 删除Token
func (c *TokenCache) Delete(accessToken string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 获取Token以找到所有�?	token, exists := c.tokens[accessToken]
	if exists {
		// 从用户Token列表中移�?		userTokens := c.userTokens[token.Owner]
		for i, userToken := range userTokens {
			if userToken == accessToken {
				c.userTokens[token.Owner] = append(userTokens[:i], userTokens[i+1:]...)
				break
			}
		}

		// 如果用户没有Token了，删除用户条目
		if len(c.userTokens[token.Owner]) == 0 {
			delete(c.userTokens, token.Owner)
		}
	}

	// 删除Token
	delete(c.tokens, accessToken)
	return nil
}

// GetUserTokens 获取用户的所有Token
func (c *TokenCache) GetUserTokens(userEmail string) ([]*yggdrasil.Token, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	accessTokens, exists := c.userTokens[userEmail]
	if !exists {
		return []*yggdrasil.Token{}, nil
	}

	var tokens []*yggdrasil.Token
	for _, accessToken := range accessTokens {
		if token, exists := c.tokens[accessToken]; exists && token.IsValid() {
			tokenCopy := *token
			tokens = append(tokens, &tokenCopy)
		}
	}

	return tokens, nil
}

// DeleteUserTokens 删除用户的所有Token
func (c *TokenCache) DeleteUserTokens(userEmail string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	accessTokens, exists := c.userTokens[userEmail]
	if !exists {
		return nil
	}

	// 删除所有Token
	for _, accessToken := range accessTokens {
		delete(c.tokens, accessToken)
	}

	// 删除用户Token列表
	delete(c.userTokens, userEmail)
	return nil
}

// GetUserTokenCount 获取用户Token数量
func (c *TokenCache) GetUserTokenCount(userEmail string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	accessTokens, exists := c.userTokens[userEmail]
	if !exists {
		return 0, nil
	}

	count := 0
	for _, accessToken := range accessTokens {
		if token, exists := c.tokens[accessToken]; exists && token.IsValid() {
			count++
		}
	}

	return count, nil
}

// CleanupExpired 清理过期Token
func (c *TokenCache) CleanupExpired() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 收集过期的Token
	var expiredTokens []string
	for accessToken, token := range c.tokens {
		if !token.IsValid() {
			expiredTokens = append(expiredTokens, accessToken)
		}
	}

	// 删除过期Token
	for _, accessToken := range expiredTokens {
		token := c.tokens[accessToken]

		// 从用户Token列表中移�?		userTokens := c.userTokens[token.Owner]
		for i, userToken := range userTokens {
			if userToken == accessToken {
				c.userTokens[token.Owner] = append(userTokens[:i], userTokens[i+1:]...)
				break
			}
		}

		// 如果用户没有Token了，删除用户条目
		if len(c.userTokens[token.Owner]) == 0 {
			delete(c.userTokens, token.Owner)
		}

		// 删除Token
		delete(c.tokens, accessToken)
	}

	return nil
}

// Close 关闭缓存连接
func (c *TokenCache) Close() error {
	// 内存缓存无需关闭操作
	return nil
}

// GetCacheType 获取缓存类型
func (c *TokenCache) GetCacheType() string {
	return "memory"
}
