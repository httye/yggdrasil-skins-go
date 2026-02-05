package models

import (
	"time"

	"gorm.io/gorm"
)

// EnhancedUser 扩展用户模型
type EnhancedUser struct {
	UUID                   string     `gorm:"primaryKey;type:varchar(36)" json:"uuid"`
	Email                  string     `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Username               string     `gorm:"uniqueIndex;type:varchar(255);not null" json:"username"`
	Password               string     `gorm:"type:varchar(255);not null" json:"-"`
	PrimaryPlayerName      string     `gorm:"uniqueIndex;type:varchar(255)" json:"primary_player_name,omitempty"`
	PlayerUUID             string     `gorm:"type:varchar(36);index" json:"player_uuid,omitempty"`
	QQNumber               string     `gorm:"type:varchar(20)" json:"qq_number,omitempty"`
	EmailVerified          bool       `gorm:"default:false;index" json:"email_verified"`
	EmailVerificationToken string     `gorm:"type:varchar(255)" json:"-"`
	AgreedToTerms          bool       `gorm:"default:false" json:"agreed_to_terms"`
	RegistrationIP         string     `gorm:"type:varchar(45)" json:"registration_ip,omitempty"`
	LastLoginIP            string     `gorm:"type:varchar(45)" json:"last_login_ip,omitempty"`
	LastLoginAt            *time.Time `json:"last_login_at,omitempty"`
	MaxProfiles            int        `gorm:"default:5" json:"max_profiles"`
	IsBanned               bool       `gorm:"default:false;index" json:"is_banned"`
	BannedReason           string     `gorm:"type:text" json:"banned_reason,omitempty"`
	BannedAt               *time.Time `gorm:"index" json:"banned_at,omitempty"`
	BannedBy               string     `gorm:"type:varchar(36)" json:"banned_by,omitempty"`
	IsAdmin                bool       `gorm:"default:false;index" json:"is_admin"`
	PermissionGroupID      int        `gorm:"default:1;index" json:"permission_group_id"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`

	// 关联关系
	Profiles          []Profile          `gorm:"foreignKey:UserUUID;references:UUID" json:"profiles,omitempty"`
	UploadedSkins     []Skin             `gorm:"foreignKey:UploaderUUID;references:UUID" json:"uploaded_skins,omitempty"`
	UploadedCapes     []Cape             `gorm:"foreignKey:UploaderUUID;references:UUID" json:"uploaded_capes,omitempty"`
	PermissionGroup   PermissionGroup    `gorm:"foreignKey:PermissionGroupID" json:"permission_group,omitempty"`
	AdminLogs         []AdminLog         `gorm:"foreignKey:AdminUUID;references:UUID" json:"admin_logs,omitempty"`
	TargetLogs        []AdminLog         `gorm:"foreignKey:TargetUserUUID;references:UUID" json:"target_logs,omitempty"`
	UserLogs          []UserLog          `gorm:"foreignKey:UserUUID;references:UUID" json:"user_logs,omitempty"`
}

// TableName 设置表名
func (EnhancedUser) TableName() string {
	return "users"
}

// PermissionGroup 权限组模型
type PermissionGroup struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Permissions JSONMap   `gorm:"type:json;not null" json:"permissions"`
	IsDefault   bool      `gorm:"default:false;index" json:"is_default"`
	IsSystem    bool      `gorm:"default:false;index" json:"is_system"`
	Priority    int       `gorm:"default:0;index" json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Users []EnhancedUser `gorm:"foreignKey:PermissionGroupID" json:"users,omitempty"`
}

// TableName 设置表名
func (PermissionGroup) TableName() string {
	return "permission_groups"
}

