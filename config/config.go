package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var C Config

type Config struct {
	Server    ServerConfig   `yaml:"server"`
	Database  DatabaseConfig `yaml:"database"`
	JWTSecret string         `yaml:"jwt_secret"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// Load 加载配置
func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 初始加载配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 配置热更新
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("检测到配置变更:", e.Name)
		if err := viper.Unmarshal(&C); err != nil {
			log.Printf("配置热更新失败: %v", err)
		}
	})

	// 绑定配置到结构体
	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}
}

// Get 获取配置实例
func Get() *Config {
	return &C
}
