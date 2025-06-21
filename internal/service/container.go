package service

import "newblog/internal/repository"

type Container struct {
	ArticleService ArticleService
	AdminService   AdminService
	WebService     WebService
	AuthService    AuthService
}

func NewServiceContainer(repo *repository.Container) *Container {
	authService := NewAuthService()
	return &Container{
		ArticleService: NewArticleService(repo.ArticleRepo),
		AdminService:   NewAdminService(repo.AdminRepo, authService),
		WebService:     NewWebService(repo.WebRepo),
		AuthService:    authService,
	}
}
