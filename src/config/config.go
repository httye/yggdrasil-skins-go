// Package config 提供配置管理
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Auth       AuthConfig       `yaml:"auth"`
	Rate       RateConfig       `yaml:"rate"`
	Storage    StorageConfig    `yaml:"storage"`
	Cache      CacheConfig      `yaml:"cache"`
	Texture    TextureConfig    `yaml:"texture"`
	Yggdrasil  YggdrasilConfig  `yaml:"yggdrasil"`
	Middleware MiddlewareConfig `yaml:"middleware"`
	Logging    LoggingConfig    `yaml:"logging"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Security   SecurityConfig   `yaml:"security"`
	Warmup     WarmupConfig     `yaml:"warmup"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type                string                     `yaml:"type"`                 // 存储类型：memory, file, database, blessing_skin
	MemoryOptions       MemoryStorageOptions       `yaml:"memory_options"`       // 内存存储选项
	FileOptions         FileStorageOptions         `yaml:"file_options"`         // 文件存储选项
	DatabaseOptions     DatabaseStorageOptions     `yaml:"database_options"`     // 数据库存储选项
	BlessingSkinOptions BlessingSkinStorageOptions `yaml:"blessingskin_options"` // BlessingSkin存储选项
}

// MemoryStorageOptions 内存存储选项
type MemoryStorageOptions struct {
	// 内存存储暂无特殊配置
}

// FileStorageOptions 文件存储选项
type FileStorageOptions struct {
	DataDir string `yaml:"data_dir"` // 数据目录
}

// DatabaseStorageOptions 数据库存储选项
type DatabaseStorageOptions struct {
	DatabaseDSN string `yaml:"database_dsn"` // 数据库连接字符串
	Debug       bool   `yaml:"debug"`        // 调试模式
}

// BlessingSkinStorageOptions BlessingSkin存储选项
type BlessingSkinStorageOptions struct {
	DatabaseDSN            string               `yaml:"database_dsn"`              // MySQL连接字符�?	Debug                  bool                 `yaml:"debug"`                     // 调试模式
	TextureBaseURLOverride bool                 `yaml:"texture_base_url_override"` // 为true时使用配置文件的texture.base_url而不是options中的site_url
	Security               BlessingSkinSecurity `yaml:"security"`                  // 安全配置
}

// BlessingSkinSecurity BlessingSkin安全配置
type BlessingSkinSecurity struct {
	Salt      string `yaml:"salt"`       // 密码加密盐�?(对应BlessingSkin的SALT)
	PwdMethod string `yaml:"pwd_method"` // 密码加密方法 (对应BlessingSkin的PWD_METHOD)
	AppKey    string `yaml:"app_key"`    // 应用密钥 (对应BlessingSkin的APP_KEY)
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Token    CacheBackendConfig  `yaml:"token"`    // Token缓存配置
	Session  CacheBackendConfig  `yaml:"session"`  // Session缓存配置
	Response ResponseCacheConfig `yaml:"response"` // 响应缓存配置
	User     UserCacheConfig     `yaml:"user"`     // 用户缓存配置
}

// CacheBackendConfig 缓存后端配置
type CacheBackendConfig struct {
	Type    string         `yaml:"type"`    // 缓存类型：memory, redis, file, database
	Options map[string]any `yaml:"options"` // 缓存选项
}

// ResponseCacheConfig 响应缓存配置
type ResponseCacheConfig struct {
	Enabled          bool          `yaml:"enabled"`           // 是否启用响应缓存
	APIMetadata      bool          `yaml:"api_metadata"`      // 是否缓存API元数�?	ErrorResponses   bool          `yaml:"error_responses"`   // 是否缓存错误响应
	ProfileResponses bool          `yaml:"profile_responses"` // 是否缓存角色响应
	CacheDuration    time.Duration `yaml:"cache_duration"`    // 缓存持续时间
	MaxCacheSize     int           `yaml:"max_cache_size"`    // 最大缓存条目数
}

// UserCacheConfig 用户缓存配置
type UserCacheConfig struct {
	Enabled         bool          `yaml:"enabled"`          // 是否启用用户缓存
	Duration        time.Duration `yaml:"duration"`         // 缓存持续时间
	MaxUsers        int           `yaml:"max_users"`        // 最大缓存用户数
	CleanupInterval time.Duration `yaml:"cleanup_interval"` // 清理间隔
}

