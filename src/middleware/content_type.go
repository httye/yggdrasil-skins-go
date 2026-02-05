// Package middleware 提供HTTP中间�?package middleware

import (
	"net/http"
	"strings"

	"github.com/httye/yggdrasil-skins-go/src/utils"

	"github.com/gin-gonic/gin"
)

// CheckContentType 检查请求内容类型的中间�?func CheckContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对POST、PUT、PATCH请求检查Content-Type
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {

			contentType := c.GetHeader("Content-Type")

			// 检查是否为JSON格式
			if !strings.Contains(strings.ToLower(contentType), "application/json") {
				utils.RespondError(c, http.StatusUnsupportedMediaType,
					utils.ErrUnsupportedMediaType, utils.MsgContentTypeRequired)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
