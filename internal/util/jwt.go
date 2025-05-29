package util

import (
	"errors"
	"log"
	"newblog/internal/config"
	"newblog/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GetToken(userId string) (*model.Token, error)
	Check(token string) (bool, error)
	Cancel()
}

type jwtService struct {
	key       []byte
	localPath string
	claims    *JWTClaims
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

func NewJwt(key string, path string) *jwtService {
	return &jwtService{
		key:       []byte(key),
		localPath: path,
	}
}

func (j *jwtService) GetToken(userId string) (*model.Token, error) {
	// 检查是否存在已保存的token
	if tokenBytes, err := os.ReadFile(j.localPath); err == nil {
		token, _ := jwt.Parse(string(tokenBytes), func(token *jwt.Token) (any, error) {
			return j.key, nil
		})
		if token.Valid {
			exp, err := token.Claims.GetExpirationTime()
			if err == nil {
				return &model.Token{
					Token: string(tokenBytes),
					Exp:   exp.Unix(),
				}, nil
			}
		}
	}

	// 生成新的token
	nowtime := time.Now()
	exptime := nowtime.Add(time.Hour)
	j.claims = &JWTClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exptime),
			IssuedAt:  jwt.NewNumericDate(nowtime),
			NotBefore: jwt.NewNumericDate(nowtime),
			Issuer:    config.Global.Auth.Issuer,
			Subject:   userId,
			ID:        "0",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims)
	token, err := t.SignedString(j.key)

	// 保存token到本地
	os.Remove(j.localPath)
	err = os.WriteFile(j.localPath, []byte(token), 0644)
	// 定时删除
	if err != nil {
		go func() {
			time.Sleep(time.Hour)
			os.Remove(j.localPath)
		}()
	}

	return &model.Token{
		Token: token,
		Exp:   exptime.Unix(),
	}, nil
}

func (j *jwtService) Check(token string) (bool, error) {
	// 检查是否存在已保存的token
	if tokenBytes, err := os.ReadFile(j.localPath); err == nil {
		if string(tokenBytes) != token {
			return false, errors.New("认证信息错误")
		}
	} else {
		return false, errors.New("认证已过期")
	}

	// 解析token
	_, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.key, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return false, errors.New("认证格式错误")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return false, errors.New("认证被篡改")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			os.Remove(j.localPath)
			return false, errors.New("认证已过期")
		default:
			log.Println(err)
			return false, errors.New("认证信息错误")
		}
	}

	return true, nil
}

func (j *jwtService) Cancel() {
	os.Remove(j.localPath)
}
