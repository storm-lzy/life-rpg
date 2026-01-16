// Package models 游戏业务模型
package models

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"size:500" json:"description"`
	GoldReward  int            `gorm:"default:0" json:"goldReward"`
	ExpReward   int            `gorm:"default:0" json:"expReward"`
	Type        string         `gorm:"size:20;default:daily" json:"type"` // daily每日 once一次性
	Category    string         `gorm:"size:50" json:"category"`
	Icon        string         `gorm:"size:50" json:"icon"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	Sort        int            `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Task) TableName() string {
	return "task"
}

// UserTask 用户任务完成记录
type UserTask struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	TaskID      uint      `gorm:"index;not null" json:"taskId"`
	Task        *Task     `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	CompletedAt time.Time `json:"completedAt"`
}

// TableName 表名
func (UserTask) TableName() string {
	return "user_task"
}

// Reward 奖励/商品
type Reward struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"size:500" json:"description"`
	Cost        int            `gorm:"default:0" json:"cost"`
	Stock       int            `gorm:"default:-1" json:"stock"` // -1无限
	Image       string         `gorm:"size:255" json:"image"`
	Category    string         `gorm:"size:50" json:"category"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	Sort        int            `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Reward) TableName() string {
	return "reward"
}

// UserLog 用户流水记录
type UserLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	Type        string    `gorm:"size:20;not null" json:"type"` // gold_in/gold_out/exp_in
	Amount      int       `gorm:"not null" json:"amount"`
	Balance     int       `json:"balance"`
	Description string    `gorm:"size:255" json:"description"`
	RefType     string    `gorm:"size:50" json:"refType"` // task/reward/admin
	RefID       uint      `json:"refId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TableName 表名
func (UserLog) TableName() string {
	return "user_log"
}

// Announcement 公告
type Announcement struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:20;default:notice" json:"type"` // notice/activity/update
	IsActive  bool           `gorm:"default:true" json:"isActive"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Announcement) TableName() string {
	return "announcement"
}

// ThemeConfig H5主题配置
type ThemeConfig struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PrimaryColor   string    `gorm:"size:20;default:#1989fa" json:"primaryColor"`
	SecondaryColor string    `gorm:"size:20;default:#ff976a" json:"secondaryColor"`
	GoldColor      string    `gorm:"size:20;default:#ffd700" json:"goldColor"`
	ExpColor       string    `gorm:"size:20;default:#07c160" json:"expColor"`
	BackgroundURL  string    `gorm:"size:255" json:"backgroundUrl"`
	LogoURL        string    `gorm:"size:255" json:"logoUrl"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TableName 表名
func (ThemeConfig) TableName() string {
	return "theme_config"
}
