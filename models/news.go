package models

import (
	"app_api/config"
	"app_api/database"
)

// News 对应数据库表结构
type News struct {
	Uniquekey       string `gorm:"primaryKey;size:255"`
	Title           string `gorm:"size:255"`
	Date            string `gorm:"size:50"`
	Category        string `gorm:"size:100"`
	AuthorName      string `gorm:"size:100;column:author_name" json:"author_name"`
	Url             string `gorm:"size:255"`
	ThumbnailPicS   string `gorm:"size:255;column:thumbnail_pic_s" json:"thumbnail_pic_s"`
	ThumbnailPicS02 string `gorm:"size:255;column:thumbnail_pic_s02" json:"thumbnail_pic_s02"`
	ThumbnailPicS03 string `gorm:"size:255;column:thumbnail_pic_s03" json:"thumbnail_pic_s03"`
	IsContent       string `gorm:"size:10;column:is_content" json:"is_content"`
}

type NewsListRequest struct {
	Type     string `form:"type" json:"type"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// GetNewsListData 获取新闻列表数据
func GetNewsListData(data NewsListRequest) ([]News, error) {
	var newsList []News

	if data.PageSize > config.Get().Page.MaxPageSize {
		data.PageSize = config.Get().Page.MaxPageSize
	} else if data.PageSize < config.Get().Page.MinPageSize {
		data.PageSize = config.Get().Page.MinPageSize
	}

	query := database.GetDB().Order("date DESC").Offset((data.Page - 1) * data.PageSize).Limit(data.PageSize)
	if data.Type != "" {
		query = query.Where("category = ?", data.Type)
	}

	if err := query.Find(&newsList).Error; err != nil {
		return nil, err
	}

	return newsList, nil
}

// BatchCreateNews 批量保存新闻数据
func BatchCreateNews(newsList []News) error {
	if len(newsList) == 0 {
		return nil
	}
	return database.GetDB().Create(&newsList).Error
}

// CreateOrUpdateNews 创建或更新单条新闻数据
func CreateOrUpdateNews(news News) error {
	return database.GetDB().Save(&news).Error
}

// GetNewsByUniqueKey 根据唯一键获取新闻
func GetNewsByUniqueKey(uniquekey string) (News, error) {
	var news News
	err := database.GetDB().Where("uniquekey = ?", uniquekey).First(&news).Error
	return news, err
}