// TextureConfig 材质配置
type TextureConfig struct {
	BaseURL       string   `yaml:"base_url"`       // 材质基础URL
	UploadEnabled bool     `yaml:"upload_enabled"` // 是否启用上传
	MaxFileSize   int64    `yaml:"max_file_size"`  // 最大文件大小（字节�?	AllowedTypes  []string `yaml:"allowed_types"`  // 允许的文件类�?}

// MiddlewareConfig 中间件配�?type MiddlewareConfig struct {
	CORS        CORSConfig        `yaml:"cors"`        // CORS配置
	RateLimit   RateLimitConfig   `yaml:"rate_limit"`  // 速率限制配置
	Performance PerformanceConfig `yaml:"performance"` // 性能监控配置
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`           // 是否启用CORS
	AllowedOrigins   []string `yaml:"allowed_origins"`   // 允许的源
	AllowedMethods   []string `yaml:"allowed_methods"`   // 允许的方�?	AllowedHeaders   []string `yaml:"allowed_headers"`   // 允许的头�?	ExposedHeaders   []string `yaml:"exposed_headers"`   // 暴露的头�?	AllowCredentials bool     `yaml:"allow_credentials"` // 是否允许凭证
	MaxAge           int      `yaml:"max_age"`           // 预检请求缓存时间
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Enabled      bool          `yaml:"enabled"`       // 是否启用速率限制
	AuthInterval time.Duration `yaml:"auth_interval"` // 认证请求间隔
	GlobalLimit  int           `yaml:"global_limit"`  // 全局请求限制（每分钟�?	BurstLimit   int           `yaml:"burst_limit"`   // 突发请求限制
}

// PerformanceConfig 性能监控配置
type PerformanceConfig struct {
	Enabled         bool `yaml:"enabled"`          // 是否启用性能监控
	DetailedMetrics bool `yaml:"detailed_metrics"` // 是否启用详细指标
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`       // 日志级别：debug, info, warn, error
	Format     string `yaml:"format"`      // 日志格式：text, json
	Output     string `yaml:"output"`      // 输出目标：stdout, stderr, file
	File       string `yaml:"file"`        // 日志文件路径（当output为file时）
	MaxSize    int    `yaml:"max_size"`    // 日志文件最大大小（MB�?	MaxBackups int    `yaml:"max_backups"` // 保留的日志文件数�?	MaxAge     int    `yaml:"max_age"`     // 日志文件保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧日志文�?}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled         bool   `yaml:"enabled"`          // 是否启用监控
	MetricsEndpoint string `yaml:"metrics_endpoint"` // 监控端点路径
	CacheStats      bool   `yaml:"cache_stats"`      // 是否启用缓存统计
	DBStats         bool   `yaml:"db_stats"`         // 是否启用数据库统�?	SystemStats     bool   `yaml:"system_stats"`     // 是否启用系统统计
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	MaxRequestSize string        `yaml:"max_request_size"` // 最大请求大�?	ReadTimeout    time.Duration `yaml:"read_timeout"`     // 读取超时
	WriteTimeout   time.Duration `yaml:"write_timeout"`    // 写入超时
	IdleTimeout    time.Duration `yaml:"idle_timeout"`     // 空闲超时
}

// WarmupConfig 预热配置
type WarmupConfig struct {
	Enabled       bool `yaml:"enabled"`        // 是否启用预热
	APIMetadata   bool `yaml:"api_metadata"`   // 是否预热API元数�?	ErrorCache    bool `yaml:"error_cache"`    // 是否预热错误缓存
	UUIDCache     bool `yaml:"uuid_cache"`     // 是否预热UUID缓存
	ProfileCache  bool `yaml:"profile_cache"`  // 是否预热角色缓存
	ConcurrentNum int  `yaml:"concurrent_num"` // 并发预热数量
}

// ServerConfig 服务器配�?type ServerConfig struct {
	Host    string `yaml:"host"`     // 监听地址
	Port    int    `yaml:"port"`     // 监听端口
	Debug   bool   `yaml:"debug"`    // 调试模式
	BaseURL string `yaml:"base_url"` // API基础路径，如 "/api/yggdrasil"
}

