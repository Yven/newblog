package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"newblog/internal/config"
	"newblog/internal/global"
	"newblog/internal/service"
	"newblog/internal/util"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth(srv service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		claims, err := srv.BearerHeaderCheck(bearer)
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

		var bodyBytes []byte
		// 获取请求体数据
		if c.Request.Body != nil {
			bodyBytes, _ = c.GetRawData()
			// 重新设置请求体，因为GetRawData会清空body
			// 使用io.NopCloser重新构造请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		var body any
		if len(bodyBytes) > 0 {
			json.Unmarshal(bodyBytes, &body)
		}

		c.Next()

		end := time.Now()
		latency := strings.Trim(fmt.Sprintf("%13v", end.Sub(start)), " ")

		global.Logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.String("latency", latency),
		)

		// gin 的 logger 中间件请求标准输出
		if config.Global.App.Env != "release" {
			param := gin.LogFormatterParams{}
			statusColor := param.StatusCodeColor()
			methodColor := param.MethodColor()
			resetColor := param.ResetColor()

			fmt.Fprintf(os.Stdout, "[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, c.Writer.Status(), resetColor,
				latency,
				c.ClientIP(),
				methodColor, c.Request.Method, resetColor,
				path,
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		}

		for _, e := range c.Errors {
			global.Logger.Error("接口错误",
				slog.Any("error", e.Err),
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.Any("token", c.Request.Header.Get("Authorization")),
				slog.Any("query", c.Request.URL.RawQuery),
				slog.Any("body", body),
				slog.String("ip", c.ClientIP()),
				slog.Int("status", c.Writer.Status()),
				slog.String("latency", latency),
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