// Profile 角色模型（Minecraft游戏角色）
type Profile struct {
	UUID      string     `gorm:"primaryKey;type:varchar(36)" json:"uuid"`
	Name      string     `gorm:"uniqueIndex;type:varchar(255);not null" json:"name"`
	UserUUID  string     `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	SkinID    *int       `gorm:"index" json:"skin_id,omitempty"`
	CapeID    *int       `gorm:"index" json:"cape_id,omitempty"`
	IsActive  bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// 关联关系
	User     EnhancedUser `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
	Skin     *Skin        `gorm:"foreignKey:SkinID" json:"skin,omitempty"`
	Cape     *Cape        `gorm:"foreignKey:CapeID" json:"cape,omitempty"`
}

// TableName 设置表名
func (Profile) TableName() string {
	return "profiles"
}

// Skin 皮肤模型
type Skin struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID           string    `gorm:"uniqueIndex;type:varchar(36);not null" json:"uuid"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Hash           string    `gorm:"uniqueIndex;type:varchar(64);not null" json:"hash"`
	Type           string    `gorm:"type:enum('steve','alex');default:'steve'" json:"type"`
	ModelType      string    `gorm:"type:enum('default','slim');default:'default'" json:"model_type"`
	FilePath       string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize       int       `gorm:"not null" json:"file_size"`
	UploadTime     time.Time `json:"upload_time"`
	UploaderUUID   string    `gorm:"type:varchar(36);not null;index" json:"uploader_uuid"`
	IsPublic       bool      `gorm:"default:false;index" json:"is_public"`
	DownloadCount  int       `gorm:"default:0" json:"download_count"`
	LikesCount     int       `gorm:"default:0" json:"likes_count"`
	IsVerified     bool      `gorm:"default:false;index" json:"is_verified"`
	VerifiedBy     *string   `gorm:"type:varchar(36)" json:"verified_by,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联关系
	Uploader   EnhancedUser `gorm:"foreignKey:UploaderUUID;references:UUID" json:"uploader,omitempty"`
	Verifier   *EnhancedUser `gorm:"foreignKey:VerifiedBy;references:UUID" json:"verifier,omitempty"`
	Profiles   []Profile    `gorm:"foreignKey:SkinID" json:"profiles,omitempty"`
	Tags       []SkinTag    `gorm:"foreignKey:SkinID" json:"tags,omitempty"`
	Likes      []SkinLike   `gorm:"foreignKey:SkinID" json:"likes,omitempty"`
}

// TableName 设置表名
func (Skin) TableName() string {
	return "skins"
}

// Cape 披风模型
type Cape struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID          string    `gorm:"uniqueIndex;type:varchar(36);not null" json:"uuid"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Hash          string    `gorm:"uniqueIndex;type:varchar(64);not null" json:"hash"`
	FilePath      string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize      int       `gorm:"not null" json:"file_size"`
	UploadTime    time.Time `json:"upload_time"`
	UploaderUUID  string    `gorm:"type:varchar(36);not null;index" json:"uploader_uuid"`
	IsPublic      bool      `gorm:"default:false;index" json:"is_public"`
	DownloadCount int       `gorm:"default:0" json:"download_count"`
	IsVerified    bool      `gorm:"default:false;index" json:"is_verified"`
	VerifiedBy    *string   `gorm:"type:varchar(36)" json:"verified_by,omitempty"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联关系
	Uploader  EnhancedUser `gorm:"foreignKey:UploaderUUID;references:UUID" json:"uploader,omitempty"`
	Verifier  *EnhancedUser `gorm:"foreignKey:VerifiedBy;references:UUID" json:"verifier,omitempty"`
	Profiles  []Profile    `gorm:"foreignKey:CapeID" json:"profiles,omitempty"`
}

// TableName 设置表名
func (Cape) TableName() string {
	return "capes"
}

// Announcement 公告模型
type Announcement struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Type        string     `gorm:"type:enum('info','warning','error','success','update','maintenance');default:'info'" json:"type"`
	Priority    int        `gorm:"default:0;index" json:"priority"`
	IsActive    bool       `gorm:"default:true;index" json:"is_active"`
	TargetGroup string     `gorm:"type:enum('all','users','admins','banned');default:'all'" json:"target_group"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	CreatedBy   string     `gorm:"type:varchar(36);not null;index" json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联关系
	Creator EnhancedUser `gorm:"foreignKey:CreatedBy;references:UUID" json:"creator,omitempty"`
}

