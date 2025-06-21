package util

import (
	"errors"
	"log"
	"newblog/internal/config"
	"newblog/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GetToken(userId string) (*model.Token, error)
	Check(token string) (*model.JWTClaims, error)
}

type jwtService struct {
	key    []byte
	claims *model.JWTClaims
}

func NewJwt(key string) *jwtService {
	return &jwtService{
		key: []byte(key),
	}
}

func (j *jwtService) GetToken(userId string) (*model.Token, error) {
	// 生成新的token
	nowtime := time.Now()
	appendTime := time.Hour
	exptime := nowtime.Add(appendTime)
	j.claims = &model.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
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
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token: token,
		Exp:   exptime.Unix(),
	}, nil
}

func (j *jwtService) Check(token string) (*model.JWTClaims, error) {
	// 解析token
	claims, err := jwt.ParseWithClaims(token, &model.JWTClaims{}, func(token *jwt.Token) (any, error) {
		return j.key, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("认证格式错误")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("认证被篡改")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.New("认证已过期")
		default:
			log.Println(err)
			return nil, errors.New("认证信息错误")
		}
	}

	return claims.Claims.(*model.JWTClaims), nil
}
