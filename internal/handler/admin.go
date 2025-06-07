package handler

import (
	"errors"
	"net/http"
	"newblog/internal/service"
	"newblog/internal/util"
	"newblog/internal/validate"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (s *AdminHandler) Login(c *gin.Context) {
	var data validate.Admin
	if err := c.ShouldBindJSON(&data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			util.Error(c, http.StatusBadRequest, errors.New("字段格式错误: "+validationErrors.Error()))
			return
		}

		util.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := s.adminService.Login(data.Username, data.Password)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, err)
		return
	}

	util.Success(c, token)
	return
}

func (s *AdminHandler) Logout(c *gin.Context) {
	s.adminService.Logout()
	util.Success(c, nil)
	return
}
