package middleware

import (
	"net/http"
	"newblog/internal/global"
	"newblog/internal/util"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := global.JwtService.BearerHeaderCheck(c.GetHeader("Authorization"))
		if err != nil {
			util.ErrorAbort(c, http.StatusUnauthorized, err)
		}

		// 从 token 中获取用户信息并保存到 context 中
		c.Set("userId", claims.Subject)

		c.Next()
	}
}
