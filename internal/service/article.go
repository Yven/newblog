package service

import (
	"newblog/internal/model"
	"newblog/internal/repository"
)

type ArticleService interface {
	List() (data *[]model.ArticleList, err error)
	Info(slug string) (data *model.Article, err error)
	Edit(slug string, newContent string) error
	Delete(slug string) error
}

type articleService struct {
	db repository.ArticleRepository
}

func NewArticleService(db repository.ArticleRepository) *articleService {
	return &articleService{
		db: db,
	}
}

func (s *articleService) List() (data *[]model.ArticleList, err error) {
	return s.db.List()
}

func (s *articleService) Info(slug string) (data *model.Article, err error) {
	return s.db.Info(slug)
}

func (s *articleService) Edit(slug string, newContent string) error {
	return s.db.Edit(slug, newContent)
}

func (s *articleService) Delete(slug string) error {
	return s.db.Delete(slug)
}
