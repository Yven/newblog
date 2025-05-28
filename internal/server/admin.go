package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	user      = os.Getenv("ADMIN_NAME")
	password  = os.Getenv("PASSWORD")
	key       = []byte(os.Getenv("JWT_SIGN_KEY"))
	tokenFile = os.Getenv("LOCAL_TOKEN_PATH")
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func (s *Server) login(c *gin.Context) {
	// 从 form 表单获取数据
	postUser := c.PostForm("username")
	postPassword := c.PostForm("password")

	// 如果 form 表单为空，尝试从 json 获取数据
	if postUser == "" || postPassword == "" {
		var jsonData struct {
			User     string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&jsonData); err == nil {
			postUser = jsonData.User
			postPassword = jsonData.Password
		}
	}

	// 如果两种方式都没有获取到数据，返回错误
	if postUser == "" || postPassword == "" {
		c.JSON(http.StatusBadRequest, Error(400, "缺少用户名或密码"))
		return
	}
	// postUser := c.PostForm("user")
	// postPassword := c.PostForm("password")

	if user == postUser && password == postPassword {
		// 检查是否存在已保存的token
		if tokenBytes, err := os.ReadFile(tokenFile); err == nil {
			token, _ := jwt.Parse(string(tokenBytes), func(token *jwt.Token) (any, error) {
				return key, nil
			})
			if token.Valid {
				exp, err := token.Claims.GetExpirationTime()
				if err == nil {
					c.JSON(http.StatusOK, Success(gin.H{"token": string(tokenBytes), "exp": exp.Unix()}))
					return
				}
			} else {
				os.Remove(tokenFile)
			}
		}

		nowtime := time.Now()
		exptime := nowtime.Add(time.Hour)
		claims := JWTClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(exptime),
				IssuedAt:  jwt.NewNumericDate(nowtime),
				NotBefore: jwt.NewNumericDate(nowtime),
				Issuer:    "yven_server",
				Subject:   postUser,
				ID:        "1",
			},
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := t.SignedString(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(500, err.Error()))
			return
		}

		os.Remove(tokenFile)
		err = os.WriteFile(tokenFile, []byte(token), 0644)
		if err != nil {
			go func() {
				time.Sleep(time.Hour)
				os.Remove(tokenFile)
			}()
		}
		c.JSON(http.StatusOK, Success(gin.H{"token": token, "exp": exptime.Unix()}))
		return
	} else {
		c.JSON(http.StatusUnauthorized, Error(401, "用户名或密码错误"))
		return
	}
}

func (s *Server) logout(c *gin.Context) {
	os.Remove(tokenFile)
	c.JSON(http.StatusOK, Success(nil))
	return
}
