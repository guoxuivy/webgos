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
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Server struct {
		Mode   string `yaml:"mode"` // "debug" 或 "release"
		LogDir string `yaml:"logdir"`
		Port   int    `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret string `yaml:"secret"`
		Expiry int    `yaml:"expiry"` // 小时
	} `yaml:"jwt"`
	// 自动迁移配置
	AutoMigrate bool `yaml:"auto_migrate"`
	// 自动同步RBAC权限点
	AutoRBACPoint bool `yaml:"auto_rbac_point"`
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

	return nil
}
