package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "缺少认证信息"))
			return
		}
		authArr := strings.Split(auth, " ")
		if len(authArr) != 2 || authArr[0] != "Bearer" || authArr[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证格式错误"))
			return
		}
		if tokenBytes, err := os.ReadFile(tokenFile); err == nil {
			if string(tokenBytes) != authArr[1] {
				c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证信息错误"))
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证已过期"))
			return
		}

		_, err := jwt.Parse(authArr[1], func(token *jwt.Token) (any, error) {
			return key, nil
		})
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证格式错误"))
				return
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证被篡改"))
				return
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				os.Remove(tokenFile)
				c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证已过期"))
				return
			default:
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, Error(401, "认证信息错误"))
				return
			}
			// } else if _, ok := token.Claims.(*JWTClaims); ok {
			// } else {
			// 	c.AbortWithStatusJSON(http.StatusInternalServerError, Error(500, "unknown claims type"))
		}

		c.Next()
	}
}
