package routes

import (
	"app_api/controllers"
	"app_api/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	// 不需要鉴权的公共路由
	userPublic := r.Group("/users")
	{
		userPublic.POST("/login", controllers.LoginUser)     // 登录
		userPublic.POST("/register", controllers.CreateUser) // 注册
	}
	user := r.Group("/users").Use(middleware.JWTAuth())
	{
		user.POST("/info", controllers.GetUserInfo) // 获取用户信息
	}
	news := r.Group("/news")
	{
		// 支持GET/POST两种方式获取新闻列表
		news.GET("/list", controllers.GetNewsList)
		news.POST("/list", controllers.GetNewsList)
	}
}
