package model

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
}
