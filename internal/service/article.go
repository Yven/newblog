package service

import (
	"log"
	"newblog/internal/cron"
	"newblog/internal/model"
	"newblog/internal/repository"
	"newblog/internal/validate"
	"time"
)

type ArticleService interface {
	List(search *validate.List, getAll bool) (data []*model.Article, err error)
	ListByYear(search *validate.List, getAll bool) (data []*model.ArticleList, err error)
	Info(slug string, getAll bool) (data *model.Article, err error)
	Edit(slug string, newContent string) error
	Delete(slug string) error
	RealDelete(slug string) error
	Recover(slug string) error
	Create(article *model.Article) (*model.Article, error)
	Sync() error
}

type articleService struct {
	db repository.ArticleRepository
}

func NewArticleService(db repository.ArticleRepository) *articleService {
	return &articleService{
		db: db,
	}
}

func (s *articleService) List(search *validate.List, getAll bool) (data []*model.Article, err error) {
	return s.db.List(search, getAll)
}

func (s *articleService) ListByYear(search *validate.List, getAll bool) (data []*model.ArticleList, err error) {
	list, err := s.db.List(search, getAll)
	if err != nil {
		return nil, err
	}

	// 根据时间分类
	var newList []*model.ArticleList
	for _, item := range list {
		var listItem model.ArticleListItem
		log.Println(item.CreateTime)
		time, _ := time.Parse("2006-01-02 15:04:05", item.CreateTime)

		listItem = model.ArticleListItem{
			ID:         item.ID,
			Slug:       item.Slug,
			Cid:        item.Cid,
			Title:      item.Title,
			Category:   item.Category,
			Date:       time.Format("01-02"),
			DeleteTime: item.DeleteTime,
			TagList:    item.TagList,
		}

		find := false
		for _, newItem := range newList {
			if newItem.Year == time.Year() {
				newItem.Item = append(newItem.Item, listItem)
				find = true
				break
			}
		}
		if !find {
			// 新增
			newList = append(newList, &model.ArticleList{Year: time.Year(), Item: []model.ArticleListItem{listItem}})
		}
	}

	return newList, nil
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
	return s.db.Insert(article)
}

func (s *articleService) Sync() error {
	nb := &cron.NotionBlog{}
	return nb.Exec()
}