// TableName 设置表名
func (Announcement) TableName() string {
	return "announcements"
}

// AdminLog 管理员操作日志模型
type AdminLog struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminUUID      string     `gorm:"type:varchar(36);not null;index" json:"admin_uuid"`
	Action         string     `gorm:"type:varchar(100);not null;index" json:"action"`
	TargetUserUUID *string    `gorm:"type:varchar(36);index" json:"target_user_uuid,omitempty"`
	Details        JSONMap    `gorm:"type:json" json:"details,omitempty"`
	IPAddress      string     `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent      string     `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`

	// 关联关系
	Admin      EnhancedUser  `gorm:"foreignKey:AdminUUID;references:UUID" json:"admin,omitempty"`
	TargetUser *EnhancedUser `gorm:"foreignKey:TargetUserUUID;references:UUID" json:"target_user,omitempty"`
}

// TableName 设置表名
func (AdminLog) TableName() string {
	return "admin_logs"
}

// ServerStatus 服务器状态模型
type ServerStatus struct {
	ID           int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerName   string     `gorm:"type:varchar(255);not null" json:"server_name"`
	ServerType   string     `gorm:"type:enum('survival','creative','minigames','auth','lobby','bedwars','skywars');not null" json:"server_type"`
	ServerAddress string    `gorm:"type:varchar(255);not null" json:"server_address"`
	ServerPort   int        `gorm:"default:25565" json:"server_port"`
	Status       string     `gorm:"type:enum('online','offline','maintenance','starting','stopping');default:'offline';index" json:"status"`
	PlayerCount  int        `gorm:"default:0" json:"player_count"`
	MaxPlayers   int        `gorm:"default:0" json:"max_players"`
	MOTD         string     `gorm:"type:text" json:"motd"`
	Version      string     `gorm:"type:varchar(100)" json:"version"`
	TPS          float32    `gorm:"default:20.0" json:"tps"`
	UptimeSeconds int64     `gorm:"default:0" json:"uptime_seconds"`
	MemoryUsed   int64      `gorm:"default:0" json:"memory_used"`
	MemoryMax    int64      `gorm:"default:0" json:"memory_max"`
	CPUUsage     float32    `gorm:"default:0.0" json:"cpu_usage"`
	LastPing     time.Time  `json:"last_ping"`
	NextPing     time.Time  `json:"next_ping"`
	IsMonitoring bool       `gorm:"default:true;index" json:"is_monitoring"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 设置表名
func (ServerStatus) TableName() string {
	return "server_status"
}

