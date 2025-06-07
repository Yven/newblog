package model

type Article struct {
	ID         int64   `json:"id"`
	Slug       string  `json:"slug"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Cid        int64   `json:"cid"`
	Category   *string `json:"category"`
	CreateTime string  `json:"create_time"`
	UpdateTime string  `json:"update_time"`
	DeleteTime *string `json:"delete_time"`
	TagList    *[]Tag  `json:"tag_list"`
}

type ArticleListItem struct {
	ID         int64   `json:"id"`
	Slug       string  `json:"slug"`
	Cid        int64   `json:"cid"`
	Title      string  `json:"title"`
	Category   *string `json:"category"`
	Date       string  `json:"date"`
	DeleteTime *string `json:"delete_time"`
	TagList    *[]Tag  `json:"tag_list"`
}

type ArticleList struct {
	Year int               `json:"year"`
	Item []ArticleListItem `json:"item"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
