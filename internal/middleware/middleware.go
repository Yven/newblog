package middleware

import (
	"errors"
	"net/http"
	"newblog/internal/global"
	"newblog/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			util.ErrorAbort(c, http.StatusUnauthorized, errors.New("缺少认证信息"))
			return
		}
		authArr := strings.Split(auth, " ")
		if len(authArr) != 2 || authArr[0] != "Bearer" || authArr[1] == "" {
			util.ErrorAbort(c, http.StatusUnauthorized, errors.New("认证格式错误"))
			return
		}

		jwt := util.NewJwt(global.Token.Key, global.Token.Path)
		if ok, err := jwt.Parse(authArr[1]); !ok {
			util.ErrorAbort(c, http.StatusUnauthorized, err)
		}

		c.Next()
	}
}