// SkinTag 皮肤标签模型
type SkinTag struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SkinID    int       `gorm:"not null;index" json:"skin_id"`
	TagName   string    `gorm:"type:varchar(50);not null" json:"tag_name"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	Skin Skin `gorm:"foreignKey:SkinID" json:"skin,omitempty"`
}

// TableName 设置表名
func (SkinTag) TableName() string {
	return "skin_tags"
}

// SkinLike 皮肤点赞模型
type SkinLike struct {
	SkinID    int       `gorm:"primaryKey" json:"skin_id"`
	UserUUID  string    `gorm:"primaryKey;type:varchar(36)" json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	Skin Skin         `gorm:"foreignKey:SkinID" json:"skin,omitempty"`
	User EnhancedUser `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

// TableName 设置表名
func (SkinLike) TableName() string {
	return "skin_likes"
}

// JSONMap JSON字段类型
type JSONMap map[string]interface{}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return fmt.Errorf("cannot scan type %T into JSONMap", value)
	}
}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (interface{}, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// UserLog 用户操作日志模型
type UserLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID  string    `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	Action    string    `gorm:"type:varchar(100);not null;index" json:"action"`
	Details   JSONMap   `gorm:"type:json" json:"details,omitempty"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent string    `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	User EnhancedUser `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

// TableName 设置表名
func (UserLog) TableName() string {
	return "user_logs"
}

// UserFullInfo 用户完整信息视图
type UserFullInfo struct {
	UUID                   string     `json:"uuid"`
	Email                  string     `json:"email"`
	Username               string     `json:"username"`
	PrimaryPlayerName      string     `json:"primary_player_name,omitempty"`
	PlayerUUID             string     `json:"player_uuid,omitempty"`
	QQNumber               string     `json:"qq_number,omitempty"`
	EmailVerified          bool       `json:"email_verified"`
	AgreedToTerms          bool       `json:"agreed_to_terms"`
	RegistrationIP         string     `json:"registration_ip,omitempty"`
	LastLoginIP            string     `json:"last_login_ip,omitempty"`
	LastLoginAt            *time.Time `json:"last_login_at,omitempty"`
	MaxProfiles            int        `json:"max_profiles"`
	IsBanned               bool       `json:"is_banned"`
	BannedReason           string     `json:"banned_reason,omitempty"`
	BannedAt               *time.Time `json:"banned_at,omitempty"`
	IsAdmin                bool       `json:"is_admin"`
	PermissionGroupID      int        `json:"permission_group_id"`
	PermissionGroupName    string     `json:"permission_group_name"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	CurrentProfiles        int        `json:"current_profiles"`
	UploadedSkins          int        `json:"uploaded_skins"`
	UploadedCapes          int        `json:"uploaded_capes"`
}

// SkinStatistics 皮肤统计视图
type SkinStatistics struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	ModelType    string    `json:"model_type"`
	UploadTime   time.Time `json:"upload_time"`
	UploaderUUID string    `json:"uploader_uuid"`
	UploaderName string    `json:"uploader_name"`
	IsPublic     bool      `json:"is_public"`
	DownloadCount int      `json:"download_count"`
	LikesCount    int       `json:"likes_count"`
	IsVerified    bool      `json:"is_verified"`
	ActiveUsers   int       `json:"active_users"`
	Tags          string    `json:"tags,omitempty"`
}

