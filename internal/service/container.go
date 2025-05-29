package service

import "newblog/internal/repository"

type Container struct {
	ArticleService ArticleService
	AdminService   AdminService
}

func NewServiceContainer(repo *repository.Container) *Container {
	return &Container{
		ArticleService: NewArticleService(repo.ArticleRepo),
		AdminService:   NewAdminService(repo.AdminRepo),
	}
}
