package handler

import (
	"net/http"
	"newblog/internal/service"
	"newblog/internal/util"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	webService service.WebService
}

func NewWebHandler(webService service.WebService) *WebHandler {
	return &WebHandler{webService: webService}
}

func (h *WebHandler) Info(c *gin.Context) {
	data, err := h.webService.Info()

	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, data)
	return
}