// ActiveAnnouncement 有效公告视图
type ActiveAnnouncement struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Type         string    `json:"type"`
	Priority     int       `json:"priority"`
	TargetGroup  string    `json:"target_group"`
	StartTime    time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	CreatedByName string   `json:"created_by_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// Common queries and scopes
func (u *EnhancedUser) BeforeCreate(tx *gorm.DB) error {
	if u.MaxProfiles == 0 {
		u.MaxProfiles = 5 // 默认角色数量限制
	}
	if u.PermissionGroupID == 0 {
		u.PermissionGroupID = 1 // 默认权限组
	}
	return nil
}

// IsProfileLimitReached 检查是否达到角色数量限制
func (u *EnhancedUser) IsProfileLimitReached(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Model(&Profile{}).Where("user_uuid = ? AND is_active = ?", u.UUID, true).Count(&count).Error
	if err != nil {
		return false, err
	}
	
	if u.MaxProfiles == -1 {
		return false, nil // 无限制
	}
	
	return int(count) >= u.MaxProfiles, nil
}

// GetActiveProfileCount 获取活跃角色数量
func (u *EnhancedUser) GetActiveProfileCount(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&Profile{}).Where("user_uuid = ? AND is_active = ?", u.UUID, true).Count(&count).Error
	return int(count), err
}

// CanCreateProfile 检查是否可以创建新角色
func (u *EnhancedUser) CanCreateProfile(db *gorm.DB) (bool, int, error) {
	currentCount, err := u.GetActiveProfileCount(db)
	if err != nil {
		return false, 0, err
	}
	
	if u.IsBanned {
		return false, currentCount, nil
	}
	
	if u.MaxProfiles == -1 {
		return true, currentCount, nil // 无限制
	}
	
	return currentCount < u.MaxProfiles, currentCount, nil
}

// BanUser 封禁用户
func (u *EnhancedUser) BanUser(reason string, adminUUID string, db *gorm.DB) error {
	now := time.Now()
	u.IsBanned = true
	u.BannedReason = reason
	u.BannedAt = &now
	u.BannedBy = adminUUID
	
	// 禁用所有角色
	if err := db.Model(&Profile{}).Where("user_uuid = ?", u.UUID).Update("is_active", false).Error; err != nil {
		return err
	}
	
	return db.Save(u).Error
}

// UnbanUser 解封用户
func (u *EnhancedUser) UnbanUser(db *gorm.DB) error {
	u.IsBanned = false
	u.BannedReason = ""
	u.BannedAt = nil
	u.BannedBy = ""
	return db.Save(u).Error
}

// HasPermission 检查用户权限
func (u *EnhancedUser) HasPermission(permission string) bool {
	if u.IsAdmin {
		return true // 管理员拥有所有权限
	}
	
	// 这里应该检查权限组的权限配置
	// 简化实现，实际需要解析u.PermissionGroup.Permissions
	return true
}

// GetUserFullInfo 获取用户完整信息
func GetUserFullInfo(db *gorm.DB, userUUID string) (*UserFullInfo, error) {
	var userInfo UserFullInfo
	err := db.Raw(`
		SELECT 
			u.uuid, u.email, u.username, u.primary_player_name, u.player_uuid, 
			u.qq_number, u.email_verified, u.agreed_to_terms, 
			u.registration_ip, u.last_login_ip, u.last_login_at,
			u.max_profiles, u.is_banned, u.banned_reason, u.banned_at, 
			u.is_admin, u.permission_group_id, pg.name as permission_group_name, 
			u.created_at, u.updated_at,
			(SELECT COUNT(*) FROM profiles p WHERE p.user_uuid = u.uuid AND p.is_active = TRUE) as current_profiles,
			(SELECT COUNT(*) FROM skins s WHERE s.uploader_uuid = u.uuid) as uploaded_skins,
			(SELECT COUNT(*) FROM capes c WHERE c.uploader_uuid = u.uuid) as uploaded_capes
		FROM users u
		LEFT JOIN permission_groups pg ON u.permission_group_id = pg.id
		WHERE u.uuid = ?
	`, userUUID).Scan(&userInfo).Error
	
	return &userInfo, err
}

// GetUserByPlayerName 通过游戏名获取用户信息
func GetUserByPlayerName(db *gorm.DB, playerName string) (*EnhancedUser, error) {
	var user EnhancedUser
	err := db.Where("primary_player_name = ?", playerName).First(&user).Error
	return &user, err
}

// GetUserByPlayerUUID 通过游戏UUID获取用户信息
func GetUserByPlayerUUID(db *gorm.DB, playerUUID string) (*EnhancedUser, error) {
	var user EnhancedUser
	err := db.Where("player_uuid = ?", playerUUID).First(&user).Error
	return &user, err
}

// LogUserAction 记录用户操作日志
func LogUserAction(db *gorm.DB, userUUID, action string, details JSONMap, ipAddress, userAgent string) error {
	logEntry := UserLog{
		UserUUID:  userUUID,
		Action:    action,
		Details:   details,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	return db.Create(&logEntry).Error
}

// GetUserLogs 获取用户操作日志
func GetUserLogs(db *gorm.DB, userUUID string, limit int) ([]UserLog, error) {
	var logs []UserLog
	err := db.Where("user_uuid = ?", userUUID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// UpdateLastLoginInfo 更新用户最后登录信息
func (u *EnhancedUser) UpdateLastLoginInfo(db *gorm.DB, ipAddress string) error {
	now := time.Now()
	u.LastLoginIP = ipAddress
	u.LastLoginAt = &now
	return db.Save(u).Error
}