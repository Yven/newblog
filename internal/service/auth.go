package service

import (
	"errors"
	"fmt"
	"newblog/internal/config"
	"newblog/internal/global"
	"newblog/internal/model"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AuthService interface {
	GetAuthFileName(subject string) string
	ReadAuthFile(subject string) string
	WriteAuthFile(subject string, token *model.Token) error
	BearerHeaderCheck(bearer string) (*model.JWTClaims, error)
}

type authService struct {
}

func NewAuthService() *authService {
	return &authService{}
}

func (s *authService) GetAuthFileName(subject string) string {
	return fmt.Sprintf("%s/auth-%s.auth", strings.TrimRight(config.Global.Auth.LocalPath, "/"), subject)
}

func (s *authService) ReadAuthFile(subject string) string {
	filename := s.GetAuthFileName(subject)
	if tokenBytes, err := os.ReadFile(filename); err == nil {
		_, err := global.JwtService.Check(string(tokenBytes))
		if err == nil {
			return string(tokenBytes)
		}
		os.Remove(filename)
	}

	return ""
}

func (s *authService) WriteAuthFile(subject string, token *model.Token) error {
	filename := s.GetAuthFileName(subject)

	// 保存token到本地
	os.Remove(filename)

	// 判断文件所在目录是否存在
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录，权限设置为 0755
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	err := os.WriteFile(filename, []byte(token.Token), 0644)
	if err != nil {
		return err
	}

	// 定时删除
	// 计算token过期时间与当前时间的差值作为定时器时间
	settime := time.Until(time.Unix(token.Exp, 0))
	time.AfterFunc(settime, func() {
		os.Remove(filename)
	})

	return nil
}

func (s *authService) BearerHeaderCheck(bearer string) (*model.JWTClaims, error) {
	authArr := strings.Split(bearer, " ")
	if len(authArr) != 2 || authArr[0] != "Bearer" || authArr[1] == "" {
		return nil, errors.New("认证格式错误")
	}
	auth := authArr[1]

	claims, err := global.JwtService.Check(auth)
	if err != nil {
		return nil, err
	}

	if tokenBytes := s.ReadAuthFile(claims.Subject); tokenBytes != "" {
		if string(tokenBytes) != auth {
			return nil, errors.New("请重新登录")
		}
	}

	return claims, nil
}
