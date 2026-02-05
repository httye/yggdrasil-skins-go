// Package blessing_skin BlessingSkin配置管理
package blessing_skin

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"
)

// OptionsManager 配置管理�?type OptionsManager struct {
	storage *Storage
	options map[string]string // 启动时批量加载的配置
	mutex   sync.RWMutex
}

// NewOptionsManager 创建配置管理�?func NewOptionsManager(storage *Storage) *OptionsManager {
	om := &OptionsManager{
		storage: storage,
		options: make(map[string]string),
	}

	// 启动时批量加载所有Yggdrasil配置
	om.loadAllOptions()

	return om
}

// loadAllOptions 启动时批量加载所有Yggdrasil配置
func (om *OptionsManager) loadAllOptions() {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// 获取所有Yggdrasil相关配置
	var options []Option
	err := om.storage.db.Where("option_name LIKE 'ygg_%' OR option_name = 'site_url'").Find(&options).Error
	if err != nil {
		log.Printf("⚠️  Failed to load options: %v", err)
		return
	}

	// 存储到内存中
	for _, option := range options {
		om.options[option.OptionName] = option.OptionValue
	}

	log.Printf("�?Loaded %d options into memory", len(om.options))
}

// YggdrasilOptions Yggdrasil配置项及其默认值（仅包含实际存在的配置项）
var YggdrasilOptions = map[string]string{
	"ygg_uuid_algorithm":          "v3",     // UUID生成算法: v3(离线模式兼容) | v4(随机)
	"ygg_token_expire_1":          "259200", // 访问令牌过期时间（秒�?天）
	"ygg_token_expire_2":          "604800", // 刷新令牌过期时间（秒�?天）
	"ygg_tokens_limit":            "10",     // 每用户最大令牌数
	"ygg_rate_limit":              "1000",   // 速率限制（毫秒）
	"ygg_skin_domain":             "",       // 皮肤域名白名单（逗号分隔�?	"ygg_search_profile_max":      "5",      // 批量查询角色最大数�?	"ygg_private_key":             "",       // RSA私钥（PEM格式�?	"ygg_show_config_section":     "true",   // 显示配置面板
	"ygg_show_activities_section": "true",   // 显示活动面板
	"ygg_enable_ali":              "true",   // 启用ALI�?	// 注意：jwt_secret 在BlessingSkin中不存在，已移除
}

// InitializeOptions 初始化Yggdrasil配置项（只读模式，批量查询优化）
func (om *OptionsManager) InitializeOptions() error {
	// 批量查询所有需要的配置�?	optionNames := make([]string, 0, len(YggdrasilOptions))
	for optionName := range YggdrasilOptions {
		optionNames = append(optionNames, optionName)
	}

	// 一次性查询所有配置项
	var existingOptions []Option
	err := om.storage.db.Select("option_name").Where("option_name IN ?", optionNames).Find(&existingOptions).Error
	if err != nil {
		return fmt.Errorf("failed to query existing options: %w", err)
	}

	// 检查缺失的配置�?	existingMap := make(map[string]bool)
	for _, option := range existingOptions {
		existingMap[option.OptionName] = true
	}

	for optionName := range YggdrasilOptions {
		if !existingMap[optionName] {
			log.Printf("⚠️  Option '%s' not found in database, using default behavior", optionName)
		}
	}

	// 只读模式：不生成或修改RSA密钥和JWT密钥
	// BlessingSkin数据库中已有完整的Yggdrasil配置
	log.Println("�?BlessingSkin options initialized in read-only mode")
	return nil
}

// GetOption 获取配置选项（从内存读取�?func (om *OptionsManager) GetOption(name string) (string, error) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	if value, exists := om.options[name]; exists {
		return value, nil
	}

	return "", fmt.Errorf("option not found: %s", name)
}

// SetOption 设置配置选项
func (om *OptionsManager) SetOption(name, value string) error {
	var option Option
	err := om.storage.db.Where("option_name = ?", name).First(&option).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新选项
		option = Option{
			OptionName:  name,
			OptionValue: value,
		}
		return om.storage.db.Create(&option).Error
	} else if err != nil {
		return err
	}

	// 更新现有选项
	option.OptionValue = value
	return om.storage.db.Save(&option).Error
}

// GetOptionWithDefault 获取配置选项（带默认值）
func (om *OptionsManager) GetOptionWithDefault(name, defaultValue string) string {
	value, err := om.GetOption(name)
	if err != nil {
		return defaultValue
	}
	return value
}
