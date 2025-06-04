package services

import (
	"app_api/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// initLogger 初始化日志记录器
func initLogger() {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic(fmt.Sprintf("初始化日志失败: %v", err))
		}
	})
}

// StartNewsScheduler 启动定时任务，立即执行一次，之后每180分钟执行一次
func StartNewsScheduler(first bool) {
	initLogger()

	if first {
		// 立即执行一次
		if err := RequestNewsAPIServices(); err != nil {
			logger.Error("首次执行新闻API服务失败", zap.Error(err))
		}
	}

	c := cron.New()
	_, err := c.AddFunc("@every 180m", func() {
		if err := RequestNewsAPIServices(); err != nil {
			logger.Error("定时执行新闻API服务失败", zap.Error(err))
		}
	})

	if err != nil {
		logger.Panic("创建定时任务失败", zap.Error(err))
	}

	c.Start()
}

// RequestNewsAPIServices 从API获取新闻数据并保存到数据库
func RequestNewsAPIServices() error {
	initLogger()

	newsTypes := map[string]string{
		"":       "推荐",
		"guonei": "国内",
		"tiyu":   "体育",
		"keji":   "科技",
		"youxi":  "游戏",
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(newsTypes))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for newsType, typeName := range newsTypes {
		wg.Add(1)
		go func(newsType, typeName string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- fmt.Errorf("处理%s类新闻超时", typeName)
				return
			default:
				if err := fetchAndSaveNews(client, newsType, typeName); err != nil {
					errChan <- err
				}
			}
		}(newsType, typeName)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// fetchAndSaveNews 获取并保存单个类型的新闻数据
func fetchAndSaveNews(client *http.Client, newsType, typeName string) error {
	apiUrl := "http://v.juhe.cn/toutiao/index"
	apiKey := "7292bbea82b0ad89f513adaa8d4a5d93"

	params := url.Values{}
	params.Set("key", apiKey)
	params.Set("type", newsType)
	params.Set("page", "1")
	params.Set("page_size", "30")
	params.Set("is_filter", "0")

	resp, err := client.Get(apiUrl + "?" + params.Encode())
	if err != nil {
		logger.Error("获取新闻API响应失败",
			zap.String("type", typeName),
			zap.Error(err))
		return fmt.Errorf("获取%s类新闻失败: %w", typeName, err)
	}
	defer resp.Body.Close()

	var result struct {
		Reason string `json:"reason"`
		Result struct {
			Data []models.News `json:"data"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("解析新闻响应失败",
			zap.String("type", typeName),
			zap.Error(err))
		return fmt.Errorf("解析%s类新闻响应失败: %w", typeName, err)
	}

	if len(result.Result.Data) > 0 {
		if err := models.BatchCreateNews(result.Result.Data); err != nil {
			logger.Error("保存新闻数据失败",
				zap.String("type", typeName),
				zap.Error(err))
			return fmt.Errorf("保存%s类新闻失败: %w", typeName, err)
		}
		logger.Info("成功保存新闻数据",
			zap.String("type", typeName),
			zap.Int("count", len(result.Result.Data)))
	} else {
		logger.Info("未找到新闻数据",
			zap.String("type", typeName))
	}

	return nil
}
