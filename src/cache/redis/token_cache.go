// Package redis Redis Token缓存实现（BlessingSkin兼容�?package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/utils"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

// TokenCache Redis Token缓存
type TokenCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewTokenCache 创建Redis Token缓存
func NewTokenCache(options map[string]any) (*TokenCache, error) {
	redisURL := "redis://localhost:6379"
	if url, ok := options["redis_url"].(string); ok && url != "" {
		redisURL = url
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &TokenCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Store 存储Token（优化版：先验证JWT，提取信息）
func (c *TokenCache) Store(token *yggdrasil.Token) error {
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

	// 序列化Token
	tokenData, err := sonic.Marshal(cacheToken)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// 计算TTL
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("token already expired")
	}

	// 存储Token（使用用户ID:TokenID作为键）
	tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", claims.UserID, claims.TokenID)
	if err := c.client.Set(c.ctx, tokenKey, tokenData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	// 更新用户Token列表（使用用户ID�?	userTokensKey := fmt.Sprintf("yggdrasil-id-%s", claims.UserID)
	if err := c.client.SAdd(c.ctx, userTokensKey, claims.TokenID).Err(); err != nil {
		return fmt.Errorf("failed to add token to user list: %w", err)
	}

	// 设置用户Token列表的过期时间（7天）
	c.client.Expire(c.ctx, userTokensKey, 7*24*time.Hour)

	return nil
}

// Get 获取Token（优化版：先验证JWT，按需查询缓存�?func (c *TokenCache) Get(accessToken string) (*yggdrasil.Token, error) {
	// 第一步：验证JWT（本地计算，极快�?	claims, err := utils.ValidateJWT(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT token: %w", err)
	}

	// 第二步：从缓存获取ClientToken等额外信�?	tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", claims.UserID, claims.TokenID)

	data, err := c.client.Get(c.ctx, tokenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("token not found in cache")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var token yggdrasil.Token
	if err := sonic.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
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

	// 从用户Token列表中移除（使用用户ID�?	userTokensKey := fmt.Sprintf("yggdrasil-id-%s", claims.UserID)
	c.client.SRem(c.ctx, userTokensKey, claims.TokenID)

	// 删除Token
	tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", claims.UserID, claims.TokenID)
	return c.client.Del(c.ctx, tokenKey).Err()
}

// GetUserTokens 获取用户的所有Token（按用户ID查询�?func (c *TokenCache) GetUserTokens(userID string) ([]*yggdrasil.Token, error) {
	userTokensKey := fmt.Sprintf("yggdrasil-id-%s", userID)

	tokenIDs, err := c.client.SMembers(c.ctx, userTokensKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []*yggdrasil.Token{}, nil
		}
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	var tokens []*yggdrasil.Token
	for _, tokenID := range tokenIDs {
		// 直接从Redis获取Token数据
		tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", userID, tokenID)
		data, err := c.client.Get(c.ctx, tokenKey).Result()
		if err != nil {
			// 清理无效的Token引用
			c.client.SRem(c.ctx, userTokensKey, tokenID)
			continue
		}

		var token yggdrasil.Token
		if err := sonic.Unmarshal([]byte(data), &token); err != nil {
			// 清理无效的Token引用
			c.client.SRem(c.ctx, userTokensKey, tokenID)
			continue
		}

		tokens = append(tokens, &token)
	}

	return tokens, nil
}

// DeleteUserTokens 删除用户的所有Token（按用户ID�?func (c *TokenCache) DeleteUserTokens(userID string) error {
	userTokensKey := fmt.Sprintf("yggdrasil-id-%s", userID)

	// 获取用户的所有TokenID
	tokenIDs, err := c.client.SMembers(c.ctx, userTokensKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // 用户没有Token
		}
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	// 删除所有Token
	for _, tokenID := range tokenIDs {
		tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", userID, tokenID)
		c.client.Del(c.ctx, tokenKey)
	}

	// 删除用户Token列表
	return c.client.Del(c.ctx, userTokensKey).Err()
}

// GetUserTokenCount 获取用户Token数量
func (c *TokenCache) GetUserTokenCount(userID string) (int, error) {
	userTokensKey := fmt.Sprintf("yggdrasil-id-%s", userID)

	count, err := c.client.SCard(c.ctx, userTokensKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get user token count: %w", err)
	}

	return int(count), nil
}

// CleanupExpired 清理过期Token
func (c *TokenCache) CleanupExpired() error {
	// Redis会自动清理过期的键，这里主要清理用户Token列表中的无效引用

	// 获取所有用户Token列表�?	keys, err := c.client.Keys(c.ctx, "yggdrasil-id-*").Result()
	if err != nil {
		return fmt.Errorf("failed to get user token keys: %w", err)
	}

	for _, userTokensKey := range keys {
		// 提取用户ID
		userID := userTokensKey[len("yggdrasil-id-"):]

		// 获取用户TokenID列表
		tokenIDs, err := c.client.SMembers(c.ctx, userTokensKey).Result()
		if err != nil {
			continue
		}

		// 检查每个Token是否仍然存在
		for _, tokenID := range tokenIDs {
			tokenKey := fmt.Sprintf("yggdrasil-token-%s:%s", userID, tokenID)
			exists, err := c.client.Exists(c.ctx, tokenKey).Result()
			if err != nil || exists == 0 {
				// Token不存在，从用户列表中移除
				c.client.SRem(c.ctx, userTokensKey, tokenID)
			}
		}

		// 如果用户Token列表为空，删除该列表
		count, _ := c.client.SCard(c.ctx, userTokensKey).Result()
		if count == 0 {
			c.client.Del(c.ctx, userTokensKey)
		}
	}

	return nil
}

// Close 关闭缓存连接
func (c *TokenCache) Close() error {
	return c.client.Close()
}

// GetCacheType 获取缓存类型
func (c *TokenCache) GetCacheType() string {
	return "redis"
}
