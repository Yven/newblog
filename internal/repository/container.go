package repository

import "database/sql"

type Container struct {
	ArticleRepo ArticleRepository
}

func NewRepositoryContainer(db *sql.DB) *Container {
	return &Container{
		ArticleRepo: NewArticleRepository(db),
	}
}
