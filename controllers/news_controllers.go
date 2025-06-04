package controllers

import (
	"app_api/models"
	"app_api/utils"
	"github.com/gin-gonic/gin"
)

// GetNewsList 获取新闻列表接口
func GetNewsList(c *gin.Context) {
	var req models.NewsListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}

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

	newsList, err := models.GetNewsListData(req)

	if err != nil {
		utils.Error(c, "查询数据库失败")
		return
	}

	utils.Success(c, "成功", gin.H{
		"result": newsList,
	})
}