// AuthConfig 认证配置
type AuthConfig struct {
	TokenExpiration     time.Duration `yaml:"token_expiration"`     // 令牌过期时间
	JWTSecret           string        `yaml:"jwt_secret"`           // JWT密钥
	TokensLimit         int           `yaml:"tokens_limit"`         // 每用户令牌数量限�?	RequireVerification bool          `yaml:"require_verification"` // 是否需要邮箱验�?}

// RateConfig 速率限制配置
type RateConfig struct {
	AuthInterval time.Duration `yaml:"auth_interval"` // 认证请求间隔
	Enabled      bool          `yaml:"enabled"`       // 是否启用速率限制
}

// YggdrasilConfig Yggdrasil相关配置
type YggdrasilConfig struct {
	Meta        MetaConfig     `yaml:"meta"`         // 元数据配�?	SkinDomains []string       `yaml:"skin_domains"` // 皮肤域名白名�?	Keys        KeysConfig     `yaml:"keys"`         // 密钥配置
	Features    FeaturesConfig `yaml:"features"`     // 功能配置
}

// MetaConfig 元数据配�?type MetaConfig struct {
	ServerName            string            `yaml:"server_name"`            // 服务器名�?	ImplementationName    string            `yaml:"implementation_name"`    // 实现名称
	ImplementationVersion string            `yaml:"implementation_version"` // 实现版本
	Links                 map[string]string `yaml:"links"`                  // 相关链接
}

// KeysConfig 密钥配置
type KeysConfig struct {
	PrivateKeyPath string `yaml:"private_key_path"` // RSA私钥文件路径
	PublicKeyPath  string `yaml:"public_key_path"`  // RSA公钥文件路径
}

// FeaturesConfig 功能配置
type FeaturesConfig struct {
	NonEmailLogin bool `yaml:"non_email_login"` // 支持非邮箱登�?}

// LoadConfig 从文件加载配�?func LoadConfig(filename string) (*Config, error) {
	// 如果配置文件不存在，创建默认配置文件
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultConfig := DefaultConfig()
		if err := SaveConfig(defaultConfig, filename); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		fmt.Printf("Created default config file: %s\n", filename)
		fmt.Println("Please review and modify the configuration, then restart the server.")
		os.Exit(0)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置到文�?func SaveConfig(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配�?	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// 验证BaseURL格式
	if c.Server.BaseURL != "" {
		if !strings.HasPrefix(c.Server.BaseURL, "/") {
			return fmt.Errorf("base_url must start with '/', got: %s", c.Server.BaseURL)
		}
		// 去除末尾的斜杠（除非是根路径�?		if c.Server.BaseURL != "/" && strings.HasSuffix(c.Server.BaseURL, "/") {
			c.Server.BaseURL = strings.TrimSuffix(c.Server.BaseURL, "/")
		}
	}

	// 验证JWT密钥
	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// 验证密钥文件路径（对于BlessingSkin存储，允许为空）
	if c.Storage.Type != "blessing_skin" {
		if c.Yggdrasil.Keys.PrivateKeyPath == "" || c.Yggdrasil.Keys.PublicKeyPath == "" {
			return fmt.Errorf("key file paths cannot be empty for non-blessing_skin storage")
		}
	}

	// 验证皮肤域名配置
	for _, domain := range c.Yggdrasil.SkinDomains {
		if err := validateDomainOrCIDR(domain); err != nil {
			return fmt.Errorf("invalid skin domain '%s': %w", domain, err)
		}
	}

	return nil
}

// validateDomainOrCIDR 验证域名格式
func validateDomainOrCIDR(input string) error {
	// 检查域名格式（简单验证）
	if input == "" {
		return fmt.Errorf("empty domain")
	}

	// 允许通配符域名（�?开头）
	if strings.HasPrefix(input, ".") {
		if len(input) < 2 {
			return fmt.Errorf("invalid wildcard domain")
		}
	}

	return nil
}

// IsAllowedSkinDomain 检查域名是否在皮肤白名单中
func (c *Config) IsAllowedSkinDomain(domain string) bool {
	// 如果白名单为空，允许所有域�?	if len(c.Yggdrasil.SkinDomains) == 0 {
		return true
	}

	for _, allowed := range c.Yggdrasil.SkinDomains {
		// 检查精确匹�?		if domain == allowed {
			return true
		}

		// 检查通配符匹配（�?开头的规则�?		if strings.HasPrefix(allowed, ".") {
			if strings.HasSuffix(domain, allowed) {
				return true
			}
		}
	}

	return false
}

// GetBaseURL 根据请求获取基础URL
func (c *Config) GetBaseURL(host string) string {
	if host == "" {
		host = fmt.Sprintf("localhost:%d", c.Server.Port)
	}
	return fmt.Sprintf("http://%s", host)
}

