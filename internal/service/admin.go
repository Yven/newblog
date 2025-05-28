package service

import (
	"errors"
	"newblog/internal/global"
	"newblog/internal/model"
	"newblog/internal/util"
)

type AdminService interface {
	Login(postUser string, postPassword string) (*model.Token, error)
	Logout() bool
}

type adminService struct{}

func NewAdminService() *adminService {
	return &adminService{}
}

func (s *adminService) Login(postUser string, postPassword string) (*model.Token, error) {
	if global.Admin.User == postUser && global.Admin.Password == postPassword {
		token, err := util.NewJwt(global.Token.Key, global.Token.Path).GetToken(postUser)
		if err != nil {
			return nil, err
		}

		return token, nil
	} else {
		return nil, errors.New("用户名或密码错误")
	}
}

func (s *adminService) Logout() bool {
	util.NewJwt(global.Token.Key, global.Token.Path).Cancel()
	return true
}
