package database

import "database/sql"

type ArticleListItem struct {
	ID       int64  `json:"id"`
	Slug     string `json:"slug"`
	Cid      int64  `json:"cid"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Date     string `json:"date"`
}

type ArticleList struct {
	Year int               `json:"year"`
	Item []ArticleListItem `json:"item"`
}

func (s *service) List() (*[]ArticleList, error) {
	query := `
SELECT a.id, a.slug, a.title, a.cid, strftime('%m-%d', a.create_time) as date, strftime('%Y', create_time) AS year, c.name AS category
FROM article AS a
LEFT JOIN category AS c ON a.cid = c.id
WHERE delete_time IS NULL
ORDER BY create_time DESC
`

	rows, err := s.db.Query(query)
	defer rows.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var list []ArticleList

	i := 0
	for rows.Next() {
		var item ArticleListItem
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
			list = append(list, ArticleList{Year: year, Item: nil})
		}
		if list[i].Year != year {
			// 新增
			i = i + 1
			list = append(list, ArticleList{Year: year, Item: nil})
		}
		list[i].Item = append(list[i].Item, item)
	}

	return &list, nil
}