// GetLinkURL 获取链接URL，如果配置中没有则使用动态生�?func (c *Config) GetLinkURL(linkType, host string) string {
	if url, exists := c.Yggdrasil.Meta.Links[linkType]; exists && url != "" {
		return url
	}

	baseURL := c.GetBaseURL(host)
	switch linkType {
	case "homepage":
		return baseURL
	case "register":
		return baseURL + "/register"
	default:
		return baseURL
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:    "0.0.0.0",
			Port:    8080,
			Debug:   false,
			BaseURL: "", // 默认为空，表示不使用基础路径
		},
		Auth: AuthConfig{
			TokenExpiration:     3 * 24 * time.Hour, // 3�?			JWTSecret:           "yggdrasil-api-secret-key-change-in-production",
			TokensLimit:         10,
			RequireVerification: false,
		},
		Rate: RateConfig{
			AuthInterval: 1 * time.Second, // 1秒间�?			Enabled:      true,
		},
		Storage: StorageConfig{
			Type:          "memory",
			MemoryOptions: MemoryStorageOptions{},
			FileOptions: FileStorageOptions{
				DataDir: "data",
			},
			DatabaseOptions: DatabaseStorageOptions{
				DatabaseDSN: "",
				Debug:       false,
			},
			BlessingSkinOptions: BlessingSkinStorageOptions{
				DatabaseDSN:            "",
				Debug:                  false,
				TextureBaseURLOverride: false,
				Security: BlessingSkinSecurity{
					Salt:      "blessing_skin_salt",
					PwdMethod: "BCRYPT",
					AppKey:    "base64:your_app_key_here",
				},
			},
		},
		Cache: CacheConfig{
			Token: CacheBackendConfig{
				Type:    "memory",
				Options: map[string]any{},
			},
			Session: CacheBackendConfig{
				Type:    "memory",
				Options: map[string]any{},
			},
			Response: ResponseCacheConfig{
				Enabled:          true,
				APIMetadata:      true,
				ErrorResponses:   true,
				ProfileResponses: true,
				CacheDuration:    10 * time.Minute,
				MaxCacheSize:     1000,
			},
			User: UserCacheConfig{
				Enabled:         true,
				Duration:        5 * time.Minute,
				MaxUsers:        500,
				CleanupInterval: 1 * time.Minute,
			},
		},
		Texture: TextureConfig{
			BaseURL:       "",
			UploadEnabled: false,
			MaxFileSize:   1024 * 1024, // 1MB
			AllowedTypes:  []string{"image/png", "image/jpeg"},
		},
		Yggdrasil: YggdrasilConfig{
			Meta: MetaConfig{
				ServerName:            "Yggdrasil API Server (Go)",
				ImplementationName:    "yggdrasil-api-go",
				ImplementationVersion: "1.0.0",
				Links: map[string]string{
					"homepage": "http://localhost:8080",
					"register": "http://localhost:8080/register",
				},
			},
			SkinDomains: []string{
				".minecraft.net", // Minecraft官方域名
				".mojang.com",    // Mojang官方域名
			},
			Keys: KeysConfig{
				PrivateKeyPath: "keys/private.pem", // 密钥文件路径
				PublicKeyPath:  "keys/public.pem",
			},
			Features: FeaturesConfig{
				NonEmailLogin: true,
			},
		},
		Middleware: MiddlewareConfig{
			CORS: CORSConfig{
				Enabled:          true,
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
				ExposedHeaders:   []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           86400, // 24小时
			},
			RateLimit: RateLimitConfig{
				Enabled:      true,
				AuthInterval: 1 * time.Second,
				GlobalLimit:  60, // 每分�?0个请�?				BurstLimit:   10, // 突发10个请�?			},
			Performance: PerformanceConfig{
				Enabled:         true,
				DetailedMetrics: false,
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			File:       "logs/yggdrasil.log",
			MaxSize:    100, // 100MB
			MaxBackups: 3,
			MaxAge:     7, // 7�?			Compress:   true,
		},
		Monitoring: MonitoringConfig{
			Enabled:         true,
			MetricsEndpoint: "/metrics",
			CacheStats:      true,
			DBStats:         true,
			SystemStats:     true,
		},
		Security: SecurityConfig{
			MaxRequestSize: "1MB",
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    60 * time.Second,
		},
		Warmup: WarmupConfig{
			Enabled:       true,
			APIMetadata:   true,
			ErrorCache:    true,
			UUIDCache:     true,
			ProfileCache:  false,
			ConcurrentNum: 5,
		},
	}
}
