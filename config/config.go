package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var C Config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Key      KeyConfig      `yaml:"key"`
	Page     PageConfig     `yaml:"page"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type PageConfig struct {
	MinPageSize int `yaml:"min_page_size"  mapstructure:"min_page_size"`
	MaxPageSize int `yaml:"max_page_size"  mapstructure:"max_page_size"`
}

type KeyConfig struct {
	NewsApikey string `yaml:"news_apikey" mapstructure:"news_apikey"`
	JwtSecret  string `yaml:"jwt_secret"  mapstructure:"jwt_secret"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	DbName       string `yaml:"dbname"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
}

// DatabaseInitializer 定义数据库初始化接口
type DatabaseInitializer interface {
	Init() error
}

var dbInitializer DatabaseInitializer

// SetDatabaseInitializer 设置数据库初始化器
func SetDatabaseInitializer(initializer DatabaseInitializer) {
	dbInitializer = initializer
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
		} else {
			fmt.Printf("配置热更新成功: %+v\n", C)

			// 通过接口重新初始化数据库
			if dbInitializer != nil {
				if err := dbInitializer.Init(); err != nil {
					log.Printf("数据库重新初始化失败: %v", err)
				}
			}
		}
	})

	// 绑定配置到结构体
	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	// 打印配置
	fmt.Printf("配置加载完成: %+v\n", C)
}

// Get 获取配置实例
func Get() *Config {
	return &C
}
