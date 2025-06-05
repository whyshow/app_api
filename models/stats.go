package models

import (
	"time"
)

// DailyRequestStat 每日总统计
type DailyRequestStat struct {
	Date          time.Time `gorm:"primaryKey;type:date"` // 主键，日期
	TotalRequests uint      `gorm:"default:0"`            // 用于记录每日的总请求数量
	UniqueIPs     uint      `gorm:"default:0"`            // 用于记录每日的唯一IP数量
	ClientType    string    `gorm:"size:255;default:''"`  // 客户端类型
}

// EndpointStat 接口统计
type EndpointStat struct {
	Date         time.Time `gorm:"primaryKey;type:date"` // 主键，日期
	Endpoint     string    `gorm:"primaryKey;size:255"`  // 主键，接口路径
	RequestCount uint      `gorm:"default:0"`            // 用于记录该接口的总请求数量
	ClientType   string    `gorm:"size:255;default:''"`  // 客户端类型
}

// IPDetail IP详细记录
type IPDetail struct {
	Date         time.Time `gorm:"primaryKey;type:date"` // 主键，日期
	IP           string    `gorm:"primaryKey;size:45"`   // 主键，IP地址
	Endpoint     string    `gorm:"primaryKey;size:255"`  // 主键，接口路径
	RequestCount uint      `gorm:"default:0"`            // 用于记录该IP在该接口上的请求数量
	ClientType   string    `gorm:"size:255;default:''"`  // 客户端类型
}
