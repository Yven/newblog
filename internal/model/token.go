package model

type Token struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}
