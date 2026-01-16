// Package database 数据库初始化种子数据
package database

import (
	"log"

	"life-rpg/models"

	"golang.org/x/crypto/bcrypt"
)

// SeedData 初始化种子数据
func SeedData() {
	// 检查是否已有角色数据
	var roleCount int64
	DB.Model(&models.SysRole{}).Count(&roleCount)
	if roleCount > 0 {
		log.Println("种子数据已存在，跳过初始化")
		return
	}

	log.Println("开始初始化种子数据...")

	// 创建角色
	roles := []models.SysRole{
		{Name: "超级管理员", Key: "admin", Sort: 1, Remark: "系统超级管理员"},
		{Name: "普通用户", Key: "user", Sort: 2, Remark: "普通游戏用户"},
	}
	DB.Create(&roles)
	log.Println("角色数据创建完成")

	// 创建菜单
	menus := []models.SysMenu{
		// 仪表盘
		{ID: 1, ParentID: 0, Name: "仪表盘", Path: "/admin/dashboard", Component: "admin/Dashboard", Icon: "Odometer", Sort: 1, Type: 2},
		// 系统管理
		{ID: 2, ParentID: 0, Name: "系统管理", Path: "/admin/system", Component: "", Icon: "Setting", Sort: 2, Type: 1},
		{ID: 21, ParentID: 2, Name: "用户管理", Path: "/admin/system/user", Component: "admin/system/User", Icon: "User", Sort: 1, Type: 2, Permission: "system:user:list"},
		{ID: 22, ParentID: 2, Name: "角色管理", Path: "/admin/system/role", Component: "admin/system/Role", Icon: "UserFilled", Sort: 2, Type: 2, Permission: "system:role:list"},
		{ID: 23, ParentID: 2, Name: "菜单管理", Path: "/admin/system/menu", Component: "admin/system/Menu", Icon: "Menu", Sort: 3, Type: 2, Permission: "system:menu:list"},
		// 游戏配置
		{ID: 3, ParentID: 0, Name: "游戏配置", Path: "/admin/game", Component: "", Icon: "TrophyBase", Sort: 3, Type: 1},
		{ID: 31, ParentID: 3, Name: "任务管理", Path: "/admin/game/task", Component: "admin/game/Task", Icon: "List", Sort: 1, Type: 2, Permission: "game:task:list"},
		{ID: 32, ParentID: 3, Name: "奖励管理", Path: "/admin/game/reward", Component: "admin/game/Reward", Icon: "Present", Sort: 2, Type: 2, Permission: "game:reward:list"},
		// 公告管理
		{ID: 4, ParentID: 0, Name: "公告管理", Path: "/admin/announcement", Component: "admin/Announcement", Icon: "Bell", Sort: 4, Type: 2, Permission: "announcement:list"},
		// 主题配置
		{ID: 5, ParentID: 0, Name: "H5主题配置", Path: "/admin/theme", Component: "admin/Theme", Icon: "Brush", Sort: 5, Type: 2, Permission: "theme:config"},
	}
	DB.Create(&menus)
	log.Println("菜单数据创建完成")

	// 为管理员角色分配所有菜单
	var allMenus []models.SysMenu
	DB.Find(&allMenus)
	for _, menu := range allMenus {
		DB.Create(&models.RoleMenu{RoleID: 1, MenuID: menu.ID})
	}
	log.Println("角色菜单关联创建完成")

	// 创建超级管理员账号
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	admin := models.SysUser{
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "超级管理员",
		RoleID:   1,
		Gold:     0,
		Exp:      0,
		Level:    1,
		Status:   1,
	}
	DB.Create(&admin)
	log.Println("管理员账号创建完成 (admin/123456)")

	// 创建示例任务
	tasks := []models.Task{
		{Title: "早起打卡", Description: "早上7点前起床", GoldReward: 10, ExpReward: 5, Type: "daily", Category: "健康", Icon: "🌅", IsActive: true, Sort: 1},
		{Title: "阅读30分钟", Description: "阅读书籍或文章30分钟", GoldReward: 15, ExpReward: 10, Type: "daily", Category: "学习", Icon: "📚", IsActive: true, Sort: 2},
		{Title: "运动锻炼", Description: "完成30分钟运动", GoldReward: 20, ExpReward: 15, Type: "daily", Category: "健康", Icon: "🏃", IsActive: true, Sort: 3},
		{Title: "喝8杯水", Description: "今日饮水达标", GoldReward: 5, ExpReward: 3, Type: "daily", Category: "健康", Icon: "💧", IsActive: true, Sort: 4},
		{Title: "完成周报", Description: "提交本周工作总结", GoldReward: 50, ExpReward: 30, Type: "once", Category: "工作", Icon: "📝", IsActive: true, Sort: 5},
	}
	DB.Create(&tasks)
	log.Println("示例任务创建完成")

	// 创建示例奖励
	rewards := []models.Reward{
		{Title: "休息15分钟", Description: "给自己一个短暂的休息", Cost: 20, Stock: -1, Category: "休闲", IsActive: true, Sort: 1},
		{Title: "看一集电视剧", Description: "追一集喜欢的剧", Cost: 50, Stock: -1, Category: "休闲", IsActive: true, Sort: 2},
		{Title: "点一杯奶茶", Description: "奖励自己一杯奶茶", Cost: 100, Stock: -1, Category: "美食", IsActive: true, Sort: 3},
		{Title: "游戏时间1小时", Description: "畅玩游戏1小时", Cost: 80, Stock: -1, Category: "娱乐", IsActive: true, Sort: 4},
		{Title: "周末外出旅行", Description: "来一场说走就走的旅行", Cost: 500, Stock: 1, Category: "旅行", IsActive: true, Sort: 5},
	}
	DB.Create(&rewards)
	log.Println("示例奖励创建完成")

	// 创建示例公告
	announcements := []models.Announcement{
		{Title: "欢迎使用Life RPG", Content: "将生活游戏化，每天完成任务获取金币和经验，兑换心仪的奖励！", Type: "notice", IsActive: true, Sort: 1},
		{Title: "新功能上线", Content: "任务大厅已支持每日任务自动刷新，快来体验吧！", Type: "update", IsActive: true, Sort: 2},
	}
	DB.Create(&announcements)
	log.Println("示例公告创建完成")

	// 创建默认主题配置
	theme := models.ThemeConfig{
		PrimaryColor:   "#1989fa",
		SecondaryColor: "#ff976a",
		GoldColor:      "#ffd700",
		ExpColor:       "#07c160",
	}
	DB.Create(&theme)
	log.Println("默认主题配置创建完成")

	log.Println("种子数据初始化完成!")
}
