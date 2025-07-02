package util

import (
	"net/http"
	"newblog/internal/model"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message error) {
	c.Error(message)
	c.JSON(code, model.Response{
		Code:    code,
		Message: message.Error(),
		Data:    nil,
	})
}

func SuccessAbort(c *gin.Context, data any) {
	c.AbortWithStatusJSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func ErrorAbort(c *gin.Context, code int, message error) {
	// c.Error(message)
	c.AbortWithStatusJSON(code, model.Response{
		Code:    code,
		Message: message.Error(),
		Data:    nil,
	})
}
