package server

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("WEBSITE_ADDR")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/list", s.list)

	r.GET("/content/:slug", s.info)

	r.POST("/login", s.login)

	authorized := r.Group("/")
	authorized.Use(Auth())
	{
		authorized.POST("/logout", s.logout)
		authorized.POST("/content/:slug", s.edit)
		authorized.DELETE("/content/:slug", s.delete)
	}

	return r
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(data any) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
