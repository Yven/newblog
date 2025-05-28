package handler

import (
	"net/http"
	"newblog/internal/middleware"
	"newblog/internal/service"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(svc *service.Container) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("WEBSITE_ADDR")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	adminHandler := NewAdminHandler(svc.AdminService)
	articleHandler := NewArticleHandler(svc.ArticleService)

	r.POST("/login", adminHandler.Login)

	r.GET("/list", articleHandler.List)
	r.GET("/content/:slug", articleHandler.Info)

	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	{
		authorized.POST("/logout", adminHandler.Logout)

		authorized.POST("/content/:slug", articleHandler.Edit)
		authorized.DELETE("/content/:slug", articleHandler.Delete)
	}

	return r
}
