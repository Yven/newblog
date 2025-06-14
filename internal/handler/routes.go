package handler

import (
	"net/http"
	"newblog/internal/config"
	"newblog/internal/middleware"
	"newblog/internal/service"
	"newblog/internal/validate"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(svc *service.Container) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Global.Server.Addr,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("idStringList", validate.IdStringList)
	}

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

		authorized.POST("/content", articleHandler.Create)
		authorized.DELETE("/content/delete/:slug", articleHandler.RealDelete)
		authorized.GET("/content/recover/:slug", articleHandler.Recover)
	}

	return r
}
