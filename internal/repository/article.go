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
	List(search *validate.List, getAll bool) ([]*model.Article, error)
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

	// 将时间戳转换为time.Time类型,并根据系统时区格式化输出
	if article.CreateTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", article.CreateTime)
		if err == nil {
			article.CreateTime = t.Local().Format("2006-01-02 15:04:05")
		}
	}
	if article.UpdateTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", article.UpdateTime)
		if err == nil {
			article.UpdateTime = t.Local().Format("2006-01-02 15:04:05")
		}
	}
	if article.DeleteTime != nil {
		t, err := time.Parse("2006-01-02 15:04:05", *article.DeleteTime)
		if err == nil {
			time := t.Local().Format("2006-01-02 15:04:05")
			article.DeleteTime = &time
		}
	}

	article.TagList, _ = NewTagRepository(a.db).List(article.ID)

	return &article, nil
}

func (a *articleRepository) List(search *validate.List, getAll bool) ([]*model.Article, error) {
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
datetime(a.create_time, 'unixepoch') as create_time,
datetime(a.update_time, 'unixepoch') as update_time,
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

	var list []*model.Article

	for rows.Next() {
		var item model.Article
		if err := rows.Scan(
			&item.ID,
			&item.Slug,
			&item.Title,
			&item.Cid,
			&item.CreateTime,
			&item.UpdateTime,
			&item.Category,
			&item.DeleteTime,
		); err != nil {
			return nil, err
		}

		item.TagList, _ = NewTagRepository(a.db).List(item.ID)

		list = append(list, &item)
	}

	return list, nil
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
INSERT INTO article (slug, title, content, cid, create_time, update_time)
VALUES (?, ?, ?, ?, ?, ?)
`

	// 设置时间
	createTime := time.Now().Unix()
	updateTime := createTime
	if article.CreateTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", article.CreateTime)
		if err == nil {
			createTime = t.Unix()
		}
	}
	if article.UpdateTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", article.UpdateTime)
		if err == nil {
			updateTime = t.Unix()
		}
	}

	res, err := a.db.Exec(query,
		article.Slug,
		article.Title,
		article.Content,
		article.Cid,
		createTime,
		updateTime,
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
