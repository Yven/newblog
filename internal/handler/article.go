package handler

import (
	"net/http"
	"newblog/internal/service"
	"newblog/internal/util"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) List(c *gin.Context) {
	data, err := h.articleService.List()

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, data)
	return
}

func (h *ArticleHandler) Info(c *gin.Context) {
	slug := c.Param("slug")
	article, err := h.articleService.Info(slug)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, article)
	return
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	slug := c.Param("slug")
	newContent := c.PostForm("content")

	err := h.articleService.Edit(slug, newContent)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
	return
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	slug := c.Param("slug")
	err := h.articleService.Delete(slug)

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, nil)
	return
}
