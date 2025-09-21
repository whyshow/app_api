package controllers

import (
	"app_api/models"
	"app_api/utils"

	"github.com/gin-gonic/gin"
)

// GetNewsList 获取新闻列表接口
func GetNewsList(c *gin.Context) {

	var req models.NewsListRequest

	// 绑定请求参数
	if err := c.ShouldBind(&req); err != nil {
		utils.ParamError(c)
		return
	}

	// 类型映射
	switch req.Type {
	case "guonei":
		req.Type = "国内"
	case "tiyu":
		req.Type = "体育"
	case "keji":
		req.Type = "科技"
	case "youxi":
		req.Type = "游戏"
	default:
		req.Type = "头条"
	}

	// 调用模型获取数据
	newsList, err := models.GetNewsListData(req)

	// 处理错误
	if err != nil {
		utils.Error(c, "查询数据库失败")
		return
	}

	// 返回成功响应
	utils.Success(c, "成功", gin.H{
		"result": newsList,
		"total":  len(newsList),
	})
}
