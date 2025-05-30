package model

type Web struct {
	Open    bool   `json:"open"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	NavList *[]Nav `json:"nav_list"`
}

type Nav struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
