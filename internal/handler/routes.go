package handler

import (
	"net/http"
	"newblog/internal/config"
	"newblog/internal/middleware"
	"newblog/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(svc *service.Container) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Global.Server.Addr},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	adminHandler := NewAdminHandler(svc.AdminService)
	articleHandler := NewArticleHandler(svc.ArticleService)
	WebHandler := NewWebHandler(svc.WebService)

	web := r.Group("/web")
	web.GET("/info", WebHandler.Info)

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
