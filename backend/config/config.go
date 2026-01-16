// Package config 配置管理
package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime int64 // 小时
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
}

// AppConfig 全局配置实例
var AppConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	AppConfig = &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "mysql01"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "life_rpg"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "life-rpg-secret-key-2024"),
			ExpireTime: 24, // 24小时
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
