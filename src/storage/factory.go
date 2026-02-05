// Package storage_factory 存储工厂实现
package storage_factory

import (
	"fmt"

	"github.com/httye/yggdrasil-skins-go/src/config"
	"github.com/httye/yggdrasil-skins-go/src/storage/blessing_skin"
	"github.com/httye/yggdrasil-skins-go/src/storage/file"
	storage "github.com/httye/yggdrasil-skins-go/src/storage/interface"
)

// DefaultStorageFactory 默认存储工厂
type DefaultStorageFactory struct{}

// NewStorageFactory 创建存储工厂
func NewStorageFactory() storage.StorageFactory {
	return &DefaultStorageFactory{}
}

// CreateStorage 创建存储实例
func (f *DefaultStorageFactory) CreateStorage(config *config.StorageConfig, textureConfig *config.TextureConfig) (storage.Storage, error) {
	switch config.Type {
	case "file":
		return f.createFileStorage(config, textureConfig)
	case "database":
		return f.createDatabaseStorage(config, textureConfig)
	case "blessing_skin":
		return f.createBlessingSkinStorage(config, textureConfig)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}

// GetSupportedTypes 获取支持的存储类�?func (f *DefaultStorageFactory) GetSupportedTypes() []string {
	return []string{"file", "database", "blessing_skin"}
}

// createFileStorage 创建文件存储
func (f *DefaultStorageFactory) createFileStorage(config *config.StorageConfig, textureConfig *config.TextureConfig) (storage.Storage, error) {
	options := map[string]any{
		"data_dir": config.FileOptions.DataDir,
	}
	return file.NewStorage(options, textureConfig)
}

// createDatabaseStorage 创建数据库存储（待实现）
func (f *DefaultStorageFactory) createDatabaseStorage(config *config.StorageConfig, textureConfig *config.TextureConfig) (storage.Storage, error) {
	return nil, fmt.Errorf("database storage not implemented yet")
}

// createBlessingSkinStorage 创建BlessingSkin存储
func (f *DefaultStorageFactory) createBlessingSkinStorage(config *config.StorageConfig, textureConfig *config.TextureConfig) (storage.Storage, error) {
	// 准备存储选项
	options := map[string]any{
		"database_dsn":              config.BlessingSkinOptions.DatabaseDSN,
		"debug":                     config.BlessingSkinOptions.Debug,
		"texture_base_url_override": config.BlessingSkinOptions.TextureBaseURLOverride,
		"salt":                      config.BlessingSkinOptions.Security.Salt,
		"pwd_method":                config.BlessingSkinOptions.Security.PwdMethod,
		"app_key":                   config.BlessingSkinOptions.Security.AppKey,
	}

	// 准备材质配置
	bsTextureConfig := &blessing_skin.TextureConfig{
		BaseURL: textureConfig.BaseURL,
	}

	// 创建BlessingSkin存储
	return blessing_skin.NewStorage(options, bsTextureConfig)
}
