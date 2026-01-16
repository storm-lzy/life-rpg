// Package models 数据模型定义
package models

import (
	"time"

	"gorm.io/gorm"
)

// SysUser 系统用户
type SysUser struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	RoleID    uint           `gorm:"index" json:"roleId"`
	Role      *SysRole       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Gold      int            `gorm:"default:0" json:"gold"`
	Exp       int            `gorm:"default:0" json:"exp"`
	Level     int            `gorm:"default:1" json:"level"`
	Status    int            `gorm:"default:1" json:"status"` // 1正常 0禁用
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SysUser) TableName() string {
	return "sys_user"
}

// SysRole 系统角色
type SysRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Key       string         `gorm:"size:50;uniqueIndex;not null" json:"key"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int            `gorm:"default:1" json:"status"`
	Remark    string         `gorm:"size:255" json:"remark"`
	Menus     []SysMenu      `gorm:"many2many:role_menu" json:"menus,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SysRole) TableName() string {
	return "sys_role"
}

// SysMenu 系统菜单
type SysMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"default:0;index" json:"parentId"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Path      string         `gorm:"size:255" json:"path"`
	Component string         `gorm:"size:255" json:"component"`
	Icon      string         `gorm:"size:50" json:"icon"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Type      int            `gorm:"default:1" json:"type"` // 1目录 2菜单 3按钮
	Permission string        `gorm:"size:100" json:"permission"`
	Visible   int            `gorm:"default:1" json:"visible"`
	Status    int            `gorm:"default:1" json:"status"`
	Children  []SysMenu      `gorm:"-" json:"children,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SysMenu) TableName() string {
	return "sys_menu"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}

// TableName 表名
func (RoleMenu) TableName() string {
	return "role_menu"
}
