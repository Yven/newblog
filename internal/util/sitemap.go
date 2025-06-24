package util

import (
	"newblog/internal/model"
	"os"
	"strings"
)

func Sitemap(path string, list []*model.Article) error {
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
	for _, item := range list {
		var str string
		str = strings.ReplaceAll(tpl, "{{url}}", url+"/#"+item.Slug)
		str = strings.ReplaceAll(str, "{{lastmod}}", item.CreateTime)
		urlList = append(urlList, str)
	}

	xml = strings.ReplaceAll(xml, "{{list}}", strings.Join(urlList, ""))

	return os.WriteFile(strings.TrimRight(path, "/")+"/sitemap.xml", []byte(xml), 0644)
}
