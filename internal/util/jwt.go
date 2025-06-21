package util

import (
	"errors"
	"log"
	"newblog/internal/config"
	"newblog/internal/model"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GetToken(userId string) (*model.Token, error)
	Check(token string) (*JWTClaims, error)
	Cancel()
	BearerHeaderCheck(auth string) (*JWTClaims, error)
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
	appendTime := time.Hour
	exptime := nowtime.Add(appendTime)
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
		time.AfterFunc(appendTime, func() {
			os.Remove(j.localPath)
		})
	}

	return &model.Token{
		Token: token,
		Exp:   exptime.Unix(),
	}, nil
}

func (j *jwtService) Check(token string) (*JWTClaims, error) {
	if tokenBytes, err := os.ReadFile(j.localPath); err == nil {
		if string(tokenBytes) != token {
			return nil, errors.New("认证信息错误")
		}
	} else {
		return nil, errors.New("认证已过期")
	}

	// 解析token
	claims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return j.key, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("认证格式错误")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("认证被篡改")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			os.Remove(j.localPath)
			return nil, errors.New("认证已过期")
		default:
			log.Println(err)
			return nil, errors.New("认证信息错误")
		}
	}

	return claims.Claims.(*JWTClaims), nil
}

func (j *jwtService) Cancel() {
	os.Remove(j.localPath)
}

func (j *jwtService) BearerHeaderCheck(auth string) (*JWTClaims, error) {
	if auth == "" {
		return nil, errors.New("缺少认证信息")
	}
	authArr := strings.Split(auth, " ")
	if len(authArr) != 2 || authArr[0] != "Bearer" || authArr[1] == "" {
		return nil, errors.New("认证格式错误")
	}

	claims, err := j.Check(authArr[1])
	if claims == nil {
		return nil, err
	}

	return claims, nil
}
