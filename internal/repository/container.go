package repository

import "database/sql"

type Container struct {
	ArticleRepo ArticleRepository
	AdminRepo   AdminRepository
	WebRepo     WebRepository
}

func NewRepositoryContainer(db *sql.DB) *Container {
	return &Container{
		ArticleRepo: NewArticleRepository(db),
		AdminRepo:   NewAdminRepository(db),
		WebRepo:     NewWebRepository(db),
	}
}
