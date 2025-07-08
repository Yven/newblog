package cron

import (
	"errors"
	"newblog/internal/global"
	"newblog/internal/model"
	"newblog/internal/repository"
	"newblog/internal/util"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Yven/notion_blog/client"
	"github.com/Yven/notion_blog/filter"
)

type NotionBlog struct{}

func (n *NotionBlog) GetRetryTimes() int {
	return 3
}

func (n *NotionBlog) Exec() error {
	key := os.Getenv("NOTION_KEY")
	DbId := os.Getenv("NOTION_DB_ID")

	notion := client.NewClient(key)
	list, err := notion.NewDb(DbId).Query(client.QueryDatabase{
		Filter: filter.Status("Status").Equal("waiting"),
	})
	if err != nil {
		return err
	}
	if len(list.GetContent()) == 0 {
		return nil
	}

	// 初始化数据库操作
	articleRepo := repository.NewArticleRepository(global.DbInstance)
	tagRepo := repository.NewTagRepository(global.DbInstance)
	cateRepo := repository.NewCategoryRepository(global.DbInstance)

	// 获取所有标签
	allTagList, _ := tagRepo.ListAll()
	inTagList := func(name string) *model.Tag {
		for _, item := range allTagList {
			if item.Name == name {
				return item
			}
		}
		return nil
	}

	var allErr []error

	for _, pageItem := range list.GetContent() {
		page, err := notion.NewBlock(pageItem.Id).Children(pageItem, client.BaseQuery{})
		if err != nil {
			return err
		}

		content := page.ToMarkdown()
		imgRootPath := page.Property.Get("Slug").(string) + "/"

		// 使用正则表达式匹配HTML中的img标签
		imgReg := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)
		matches := imgReg.FindAllStringSubmatch(content, -1)

		cos := util.NewCos(
			os.Getenv("COS_URL"),
			os.Getenv("COS_SECRET_ID"),
			os.Getenv("COS_SECRET_Key"),
		)

		// 遍历所有匹配到的img标签
		for _, match := range matches {
			if len(match) > 1 {
				imgSrc := match[1]

				fileContent, downloadErr := util.DownloadFile(imgSrc)
				if downloadErr != nil {
					allErr = append(allErr, downloadErr)
					continue
				}

				// 上传图片到 COS 中，替换图片
				imgName := strings.Split(imgSrc[strings.LastIndex(imgSrc, "/")+1:], "?")[0]
				path, uploadErr := cos.UploadStream(imgRootPath+imgName, strings.NewReader(string(fileContent)))
				if uploadErr != nil {
					allErr = append(allErr, uploadErr)
					continue
				}
				content = strings.Replace(content, imgSrc, path, -1)
			}
		}

		// 处理标签
		var tagList []*model.Tag
		for _, tag := range page.Property.Get("Tag").([]string) {
			var tagModel *model.Tag
			if ex := inTagList(tag); ex != nil {
				tagModel = ex
			} else {
				tagModel, _ = tagRepo.Insert(tag)
			}

			tagList = append(tagList, tagModel)
		}

		// 处理分类
		categoryName := page.Property.Get("Category").(string)
		categoryModel, err := cateRepo.GetByName(categoryName)
		var cid int64
		if err != nil {
			allErr = append(allErr, err)
			cid = 0
		} else {
			cid = categoryModel.ID
		}
		if categoryModel == nil {
			cid, _ = cateRepo.Insert(categoryName)
		}

		// 新增
		_, insertErr := articleRepo.Insert(&model.Article{
			Slug:       page.Property.Get("Slug").(string),
			Title:      page.Property.Get("Name").(string),
			Content:    content,
			Cid:        cid,
			CreateTime: page.Property.CreatedTime,
			UpdateTime: page.Property.LastEditedTime,
			TagList:    tagList,
		})
		if insertErr != nil {
			allErr = append(allErr, insertErr)
			continue
		}

		// 更新文章状态
		updateData := filter.Date("Publish Time").Set("start", time.Now().Format("2006-01-02")).And(filter.Status("Status").Set("name", "publish"))
		updateErr := notion.NewPage(pageItem.Id).Update(client.UpdatePage{
			Properties: updateData,
		})
		if updateErr != nil {
			allErr = append(allErr, updateErr)
			continue
		}
	}

	// 更新 sitemap
	data, listErr := articleRepo.List(nil, false)
	if listErr == nil {
		util.Sitemap("./public", data)
	}

	// 返回所有操作中的错误，方便记录日志
	if len(allErr) > 0 {
		errMsgs := make([]string, len(allErr))
		for i, err := range allErr {
			errMsgs[i] = err.Error()
		}
		return errors.New(strings.Join(errMsgs, "\n;"))
	} else {
		return nil
	}
}
