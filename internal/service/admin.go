package service

import (
	"newblog/internal/global"
	"newblog/internal/model"
	"newblog/internal/repository"
)

type AdminService interface {
	Login(postUser string, postPassword string) (*model.Token, error)
	Logout() bool
}

type adminService struct {
	db repository.AdminRepository
}

func NewAdminService(db repository.AdminRepository) *adminService {
	return &adminService{
		db: db,
	}
}

func (s *adminService) Login(postUser string, postPassword string) (*model.Token, error) {
	admin, err := s.db.Info(postUser, postPassword)
	if err != nil {
		return nil, err
	}

	token, err := global.JwtService.GetToken(admin.Username)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *adminService) Logout() bool {
	global.JwtService.Cancel()
	return true
}
