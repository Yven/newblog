package util

import (
	"newblog/internal/model"
	"os"
	"strconv"
	"strings"
)

func Sitemap(path string, list *[]model.ArticleList) error {
	url := "https://yvenchang.cn"
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>` + url + `</loc>
    <lastmod>2025-06-12</lastmod>
  </sitemap>{{list}}
</sitemapindex>`

	tpl := `
  <sitemap>
    <loc>{{url}}</loc>
    <lastmod>{{lastmod}}</lastmod>
  </sitemap>`

	var urlList []string
	for _, x := range *list {
		year := x.Year
		for _, item := range x.Item {
			var str string
			str = strings.ReplaceAll(tpl, "{{url}}", url+"/#"+item.Slug)
			str = strings.ReplaceAll(str, "{{lastmod}}", strconv.Itoa(year)+"-"+item.Date)
			urlList = append(urlList, str)
		}
	}

	xml = strings.ReplaceAll(xml, "{{list}}", strings.Join(urlList, ""))

	return os.WriteFile(strings.TrimRight(path, "/")+"/sitemap.xml", []byte(xml), 0644)
}
