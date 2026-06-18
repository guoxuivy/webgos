package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Database struct {
		// 主库配置
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Dialect  string `yaml:"dialect"`
		// 连接池配置（主库）
		MaxOpenConns int `yaml:"max_open_conns"`
		MaxIdleConns int `yaml:"max_idle_conns"`
		MaxLifetime  int `yaml:"max_lifetime"`
		// 读写分离策略（新增）
		ReadWriteSeparation bool   `yaml:"read_write_separation"` // 是否启用读写分离
		SlaveLoadBalance    string `yaml:"slave_load_balance"`    // 备库负载均衡策略：random, round_robin（轮询）
		// 备库配置（新增）
		Slaves []struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			DBName   string `yaml:"dbname"`
			// 备库连接池配置
			MaxOpenConns int `yaml:"max_open_conns"`
			MaxIdleConns int `yaml:"max_idle_conns"`
			MaxLifetime  int `yaml:"max_lifetime"`
		} `yaml:"slaves"`
	} `yaml:"database"`
	Server struct {
		Mode string `yaml:"mode"` // "debug" 或 "release"
		Port int    `yaml:"port"` // 服务器端口
		Swag bool   `yaml:"swag"` // 是否启用 Swagger 文档接口
	} `yaml:"server"`
	Runtime struct {
		Dir string `yaml:"dir"` // 运行时数据目录，日志、黑名单等文件均存放于此
	} `yaml:"runtime"`
	Log struct {
		Level    string `yaml:"level"`     // 日志级别:  Error, Warn, Info
		Access   bool   `yaml:"access"`    // 是否启用访问日志
		LevelSQL string `yaml:"level_sql"` // SQL日志级别: Silent, Error, Warn, Info
	} `yaml:"log"`
	JWT struct {
		Secret string `yaml:"secret"`
		Expiry int    `yaml:"expiry"` // 小时
	} `yaml:"jwt"`
	Website struct {
		Dir           string `yaml:"dir"`             // website 根目录
		UploadUrl     string `yaml:"upload_url"`      // 临时文件目录
		UploadTempUrl string `yaml:"upload_temp_url"` // 临时文件目录
	} `yaml:"website"`

	// 自动迁移配置
	AutoMigrate bool `yaml:"auto_migrate"`
	// 自动同步RBAC权限点
	AutoRBACPoint bool `yaml:"auto_rbac_point"`
	// 超级管理员账号
	SuperAccount string `yaml:"super_account"`
}

var GlobalConfig *Config

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	GlobalConfig = &Config{}
	if err := yaml.Unmarshal(yamlFile, GlobalConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(GlobalConfig); err != nil {
		return nil, err
	}

	return GlobalConfig, nil
}

// validateConfig 验证配置有效性
func validateConfig(config *Config) error {
	errs := make([]string, 0)

	if config.Database.Host == "" {
		errs = append(errs, "database host is required")
	}
	if config.Database.Port == 0 {
		errs = append(errs, "database port is required")
	}
	if config.Database.Username == "" {
		errs = append(errs, "database username is required")
	}
	if config.Database.Password == "" {
		errs = append(errs, "database password is required")
	}
	if config.Database.DBName == "" {
		errs = append(errs, "database dbname is required")
	}
	if config.Server.Port == 0 {
		errs = append(errs, "server port is required")
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid config: %s", strings.Join(errs, ", "))
	}

	// 设置默认值
	if config.Server.Mode == "" {
		config.Server.Mode = "debug" // 设置默认值
	}

	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 20
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 5
	}
	if config.Database.MaxLifetime == 0 {
		config.Database.MaxLifetime = 60
	}

	if config.Runtime.Dir == "" {
		config.Runtime.Dir = "./runtime"
	}
	if config.Log.Level == "" {
		config.Log.Level = "Info"
	}
	if config.Log.LevelSQL == "" {
		config.Log.LevelSQL = "Info"
	}
	if config.JWT.Secret == "" {
		config.JWT.Secret = "sean_secret_key"
	}

	// 设置上传目录默认值
	if config.Website.Dir == "" {
		config.Website.Dir = "./public"
	}
	if config.Website.UploadUrl == "" {
		config.Website.UploadUrl = "/upload"
	}
	if config.Website.UploadTempUrl == "" {
		config.Website.UploadTempUrl = "/upload/temp"
	}

	return nil
}
