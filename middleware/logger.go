package middleware

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

var (
	csvWriter *csv.Writer
	ljLogger  *lumberjack.Logger
	mu        sync.Mutex
)

func init() {
	ljLogger = &lumberjack.Logger{
		Filename:   "logs/requests.csv",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	csvWriter = csv.NewWriter(ljLogger)
	writeCSVHeader()
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取原始请求体
		var bodyBuffer bytes.Buffer
		body := c.Request.Body
		defer body.Close()
		if _, err := bodyBuffer.ReadFrom(body); err == nil {
			c.Request.Body = &readClose{bytes.NewReader(bodyBuffer.Bytes())}
		}

		c.Next()
		latency := time.Since(start)

		// 清洗请求体内容
		var bodyStr string
		if len(bodyBuffer.Bytes()) > 0 {
			bodyStr = string(bodyBuffer.Bytes())
			bodyStr = strings.NewReplacer(
				"\r", "",
				"\n", "",
				`\`, "",
				`"`, "'",
				`\x5c`, "",
			).Replace(bodyStr)
		}

		// 构建CSV记录
		record := []string{
			start.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			strconv.Itoa(c.Writer.Status()),
			c.ClientIP(),
			c.Request.UserAgent(),
			latency.String(),
			c.GetString("X-Trace-ID"),
			bodyStr,
		}

		// 安全写入CSV
		mu.Lock()
		defer mu.Unlock()
		if err := csvWriter.Write(record); err != nil {
			zap.L().Error("CSV写入失败", zap.Error(err))
		}
		csvWriter.Flush()
	}
}

func writeCSVHeader() {
	header := []string{
		"timestamp", "method", "path", "status",
		"client_ip", "user_agent", "latency", "trace_id", "body",
	}
	mu.Lock()
	defer mu.Unlock()
	csvWriter.Write(header)
	csvWriter.Flush()
}

type readClose struct{ *bytes.Reader }

func (r *readClose) Close() error { return nil }
