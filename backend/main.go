// Life RPG Backend - 人生游戏化系统后端
// 技术栈: Go + Gin + GORM + MySQL + JWT

package main

import (
	"log"

	"life-rpg/config"
	"life-rpg/database"
	"life-rpg/middleware"
	"life-rpg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()
	log.Println("配置加载完成")

	// 初始化数据库
	database.InitDB()

	// 初始化种子数据
	database.SeedData()

	// 创建 Gin 引擎
	r := gin.Default()

	// 应用 CORS 中间件
	r.Use(middleware.CORS())

	// 配置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
