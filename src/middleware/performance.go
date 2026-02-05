// Package middleware 性能监控中间�?package middleware

import (
	"time"

	"github.com/httye/yggdrasil-skins-go/src/utils"

	"github.com/gin-gonic/gin"
)

// PerformanceMonitor 性能监控中间�?// 自动记录所有API请求的性能指标
func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录性能指标
		duration := time.Since(start)
		isError := c.Writer.Status() >= 400

		// 记录到全局性能监控
		utils.GlobalMetrics.RecordRequest(duration, isError)
	}
}

// DatabaseQueryMonitor 数据库查询监控装饰器
// 用于包装数据库查询方法，自动记录查询性能
func DatabaseQueryMonitor(queryFunc func() error) error {
	start := time.Now()
	err := queryFunc()
	duration := time.Since(start)

	// 记录数据库查询性能
	utils.GlobalMetrics.RecordDBQuery(duration)

	return err
}

// CacheMonitor 缓存监控工具
type CacheMonitor struct{}

// RecordHit 记录缓存命中
func (cm *CacheMonitor) RecordHit() {
	utils.GlobalMetrics.RecordCacheHit()
}

// RecordMiss 记录缓存未命�?func (cm *CacheMonitor) RecordMiss() {
	utils.GlobalMetrics.RecordCacheMiss()
}

// 全局缓存监控实例
var GlobalCacheMonitor = &CacheMonitor{}
