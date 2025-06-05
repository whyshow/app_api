package main

import (
	"app_api/config"
	"app_api/database"
	"app_api/middleware"
	"app_api/models"
	"app_api/routes"
	"app_api/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	// 初始化数据库
	if err := database.Init(); err != nil {
		panic(err)
	}
	// 执行数据库迁移（在此处传入需要迁移的模型）
	if err := database.UpgradeDatabase([]interface{}{&models.User{}, &models.News{}}); err != nil {
		panic(err)
	}

	// 启动定时请求新闻API服务器
	services.StartNewsScheduler(false)

	// 创建Gin实例
	router := gin.New()
	// 使用全局中间件（按顺序执行）
	router.Use(
		middleware.Logger(), // 自定义日志中间件
	)

	routes.Routes(router)

	// 启动服务
	router.Run(":" + config.Get().Server.Port)
}
