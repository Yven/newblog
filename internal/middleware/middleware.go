package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"newblog/internal/global"
	"newblog/internal/util"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := global.JwtService.BearerHeaderCheck(c.GetHeader("Authorization"))
		if err != nil {
			util.ErrorAbort(c, http.StatusUnauthorized, err)
			return
		}

		// 从 token 中获取用户信息并保存到 context 中
		c.Set("userId", claims.Subject)

		c.Next()
	}
}

func SlogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)

		global.Logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
		)

		for _, e := range c.Errors {
			global.Logger.Error("接口错误",
				slog.Any("error", e.Err),
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.String("ip", c.ClientIP()),
				slog.Int("status", c.Writer.Status()),
				slog.Duration("latency", latency),
			)
		}
	}
}

func SlogRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				global.Logger.Error("接口异常终止",
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.String("error", err.Error()),
					slog.String("stack", "...\n"+string(buf)),
				)
				// c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				// 	"error": "内部错误，请稍后再试",
				// })
			}
		}()
		c.Next()
	}
}
