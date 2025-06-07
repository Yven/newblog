package repository

import (
	"database/sql"
	"fmt"
	"newblog/internal/model"
	"strings"
)

type TagRepository interface {
	List(aid int64) (*[]model.Tag, error)
	ListAll() (*[]model.Tag, error)
	Exist(tags *[]model.Tag) (bool, error)
	Insert(name string) (*model.Tag, error)
	DeleteRelate(aid int64) error
	Relate(aid int64, tags *[]model.Tag) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

func (a *tagRepository) ListAll() (*[]model.Tag, error) {
	query := `
SELECT id, name
FROM tag
`

	tagsRow, err := a.db.Query(query)
	defer tagsRow.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var tags []model.Tag
	for tagsRow.Next() {
		var tag model.Tag
		if err := tagsRow.Scan(
			&tag.ID,
			&tag.Name,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return &tags, nil
}

func (a *tagRepository) List(aid int64) (*[]model.Tag, error) {
	query := `
SELECT t.id, t.name
FROM tag AS t
LEFT JOIN article_tag AS at ON t.id = at.tid
WHERE at.aid = ?
`

	tagsRow, err := a.db.Query(query, aid)
	defer tagsRow.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var tags []model.Tag
	for tagsRow.Next() {
		var tag model.Tag
		if err := tagsRow.Scan(
			&tag.ID,
			&tag.Name,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return &tags, nil
}

func (a *tagRepository) Exist(tags *[]model.Tag) (bool, error) {
	placeholders := make([]string, len(*tags))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
SELECT COUNT(id) as count
FROM tag
WHERE id IN (%s)
`, strings.Join(placeholders, ","))

	ids := make([]any, len(*tags))
	for i, tag := range *tags {
		ids[i] = tag.ID
	}

	var count int
	if err := a.db.QueryRow(query, ids...).Scan(&count); err != nil {
		return false, err
	}

	return count == len(*tags), nil
}

func (a *tagRepository) Insert(name string) (*model.Tag, error) {
	query := `
INSERT INTO tag(name)
VALUES(?)
`

	res, err := a.db.Exec(query, name)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &model.Tag{
		ID:   id,
		Name: name,
	}, nil
}

func (a *tagRepository) DeleteRelate(aid int64) error {
	query := `
DELETE FROM article_tag
WHERE aid = ?
`

	_, err := a.db.Exec(query, aid)
	if err != nil {
		return err
	}

	return nil
}

func (a *tagRepository) Relate(aid int64, tags *[]model.Tag) error {
	query := `
INSERT INTO article_tag(aid, tid)
VALUES(?, ?)
`

	for _, tag := range *tags {
		_, err := a.db.Exec(query, aid, tag.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
