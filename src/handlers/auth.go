package handlers

import (
	"strings"
	"time"

	"github.com/httye/yggdrasil-skins-go/src/cache"
	storage "github.com/httye/yggdrasil-skins-go/src/storage/interface"
	"github.com/httye/yggdrasil-skins-go/src/utils"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理�?type AuthHandler struct {
	storage      storage.Storage
	tokenCache   cache.TokenCache
	sessionCache cache.SessionCache
}

// NewAuthHandler 创建新的认证处理�?func NewAuthHandler(storage storage.Storage, tokenCache cache.TokenCache, sessionCache cache.SessionCache) *AuthHandler {
	return &AuthHandler{
		storage:      storage,
		tokenCache:   tokenCache,
		sessionCache: sessionCache,
	}
}

// Authenticate 用户登录认证
func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req yggdrasil.AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondIllegalArgument(c, "Invalid request format")
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Password == "" {
		utils.RespondIllegalArgument(c, utils.MsgEmptyCredentials)
		return
	}

	// 直接使用 AuthenticateUser 方法（已包含密码验证和单查询优化�?	user, err := h.storage.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		utils.RespondInvalidCredentials(c)
		return
	}

	// 生成客户端令�?	clientToken := req.ClientToken
	if clientToken == "" {
		clientToken = utils.GenerateRandomUUID()
	}

	// 准备可用角色列表
	availableProfiles := make([]yggdrasil.Profile, len(user.Profiles))
	copy(availableProfiles, user.Profiles)

	// 确定选中的角�?	var selectedProfile *yggdrasil.Profile
	var profileID string

	// 如果用户只有一个角色，自动选择
	if len(availableProfiles) == 1 {
		selectedProfile = &availableProfiles[0]
		profileID = selectedProfile.ID
	}

	// 如果通过角色名登录，自动选择对应角色
	if !strings.Contains(req.Username, "@") {
		for i := range availableProfiles {
			if availableProfiles[i].Name == req.Username {
				selectedProfile = &availableProfiles[i]
				profileID = selectedProfile.ID
				break
			}
		}
	}

	// 生成访问令牌（使用配置中的过期时间）
	accessToken, err := utils.GenerateJWT(user.ID, profileID, 3*24*time.Hour) // 使用默认3天有效期
	if err != nil {
		utils.RespondError(c, 500, "InternalServerError", "Failed to generate token")
		return
	}

	// 存储令牌
	token := &yggdrasil.Token{
		AccessToken: accessToken,
		ClientToken: clientToken,
		ProfileID:   profileID,
		Owner:       user.ID, // 使用用户ID而不是邮�?		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(3 * 24 * time.Hour),
	}

	if err := h.tokenCache.Store(token); err != nil {
		utils.RespondError(c, 500, "InternalServerError", "Failed to store token")
		return
	}

	// 构建响应
	response := yggdrasil.AuthenticateResponse{
		AccessToken:       accessToken,
		ClientToken:       clientToken,
		AvailableProfiles: availableProfiles,
		SelectedProfile:   selectedProfile,
	}

	// 如果请求用户信息
	if req.RequestUser {
		response.User = &yggdrasil.UserInfo{
			ID:         user.ID,
			Properties: []yggdrasil.ProfileProperty{},
		}
	}

	utils.RespondJSONFast(c, response)
}

// Refresh 刷新访问令牌
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req yggdrasil.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondIllegalArgument(c, "Invalid request format")
		return
	}

	// 获取并验证令�?	token, err := h.tokenCache.Get(req.AccessToken)
	if err != nil || !token.IsValid() {
		utils.RespondInvalidToken(c)
		return
	}

	// 验证客户端令牌（如果提供�?	if req.ClientToken != "" && token.ClientToken != req.ClientToken {
		utils.RespondForbiddenOperation(c, utils.MsgTokenNotMatched)
		return
	}

	// 获取用户信息
	user, err := h.storage.GetUserByID(token.Owner)
	if err != nil {
		utils.RespondForbiddenOperation(c, utils.MsgUserNotExisted)
		return
	}

	// 删除旧令�?	h.tokenCache.Delete(req.AccessToken)

	// 确定新令牌的角色绑定
	profileID := token.ProfileID
	var selectedProfile *yggdrasil.Profile

	// 如果请求中指定了角色选择
	if req.SelectedProfile != nil {
		// 验证角色是否属于用户
		found := false
		for _, profile := range user.Profiles {
			if profile.ID == req.SelectedProfile.ID {
				selectedProfile = req.SelectedProfile
				profileID = req.SelectedProfile.ID
				found = true
				break
			}
		}
		if !found {
			utils.RespondForbiddenOperation(c, "Selected profile does not belong to user")
			return
		}
	} else if profileID != "" {
		// 保持原有的角色绑�?		for _, profile := range user.Profiles {
			if profile.ID == profileID {
				selectedProfile = &profile
				break
			}
		}
	}

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateJWT(user.ID, profileID, 3*24*time.Hour)
	if err != nil {
		utils.RespondError(c, 500, "InternalServerError", "Failed to generate token")
		return
	}

	// 存储新令�?	newToken := &yggdrasil.Token{
		AccessToken: newAccessToken,
		ClientToken: token.ClientToken,
		ProfileID:   profileID,
		Owner:       user.ID, // 使用用户ID而不是邮�?		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(3 * 24 * time.Hour),
	}

	if err := h.tokenCache.Store(newToken); err != nil {
		utils.RespondError(c, 500, "InternalServerError", "Failed to store token")
		return
	}

	// 构建响应
	response := yggdrasil.RefreshResponse{
		AccessToken:     newAccessToken,
		ClientToken:     token.ClientToken,
		SelectedProfile: selectedProfile,
	}

	// 如果请求用户信息
	if req.RequestUser {
		response.User = &yggdrasil.UserInfo{
			ID:         user.ID,
			Properties: []yggdrasil.ProfileProperty{},
		}
	}

	utils.RespondJSONFast(c, response)
}

// Validate 验证令牌
func (h *AuthHandler) Validate(c *gin.Context) {
	var req yggdrasil.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondIllegalArgument(c, "Invalid request format")
		return
	}

	// 获取并验证令�?	token, err := h.tokenCache.Get(req.AccessToken)
	if err != nil || !token.IsValid() {
		utils.RespondInvalidToken(c)
		return
	}

	// 验证客户端令牌（如果提供�?	if req.ClientToken != "" && token.ClientToken != req.ClientToken {
		utils.RespondForbiddenOperation(c, utils.MsgTokenNotMatched)
		return
	}

	// 令牌有效，返�?04
	utils.RespondNoContent(c)
}

// Invalidate 撤销令牌
func (h *AuthHandler) Invalidate(c *gin.Context) {
	var req yggdrasil.InvalidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondIllegalArgument(c, "Invalid request format")
		return
	}

	// 删除令牌（无论是否存在都返回204�?	h.tokenCache.Delete(req.AccessToken)
	utils.RespondNoContent(c)
}

// Signout 全局登出
func (h *AuthHandler) Signout(c *gin.Context) {
	var req yggdrasil.SignoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondIllegalArgument(c, "Invalid request format")
		return
	}

	// 验证用户凭据（使用统一的认证方法）
	user, err := h.storage.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		utils.RespondInvalidCredentials(c)
		return
	}

	// 删除用户的所有令�?	h.tokenCache.DeleteUserTokens(user.ID)
	utils.RespondNoContent(c)
}
