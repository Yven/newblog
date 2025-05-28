package handler

import (
	"errors"
	"net/http"
	"newblog/internal/service"
	"newblog/internal/util"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (s *AdminHandler) Login(c *gin.Context) {
	// 从 form 表单获取数据
	postUser := c.PostForm("username")
	postPassword := c.PostForm("password")

	// 如果 form 表单为空，尝试从 json 获取数据
	if postUser == "" || postPassword == "" {
		var jsonData struct {
			User     string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&jsonData); err == nil {
			postUser = jsonData.User
			postPassword = jsonData.Password
		}
	}

	// 如果两种方式都没有获取到数据，返回错误
	if postUser == "" || postPassword == "" {
		util.Error(c, http.StatusBadRequest, errors.New("缺少用户名或密码"))
		return
	}

	token, err := s.adminService.Login(postUser, postPassword)
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
