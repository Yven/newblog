package service

import (
	"newblog/internal/global"
	"newblog/internal/model"
	"newblog/internal/repository"
	"os"
)

type AdminService interface {
	Login(postUser string, postPassword string) (*model.Token, error)
	Logout(subject string) bool
}

type adminService struct {
	db          repository.AdminRepository
	authService AuthService
}

func NewAdminService(db repository.AdminRepository, authService AuthService) *adminService {
	return &adminService{
		db:          db,
		authService: authService,
	}
}

func (s *adminService) Login(postUser string, postPassword string) (*model.Token, error) {
	admin, err := s.db.Info(postUser, postPassword)
	if err != nil {
		return nil, err
	}

	if tokenBytes := s.authService.ReadAuthFile(postUser); tokenBytes != "" {
		claim, checkErr := global.JwtService.Check(tokenBytes)
		if checkErr == nil {
			return &model.Token{
				Token: tokenBytes,
				Exp:   claim.ExpiresAt.Unix(),
			}, nil
		}
	}

	token, err := global.JwtService.GetToken(admin.Username)
	if err != nil {
		return nil, err
	}

	if err := s.authService.WriteAuthFile(postUser, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *adminService) Logout(subject string) bool {
	os.Remove(s.authService.GetAuthFileName(subject))
	return true
}
