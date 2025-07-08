package handler

import (
	"errors"
	"log"
	"net/http"
	"newblog/internal/config"
	"newblog/internal/model"
	"newblog/internal/service"
	"newblog/internal/util"
	"newblog/internal/validate"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/feeds"
)

type ArticleHandler struct {
	articleService service.ArticleService
	authService    service.AuthService
}

func NewArticleHandler(articleService service.ArticleService, authService service.AuthService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		authService:    authService,
	}
}

func (h *ArticleHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	tag := c.Query("tag")
	category := c.Query("category")

	tid, _ := strconv.ParseInt(tag, 10, 64)
	cid, _ := strconv.ParseInt(category, 10, 64)
	search := &validate.List{
		Keyword:  keyword,
		Tag:      tid,
		Category: cid,
	}

	_, err := h.authService.BearerHeaderCheck(c.GetHeader("Authorization"))
	getAll := false
	if err == nil {
		getAll = true
	}

	data, err := h.articleService.ListByYear(search, getAll)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, data)
}

func (h *ArticleHandler) Info(c *gin.Context) {
	slug := c.Param("slug")

	_, err := h.authService.BearerHeaderCheck(c.GetHeader("Authorization"))
	getAll := false
	if err == nil {
		getAll = true
	}

	article, err := h.articleService.Info(slug, getAll)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, article)
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	slug := c.Param("slug")

	var data struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			util.Error(c, http.StatusBadRequest, errors.New("字段格式错误: "+validationErrors.Error()))
			return
		}

		util.Error(c, http.StatusBadRequest, err)
		return
	}

	err := h.articleService.Edit(slug, data.Content)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	slug := c.Param("slug")
	err := h.articleService.Delete(slug)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
}

func (h *ArticleHandler) RealDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := h.articleService.RealDelete(slug)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
}

func (h *ArticleHandler) Recover(c *gin.Context) {
	slug := c.Param("slug")
	err := h.articleService.Recover(slug)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var data validate.Article

	if err := c.ShouldBindJSON(&data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			util.Error(c, http.StatusBadRequest, errors.New("字段格式错误: "+validationErrors.Error()))
			return
		}

		util.Error(c, http.StatusBadRequest, err)
		return
	}

	tagList := strings.Split(data.TagList, ",")
	var tags []*model.Tag
	for _, tagStr := range tagList {
		tagStr = strings.TrimSpace(tagStr)
		if tagStr == "" {
			continue
		}
		tagId, _ := strconv.ParseInt(tagStr, 10, 64)
		tags = append(tags, &model.Tag{ID: tagId})
	}

	article := &model.Article{
		Slug:    data.Slug,
		Title:   data.Title,
		Content: data.Content,
		Cid:     data.Cid,
		TagList: tags,
	}

	res, err := h.articleService.Create(article)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, res)
}

func (h *ArticleHandler) Sync(c *gin.Context) {
	err := h.articleService.Sync()

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
}

func (h *ArticleHandler) Feed(c *gin.Context) {
	now := time.Now()
	url := "https://yvenchang.cn"
	author := feeds.Author{Name: "Yven Chang", Email: "yvenchang@163.com"}

	feed := &feeds.Feed{
		Title:       config.Global.Web.Title,
		Link:        &feeds.Link{Href: url},
		Description: config.Global.Web.Desc,
		Author:      &author,
		Created:     now,
	}

	articles, err := h.articleService.List(nil, false)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	for _, article := range articles {
		create, _ := time.Parse("2006-01-02 15:04:05", article.CreateTime)
		update, _ := time.Parse("2006-01-02 15:04:05", article.UpdateTime)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: url + "/#" + article.Slug},
			Description: "",
			Author:      &author,
			Created:     create,
			Updated:     update,
			Content:     article.Content,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rss))
}
