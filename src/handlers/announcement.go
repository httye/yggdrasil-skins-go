package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/httye/yggdrasil-skins-go/src/models"
	"github.com/httye/yggdrasil-skins-go/src/utils"
)

// AnnouncementHandler 公告处理�?type AnnouncementHandler struct {
	db *gorm.DB
}

// NewAnnouncementHandler 创建公告处理�?func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{db: db}
}

// GetAnnouncements 获取公告列表
func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	isActive := c.Query("is_active")
	announcementType := c.Query("type")
	targetGroup := c.Query("target_group")
	sort := c.DefaultQuery("sort", "priority")
	order := c.DefaultQuery("order", "desc")

	// 构建查询
	query := h.db.Model(&models.Announcement{})

	// 应用筛选条�?	if isActive != "" {
		active := isActive == "true"
		if active {
			// 只获取有效的公告（激活状态且时间在有效期内）
			now := time.Now()
			query = query.Where("is_active = ? AND start_time <= ? AND (end_time IS NULL OR end_time > ?)", true, now, now)
		} else {
			query = query.Where("is_active = ?", false)
		}
	}

	if announcementType != "" {
		query = query.Where("type = ?", announcementType)
	}

	if targetGroup != "" {
		query = query.Where("target_group = ?", targetGroup)
	}

	// 应用排序
	if order == "desc" {
		query = query.Order(sort + " DESC")
	} else {
		query = query.Order(sort + " ASC")
	}

	// 执行分页查询
	var announcements []models.Announcement
	var total int64

	if err := query.Count(&total).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to count announcements")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&announcements).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch announcements")
		return
	}

	// 计算总页�?	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"data":        announcements,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetActiveAnnouncements 获取当前有效的公�?func (h *AnnouncementHandler) GetActiveAnnouncements(c *gin.Context) {
	// 获取当前用户信息（如果已登录�?	userUUID := c.GetString("user_uuid")
	var userGroup string = "all"

	if userUUID != "" {
		var user models.EnhancedUser
		if err := h.db.Where("uuid = ?", userUUID).First(&user).Error; err == nil {
			if user.IsBanned {
				userGroup = "banned"
			} else if user.IsAdmin {
				userGroup = "admins"
			} else {
				userGroup = "users"
			}
		}
	}

	now := time.Now()
	var announcements []models.ActiveAnnouncement

	// 查询有效的公�?	err := h.db.Raw(`
		SELECT 
			a.id, a.title, a.content, a.type, a.priority,
			a.target_group, a.start_time, a.end_time,
			u.username as created_by_name, a.created_at
		FROM announcements a
		LEFT JOIN users u ON a.created_by = u.uuid
		WHERE a.is_active = TRUE 
		  AND a.start_time <= ? 
		  AND (a.end_time IS NULL OR a.end_time > ?)
		  AND (a.target_group = ? OR a.target_group = 'all')
		ORDER BY a.priority DESC, a.created_at DESC
	`, now, now, userGroup).Scan(&announcements).Error

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch active announcements")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcements": announcements,
		"count":         len(announcements),
	})
}

// GetAnnouncement 获取单个公告详情
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")
	id, err := strconv.Atoi(announcementID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	var announcement models.Announcement
	if err := h.db.Preload("Creator").First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "ANNOUNCEMENT_NOT_FOUND", "Announcement not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch announcement")
		return
	}

	c.JSON(http.StatusOK, announcement)
}

// CreateAnnouncement 创建新公�?func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	adminUUID := c.GetString("user_uuid")
	if adminUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Admin authentication required")
		return
	}

	var request struct {
		Title       string     `json:"title" binding:"required,min=1,max=255"`
		Content     string     `json:"content" binding:"required,min=1"`
		Type        string     `json:"type" binding:"required,oneof=info warning error success update maintenance"`
		Priority    int        `json:"priority" binding:"min=0,max=100"`
		TargetGroup string     `json:"target_group" binding:"required,oneof=all users admins banned"`
		StartTime   *time.Time `json:"start_time"`
		EndTime     *time.Time `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// 验证时间逻辑
	if request.StartTime != nil && request.EndTime != nil {
		if request.EndTime.Before(*request.StartTime) {
			utils.RespondError(c, http.StatusBadRequest, "INVALID_TIME", "End time must be after start time")
			return
		}
	}

	// 创建公告
	announcement := models.Announcement{
		Title:       request.Title,
		Content:     request.Content,
		Type:        request.Type,
		Priority:    request.Priority,
		TargetGroup: request.TargetGroup,
		StartTime:   time.Now(), // 默认从当前时间开�?		EndTime:     request.EndTime,
		CreatedBy:   adminUUID,
	}

	if request.StartTime != nil {
		announcement.StartTime = *request.StartTime
	}

	if err := h.db.Create(&announcement).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create announcement")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Announcement created successfully",
		"announcement": announcement,
	})
}

