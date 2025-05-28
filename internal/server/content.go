package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) list(c *gin.Context) {
	data, err := s.db.List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(&data))
	return
}

func (s *Server) info(c *gin.Context) {
	slug := c.Param("slug")
	article, err := s.db.Search(slug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(article))
	return
}

func (s *Server) edit(c *gin.Context) {
	slug := c.Param("slug")
	newContent := c.PostForm("content")
	err := s.db.Edit(slug, newContent)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
	return
}

func (s *Server) delete(c *gin.Context) {
	slug := c.Param("slug")
	err := s.db.Delete(slug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
	return
}
