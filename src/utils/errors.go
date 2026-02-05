// Package utils 提供通用工具函数
package utils

import (
	"net/http"

	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"

	"github.com/gin-gonic/gin"
)

// 预定义的错误类型
const (
	ErrForbiddenOperation   = "ForbiddenOperationException"
	ErrIllegalArgument      = "IllegalArgumentException"
	ErrNotFound             = "NotFoundException"
	ErrUnauthorized         = "UnauthorizedException"
	ErrUnsupportedMediaType = "Unsupported Media Type"
)

// 预定义的错误消息
const (
	MsgInvalidToken           = "Invalid token."
	MsgInvalidCredentials     = "Invalid credentials. Invalid username or password."
	MsgTokenAlreadyHasProfile = "Access token already has a profile assigned."
	MsgPlayerNotExisted       = "Player not existed."
	MsgUserNotExisted         = "User not existed."
	MsgUserBanned             = "User has been banned."
	MsgTokenNotMatched        = "Token does not match."
	MsgEmptyCredentials       = "Username or password cannot be empty."
	MsgUnsupportedMediaType   = "Unsupported Media Type"
	MsgContentTypeRequired    = "Content-Type must be application/json"
	MsgRateLimitExceeded      = "Rate limit exceeded. Please try again later."
)

// RespondError 返回错误响应
func RespondError(c *gin.Context, statusCode int, errorType, message string) {
	errorResp := yggdrasil.ErrorResponse{
		Error:        errorType,
		ErrorMessage: message,
	}

	// 使用高性能JSON响应
	if jsonData, err := FastMarshal(errorResp); err == nil {
		c.Data(statusCode, "application/json", jsonData)
	} else {
		// 降级到标准JSON
		c.JSON(statusCode, errorResp)
	}
}

// RespondErrorWithCause 返回带原因的错误响应
func RespondErrorWithCause(c *gin.Context, statusCode int, errorType, message, cause string) {
	errorResp := yggdrasil.ErrorResponse{
		Error:        errorType,
		ErrorMessage: message,
		Cause:        cause,
	}

	// 使用高性能JSON响应
	if jsonData, err := FastMarshal(errorResp); err == nil {
		c.Data(statusCode, "application/json", jsonData)
	} else {
		// 降级到标准JSON
		c.JSON(statusCode, errorResp)
	}
}

// RespondForbiddenOperation 返回禁止操作错误
func RespondForbiddenOperation(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, ErrForbiddenOperation, message)
}

// RespondIllegalArgument 返回非法参数错误
func RespondIllegalArgument(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrIllegalArgument, message)
}

// RespondNotFound 返回未找到错�?func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, ErrNotFound, message)
}

// RespondUnauthorized 返回未授权错�?func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, ErrUnauthorized, message)
}

// RespondInvalidToken 返回无效令牌错误
func RespondInvalidToken(c *gin.Context) {
	RespondForbiddenOperation(c, MsgInvalidToken)
}

// RespondInvalidCredentials 返回无效凭据错误
func RespondInvalidCredentials(c *gin.Context) {
	RespondForbiddenOperation(c, MsgInvalidCredentials)
}

// RespondNoContent 返回无内容响�?func RespondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondJSON 返回JSON响应
func RespondJSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}
