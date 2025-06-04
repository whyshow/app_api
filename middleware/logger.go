package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var logger *zap.Logger

func init() {
	// 初始化生产环境日志器
	if gin.Mode() == gin.DebugMode {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
}

// Logger 中间件，记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取请求体
		var bodyBuffer bytes.Buffer
		body := c.Request.Body
		defer body.Close()
		if _, err := bodyBuffer.ReadFrom(body); err == nil {
			c.Request.Body = &readClose{bytes.NewReader(bodyBuffer.Bytes())}
		}

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("trace-id", c.GetString("X-Trace-ID")),
		}

		if len(bodyBuffer.Bytes()) > 0 {
			fields = append(fields, zap.ByteString("body", bodyBuffer.Bytes()))
		}

		if latency > time.Second {
			logger.Warn("Slow request", fields...)
		} else {
			logger.Info("Request handled", fields...)
		}
	}
}

// 实现可关闭的Reader
type readClose struct{ *bytes.Reader }

func (r *readClose) Close() error { return nil }
