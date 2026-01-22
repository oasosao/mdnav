package middleware

import (
	"mdnav/internal/handler/base"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 增强版日志中间件
func ZapLoggerWithConfig(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 将logger存入context，供后续处理函数使用
		c.Set("zapLogger", logger)

		// 请求处理前可以记录请求体（注意性能）
		// body, _ := io.ReadAll(c.Request.Body)
		// c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next() // 处理请求

		latency := time.Since(start)

		// 获取错误信息（如果有）
		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		// 创建日志字段
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("latency_human", latency.String()),
			zap.String("time", time.Now().Format("2006-01-02 15:04:05")),
			zap.String("err_msg", errMsg),
		}

		// 根据状态码选择日志级别
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Error("服务器错误", fields...)
		case status >= 400:
			logger.Warn("客户端错误", fields...)
		}

		if status >= 400 {
			base.Error(c, logger, status)
		}
	}

}
