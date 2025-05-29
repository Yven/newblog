package service

import "newblog/internal/repository"

type Container struct {
	ArticleService ArticleService
	AdminService   AdminService
	WebService     WebService
}

func NewServiceContainer(repo *repository.Container) *Container {
	return &Container{
		ArticleService: NewArticleService(repo.ArticleRepo),
		AdminService:   NewAdminService(repo.AdminRepo),
		WebService:     NewWebService(repo.WebRepo),
	}
}
