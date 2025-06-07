package repository

import (
	"database/sql"
	"newblog/internal/model"
)

type CategoryRepository interface {
	Exist(id int64) (bool, error)
	Insert(name string) (int64, error)
	Delete(aid int64) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (a *categoryRepository) ListAll() (*[]model.Category, error) {
	query := `
SELECT id, name
FROM category
`

	cateRow, err := a.db.Query(query)
	defer cateRow.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var cates []model.Category
	for cateRow.Next() {
		var cate model.Category
		if err := cateRow.Scan(
			&cate.ID,
			&cate.Name,
		); err != nil {
			return nil, err
		}
		cates = append(cates, cate)
	}

	return &cates, nil
}

func (a *categoryRepository) Exist(id int64) (bool, error) {
	query := `
SELECT COUNT(id) as count
FROM category
WHERE id = ?
`

	var exist int
	err := a.db.QueryRow(query, id).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist != 0, nil
}

func (a *categoryRepository) Insert(name string) (int64, error) {
	query := `
INSERT INTO category(name)
VALUES(?)
`

	res, err := a.db.Exec(query, name)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (a *categoryRepository) Delete(id int64) error {
	query := `
DELETE FROM category
WHERE id = ?
`

	_, err := a.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