// UpdateAnnouncement 更新公告
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	adminUUID := c.GetString("user_uuid")
	if adminUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Admin authentication required")
		return
	}

	announcementID := c.Param("id")
	id, err := strconv.Atoi(announcementID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	var request struct {
		Title       string     `json:"title" binding:"min=1,max=255"`
		Content     string     `json:"content" binding:"min=1"`
		Type        string     `json:"type" binding:"oneof=info warning error success update maintenance"`
		Priority    int        `json:"priority" binding:"min=0,max=100"`
		IsActive    *bool      `json:"is_active"`
		TargetGroup string     `json:"target_group" binding:"oneof=all users admins banned"`
		StartTime   *time.Time `json:"start_time"`
		EndTime     *time.Time `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// 查找现有公告
	var announcement models.Announcement
	if err := h.db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "ANNOUNCEMENT_NOT_FOUND", "Announcement not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch announcement")
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Content != "" {
		updates["content"] = request.Content
	}
	if request.Type != "" {
		updates["type"] = request.Type
	}
	if request.Priority != 0 || c.PostForm("priority") != "" {
		updates["priority"] = request.Priority
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}
	if request.TargetGroup != "" {
		updates["target_group"] = request.TargetGroup
	}
	if request.StartTime != nil {
		updates["start_time"] = *request.StartTime
	}
	if request.EndTime != nil {
		updates["end_time"] = *request.EndTime
	}

	// 验证时间逻辑
	if request.StartTime != nil && request.EndTime != nil {
		if request.EndTime.Before(*request.StartTime) {
			utils.RespondError(c, http.StatusBadRequest, "INVALID_TIME", "End time must be after start time")
			return
		}
	}

	if err := h.db.Model(&announcement).Updates(updates).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update announcement")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Announcement updated successfully",
		"announcement": announcement,
	})
}

// DeleteAnnouncement 删除公告
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	adminUUID := c.GetString("user_uuid")
	if adminUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Admin authentication required")
		return
	}

	announcementID := c.Param("id")
	id, err := strconv.Atoi(announcementID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	// 查找公告
	var announcement models.Announcement
	if err := h.db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "ANNOUNCEMENT_NOT_FOUND", "Announcement not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch announcement")
		return
	}

	// 删除公告
	if err := h.db.Delete(&announcement).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete announcement")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Announcement deleted successfully",
	})
}

// ToggleAnnouncementStatus 切换公告状�?func (h *AnnouncementHandler) ToggleAnnouncementStatus(c *gin.Context) {
	adminUUID := c.GetString("user_uuid")
	if adminUUID == "" {
		utils.RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Admin authentication required")
		return
	}

	announcementID := c.Param("id")
	id, err := strconv.Atoi(announcementID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "INVALID_ID", "Invalid announcement ID")
		return
	}

	// 查找公告
	var announcement models.Announcement
	if err := h.db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "ANNOUNCEMENT_NOT_FOUND", "Announcement not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch announcement")
		return
	}

	// 切换状�?	newStatus := !announcement.IsActive
	if err := h.db.Model(&announcement).Update("is_active", newStatus).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update announcement status")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Announcement status updated successfully",
		"is_active":  newStatus,
		"announcement": announcement,
	})
}

// GetAnnouncementTypes 获取公告类型列表
func (h *AnnouncementHandler) GetAnnouncementTypes(c *gin.Context) {
	types := []string{"info", "warning", "error", "success", "update", "maintenance"}
	c.JSON(http.StatusOK, gin.H{
		"types": types,
	})
}

// GetTargetGroups 获取目标用户组列�?func (h *AnnouncementHandler) GetTargetGroups(c *gin.Context) {
	groups := []string{"all", "users", "admins", "banned"}
	c.JSON(http.StatusOK, gin.H{
		"target_groups": groups,
	})
}
