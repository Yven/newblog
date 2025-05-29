package repository

import (
	"database/sql"
	"errors"
	"newblog/internal/config"
	"newblog/internal/model"
)

type AdminRepository interface {
	Info(name string, password string) (*model.Admin, error)
}

type adminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (a *adminRepository) Info(name string, password string) (*model.Admin, error) {
	if name != config.Global.Auth.User || password != config.Global.Auth.Password {
		return nil, errors.New("用户名或密码错误")
	}

	return &model.Admin{
		ID:       config.Global.Auth.Id,
		Username: config.Global.Auth.User,
	}, nil
}
