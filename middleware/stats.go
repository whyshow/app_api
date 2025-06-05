package middleware

import (
	"app_api/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// ClientInfo 客户端信息结构体
type ClientInfo struct {
	Type    string `json:"client_type"`
	OS      string `json:"os"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

func StatsMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 获取并校验端点路径
		endpoint := c.FullPath()
		if endpoint == "" {
			return
		}

		// 获取客户端信息
		info := parseClientInfo(c)

		// 生成JSON字符串
		infoJSON, _ := json.Marshal(info)

		date := time.Now().UTC().Truncate(24 * time.Hour)
		ip := c.ClientIP()

		// 原子化事务处理
		db.Transaction(func(tx *gorm.DB) error {
			// 更新IP详细记录（新增client_type主键）
			tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "date"}, {Name: "ip"}, {Name: "endpoint"}, {Name: "client_type"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"request_count": gorm.Expr("ip_details.request_count + 1")}),
			}).Create(&models.IPDetail{
				Date:         date,
				IP:           ip,
				Endpoint:     endpoint,
				ClientType:   string(infoJSON),
				RequestCount: 1,
			})

			// 更新接口统计（新增client_type主键）
			tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "date"}, {Name: "endpoint"}, {Name: "client_type"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"request_count": gorm.Expr("endpoint_stats.request_count + 1")}),
			}).Create(&models.EndpointStat{
				Date:         date,
				Endpoint:     endpoint,
				ClientType:   string(infoJSON),
				RequestCount: 1,
			})

			// 更新总统计（按客户端类型分区）
			tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "date"}, {Name: "client_type"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"total_requests": gorm.Expr("daily_request_stats.total_requests + 1"),
					"unique_ips":     gorm.Expr("(SELECT COUNT(DISTINCT ip) FROM ip_details WHERE date = ? AND client_type = ?)", date, string(infoJSON)),
				}),
			}).Create(&models.DailyRequestStat{
				Date:          date,
				ClientType:    string(infoJSON),
				TotalRequests: 1,
			})

			return nil
		})
	}
}

func parseClientInfo(c *gin.Context) ClientInfo {
	// 初始化解析器（实际使用时应复用实例）
	parser := uaparser.NewFromSaved()

	ua := c.GetHeader("User-Agent")
	client := parser.Parse(ua)

	return ClientInfo{
		Type: firstNonEmpty(
			c.GetHeader("X-Client-Type"),
			client.Device.Family,
			"unknown",
		),
		Name: firstNonEmpty(client.UserAgent.Family, "unknown"),
		OS:   firstNonEmpty(client.Os.Family, "unknown"),
		Version: firstNonEmpty(
			client.UserAgent.Major+"."+client.UserAgent.Minor,
			"unknown",
		),
	}
}

// 辅助函数取第一个非空值
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
