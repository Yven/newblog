package service

import (
	"newblog/internal/model"
	"newblog/internal/repository"
)

type WebService interface {
	Info() (*model.Web, error)
}

type WebServiceContainer struct {
	db repository.WebRepository
}

func NewWebService(db repository.WebRepository) *WebServiceContainer {
	return &WebServiceContainer{
		db: db,
	}
}

func (w *WebServiceContainer) Info() (*model.Web, error) {
	return w.db.Info()
}
