// Package database 数据库连接管理
package database

import (
	"fmt"
	"log"

	"life-rpg/config"
	"life-rpg/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.AppConfig.DB

	// 首先连接MySQL（不指定数据库），创建数据库
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接MySQL: %v", err)
	}

	// 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	if err := rootDB.Exec(createDBSQL).Error; err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}
	log.Printf("数据库 %s 已就绪", cfg.DBName)

	// 关闭root连接
	sqlDB, _ := rootDB.DB()
	sqlDB.Close()

	// 连接到目标数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移
	autoMigrate()
}

// autoMigrate 自动迁移数据库结构
func autoMigrate() {
	err := DB.AutoMigrate(
		&models.SysRole{},
		&models.SysMenu{},
		&models.RoleMenu{},
		&models.SysUser{},
		&models.Task{},
		&models.UserTask{},
		&models.Reward{},
		&models.UserLog{},
		&models.Announcement{},
		&models.ThemeConfig{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")
}
