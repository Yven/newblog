package service

import (
	"newblog/internal/model"
	"newblog/internal/repository"
	"newblog/internal/util"
	"newblog/internal/validate"
)

type ArticleService interface {
	List(search *validate.List, getAll bool) (data *[]model.ArticleList, err error)
	Info(slug string, getAll bool) (data *model.Article, err error)
	Edit(slug string, newContent string) error
	Delete(slug string) error
	RealDelete(slug string) error
	Recover(slug string) error
	Create(article *model.Article) (*model.Article, error)
}

type articleService struct {
	db repository.ArticleRepository
}

func NewArticleService(db repository.ArticleRepository) *articleService {
	return &articleService{
		db: db,
	}
}

func (s *articleService) List(search *validate.List, getAll bool) (data *[]model.ArticleList, err error) {
	return s.db.List(search, getAll)
}

func (s *articleService) Info(slug string, getAll bool) (data *model.Article, err error) {
	return s.db.Info(slug, getAll)
}

func (s *articleService) Edit(slug string, newContent string) error {
	return s.db.Edit(slug, newContent)
}

func (s *articleService) Delete(slug string) error {
	return s.db.Delete(slug)
}

func (s *articleService) RealDelete(slug string) error {
	return s.db.RealDelete(slug)
}

func (s *articleService) Recover(slug string) error {
	return s.db.Recover(slug)
}

func (s *articleService) Create(article *model.Article) (*model.Article, error) {
	res, err := s.db.Insert(article)

	go func() {
		data, listErr := s.db.List(nil, false)
		if listErr == nil {
			util.Sitemap("./public", data)
		}
	}()

	return res, err
}
