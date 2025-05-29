package repository

import (
	"database/sql"
	"newblog/internal/config"
	"newblog/internal/model"
)

type WebRepository interface {
	Info() (*model.Web, error)
}

type webRepository struct {
	db *sql.DB
}

func NewWebRepository(db *sql.DB) WebRepository {
	return &webRepository{db: db}
}

func (w *webRepository) Info() (*model.Web, error) {
	conf := config.Global.Web

	var navList []model.Nav
	for _, v := range conf.NavList {
		navList = append(navList, model.Nav{
			Title: v.Title,
			Path:  v.Path,
		})
	}

	return &model.Web{
		Title:   conf.Title,
		Desc:    conf.Desc,
		NavList: &navList,
	}, nil
}
