package database

import (
	"database/sql"
	"time"
)

func (s *service) init() error {
	q := `
CREATE TABLE IF NOT EXISTS article(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	cid INTEGER DEFAULT 0,
	create_time TEXT NOT NULL DEFAULT (datetime('now')),
	update_time TEXT NOT NULL DEFAULT (datetime('now')),
	delete_time TEXT
);
CREATE TABLE IF NOT EXISTS category(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS article_tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	aid INTEGER NOT NULL,
	tid INTEGER NOT NULL
);
`

	_, err := s.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

type Article struct {
	ID           int64   `json:"id"`
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	Cid          int64   `json:"cid"`
	CategoryName string  `json:"category_name"`
	CreateTime   string  `json:"create_time"`
	UpdateTime   string  `json:"update_time"`
	DeleteTime   *string `json:"delete_time,omitempty"`
}

func (s *service) Search(slug string) (*Article, error) {
	query := `
SELECT a.id, a.slug, a.title, a.content, a.cid, a.create_time, a.update_time, a.delete_time, c.name AS category_name
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
WHERE slug = ? AND delete_time IS NULL
LIMIT 1
`

	row := s.db.QueryRow(query, slug)

	var article Article
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

func (s *service) Edit(slug string, newContent string) error {
	query := `
UPDATE article
SET content = ?
WHERE slug = ? AND delete_time IS NULL
`

	_, err := s.db.Exec(query, newContent, slug)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) Delete(slug string) error {
	query := `
UPDATE article
SET delete_time = ?
WHERE slug = ?
`

	var format = "2006-01-02 15:04:05"
	_, err := s.db.Exec(query, time.Now().Format(format), slug)
	if err != nil {
		return err
	}

	return nil
}
