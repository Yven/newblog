package repository

import (
	"database/sql"
	"newblog/internal/model"
	"time"
)

type ArticleRepository interface {
	List(keyword string) (*[]model.ArticleList, error)
	Info(slug string) (*model.Article, error)
	Edit(slug string, newContent string) error
	Delete(slug string) error
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (a *articleRepository) Info(slug string) (*model.Article, error) {
	query := `
SELECT a.id, a.slug, a.title, a.content, a.cid,
datetime(a.create_time, 'unixepoch') as create_time,
datetime(a.update_time, 'unixepoch') as update_time,
a.delete_time, c.name AS category_name
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
WHERE slug = ? AND delete_time IS NULL
LIMIT 1
`

	row := a.db.QueryRow(query, slug)

	var article model.Article
	err := row.Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Content,
		&article.Cid,
		&article.CreateTime,
		&article.UpdateTime,
		&article.DeleteTime,
		&article.CategoryName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &article, nil
}

func (a *articleRepository) List(keyword string) (*[]model.ArticleList, error) {
	addWhere := ""
	if keyword != "" {
		addWhere = "AND a.title LIKE ?"
	}

	query := `
SELECT a.id, a.slug, a.title, a.cid,
strftime('%m-%d', datetime(a.create_time, 'unixepoch')) as date,
strftime('%Y', datetime(a.create_time, 'unixepoch')) as year,
c.name AS category
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
WHERE delete_time IS NULL ` + addWhere + `
ORDER BY create_time DESC
`

	rows, err := a.db.Query(query, "%"+keyword+"%")
	defer rows.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var list []model.ArticleList

	i := 0
	for rows.Next() {
		var item model.ArticleListItem
		var year int
		if err := rows.Scan(
			&item.ID,
			&item.Slug,
			&item.Title,
			&item.Cid,
			&item.Date,
			&year,
			&item.Category,
		); err != nil {
			return nil, err
		}
		if len(list) == 0 {
			// 初始化
			list = append(list, model.ArticleList{Year: year, Item: nil})
		}
		if list[i].Year != year {
			// 新增
			i = i + 1
			list = append(list, model.ArticleList{Year: year, Item: nil})
		}
		list[i].Item = append(list[i].Item, item)
	}

	return &list, nil
}

func (a *articleRepository) Edit(slug string, newContent string) error {
	query := `
UPDATE article
SET content = ?
WHERE slug = ? AND delete_time IS NULL
`

	_, err := a.db.Exec(query, newContent, slug)
	if err != nil {
		return err
	}

	return nil
}
func (a *articleRepository) Delete(slug string) error {
	query := `
UPDATE article
SET delete_time = ?
WHERE slug = ?
`

	_, err := a.db.Exec(query, time.Now().Unix(), slug)
	if err != nil {
		return err
	}

	return nil
}
