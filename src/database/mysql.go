package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/spf13/viper"
)

// MySQLConfig MySQL数据库配�?type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Charset         string        `mapstructure:"charset"`
	Collation       string        `mapstructure:"collation"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	LogLevel        string        `mapstructure:"log_level"`
}

// MySQLManager MySQL数据库管理器
type MySQLManager struct {
	DB *gorm.DB
}

// NewMySQLManager 创建MySQL管理�?func NewMySQLManager(config *MySQLConfig) (*MySQLManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.Collation,
	)

	// 配置GORM日志级别
	var logLevel logger.LogLevel
	switch config.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// GORM配置
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		PrepareStmt: true,                              // 预编译语�?		CreateBatchSize: 100,                           // 批量创建大小
		QueryFields: true,                              // 查询所有字�?		DisableForeignKeyConstraintWhenMigrating: true, // 迁移时禁用外键约�?	}

	// 连接数据�?	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 获取底层SQL数据库连�?	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// 设置连接池参�?	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	log.Printf("Successfully connected to MySQL database: %s@%s:%d/%s",
		config.Username, config.Host, config.Port, config.Database)

	return &MySQLManager{DB: db}, nil
}

// AutoMigrate 自动迁移数据库表结构
func (m *MySQLManager) AutoMigrate(models ...interface{}) error {
	return m.DB.AutoMigrate(models...)
}

// Close 关闭数据库连�?func (m *MySQLManager) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction 执行事务
func (m *MySQLManager) Transaction(fc func(tx *gorm.DB) error) error {
	return m.DB.Transaction(fc)
}

// GetDB 获取数据库实�?func (m *MySQLManager) GetDB() *gorm.DB {
	return m.DB
}

// HealthCheck 数据库健康检�?func (m *MySQLManager) HealthCheck() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetStats 获取数据库统计信�?func (m *MySQLManager) GetStats() map[string]interface{} {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections": stats.OpenConnections,
		"in_use": stats.InUse,
		"idle": stats.Idle,
		"wait_count": stats.WaitCount,
		"wait_duration": stats.WaitDuration.String(),
		"max_idle_closed": stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}

// LoadMySQLConfig 从配置文件加载MySQL配置
func LoadMySQLConfig(viper *viper.Viper) (*MySQLConfig, error) {
	config := &MySQLConfig{}
	
	// 设置默认�?	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.collation", "utf8mb4_unicode_ci")
	viper.SetDefault("database.mysql.max_open_conns", 25)
	viper.SetDefault("database.mysql.max_idle_conns", 5)
	viper.SetDefault("database.mysql.conn_max_lifetime", "300s")
	viper.SetDefault("database.mysql.log_level", "info")

	// 读取配置
	if err := viper.UnmarshalKey("database.mysql", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MySQL config: %w", err)
	}

	// 验证必要配置
	if config.Database == "" {
		return nil, fmt.Errorf("database name is required")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("database username is required")
	}

	// 解析时间配置
	if connMaxLifetime := viper.GetString("database.mysql.conn_max_lifetime"); connMaxLifetime != "" {
		duration, err := time.ParseDuration(connMaxLifetime)
		if err != nil {
			return nil, fmt.Errorf("invalid conn_max_lifetime: %w", err)
		}
		config.ConnMaxLifetime = duration
	}

	return config, nil
}

// MySQLModel 基础模型
type MySQLModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// PaginatedResult 分页查询结果
type PaginatedResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Sort     string `form:"sort"`
	Order    string `form:"order,default=desc" binding:"oneof=asc desc"`
}

// ApplyPagination 应用分页
func ApplyPagination(query *gorm.DB, params *PaginationParams) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)
	
	if params.Sort != "" {
		order := params.Sort
		if params.Order == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}
	
	return query
}

// Paginate 执行分页查询
func Paginate(query *gorm.DB, params *PaginationParams, result interface{}) (*PaginatedResult, error) {
	var total int64
	
	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 应用分页并查询数�?	query = ApplyPagination(query, params)
	if err := query.Find(result).Error; err != nil {
		return nil, err
	}
	
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}
	
	return &PaginatedResult{
		Data:       result,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}
