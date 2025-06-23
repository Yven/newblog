package repository

import (
	"database/sql"
	"errors"
	"newblog/internal/model"
	"newblog/internal/validate"
	"strings"
	"time"
)

type ArticleRepository interface {
	List(search *validate.List, getAll bool) (*[]model.ArticleList, error)
	Info(slug string, getAll bool) (*model.Article, error)
	Edit(slug string, newContent string) error
	Delete(slug string) error
	RealDelete(slug string) error
	Recover(slug string) error
	Insert(article *model.Article) (*model.Article, error)
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (a *articleRepository) Info(slug string, getAll bool) (*model.Article, error) {
	deleteWhere := ""
	if !getAll {
		deleteWhere = "AND a.delete_time IS NULL"
	}

	query := `
SELECT a.id, a.slug, a.title, a.content, a.cid,
datetime(a.create_time, 'unixepoch') as create_time,
datetime(a.update_time, 'unixepoch') as update_time,
strftime('%Y-%m-%d %H:%M:%S', datetime(a.delete_time, 'unixepoch')) as delete_time,
c.name AS category
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
WHERE slug = ? ` + deleteWhere + `
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
		&article.Category,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	article.TagList, _ = NewTagRepository(a.db).List(article.ID)

	return &article, nil
}

func (a *articleRepository) List(search *validate.List, getAll bool) (*[]model.ArticleList, error) {
	var where []string
	var args []any

	if search != nil {
		if search.Keyword != "" {
			where = append(where, "a.title LIKE ?")
			args = append(args, "%"+search.Keyword+"%")
		}
		if search.Category != 0 {
			where = append(where, "a.cid = ?")
			args = append(args, search.Category)
		}
		if search.Tag != 0 {
			where = append(where, "at.tid = ?")
			args = append(args, search.Tag)
		}
	}

	if !getAll {
		where = append(where, "a.delete_time IS NULL")
	}

	whereStr := strings.Join(where, " AND ")
	if whereStr != "" {
		whereStr = "WHERE " + whereStr
	}

	query := `
SELECT a.id, a.slug, a.title, a.cid,
strftime('%m-%d', datetime(a.create_time, 'unixepoch')) as date,
strftime('%Y', datetime(a.create_time, 'unixepoch')) as year,
c.name AS category,
strftime('%Y-%m-%d %H:%M:%S', datetime(a.delete_time, 'unixepoch')) as delete_time
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
LEFT JOIN article_tag AS at ON at.aid = a.id
 ` + whereStr + `
GROUP BY a.id
ORDER BY create_time DESC
`

	rows, err := a.db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

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
			&item.DeleteTime,
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
		item.TagList, _ = NewTagRepository(a.db).List(item.ID)

		list[i].Item = append(list[i].Item, item)
	}

	return &list, nil
}

func (a *articleRepository) Edit(slug string, newContent string) error {
	query := `
UPDATE article
SET content = ?
WHERE slug = ?
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

func (a *articleRepository) RealDelete(slug string) error {
	query := `
DELETE FROM article
WHERE slug = ?
`

	data, err := a.Info(slug, true)
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query, slug)
	if err != nil {
		return err
	}

	err = NewTagRepository(a.db).DeleteRelate(data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *articleRepository) Recover(slug string) error {
	query := `
UPDATE article
SET delete_time = ?
WHERE slug = ?
`

	_, err := a.db.Exec(query, nil, slug)
	if err != nil {
		return err
	}

	return nil
}

func (a *articleRepository) Insert(article *model.Article) (*model.Article, error) {
	if exist, _ := a.Info(article.Slug, true); exist != nil {
		return nil, errors.New("slug已存在")
	}
	if exist, _ := NewCategoryRepository(a.db).Exist(article.Cid); !exist {
		return nil, errors.New("类别不存在")
	}
	if exist, _ := NewTagRepository(a.db).Exist(article.TagList); !exist {
		return nil, errors.New("标签不存在")
	}

	query := `
INSERT INTO article (slug, title, content, cid)
VALUES (?, ?, ?, ?)
`

	res, err := a.db.Exec(query,
		article.Slug,
		article.Title,
		article.Content,
		article.Cid,
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	err = NewTagRepository(a.db).Relate(id, article.TagList)
	if err != nil {
		return nil, err
	}

	return a.Info(article.Slug, true)
}
