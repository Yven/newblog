package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"newblog/internal/config"
	"newblog/internal/cron"
	"newblog/internal/global"
	"newblog/internal/handler"
	"newblog/internal/logger"
	"newblog/internal/repository"
	"newblog/internal/service"
	"newblog/internal/util"
)

func gracefulShutdown(apiServer *http.Server, cronService *cron.CronService, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// 停止定时任务
	cronService.Stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// 配置初始化
	config.InitConfig()
	// JWT 初始化
	global.JwtService = util.NewJwt(config.Global.Auth.SignKey)
	// 数据库初始化
	db := repository.InitDb()

	// 日志初始化
	logger.Init(config.Global.Log.Path, config.Global.Log.Level)

	// 容器初始化
	repo := repository.NewRepositoryContainer(db)
	svc := service.NewServiceContainer(repo)

	// 定时任务初始化
	cronService := cron.NewCronService()
	cronService.Register()

	// http 服务初始化
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Global.Server.Port),
		Handler:      handler.RegisterRoutes(svc),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, cronService, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
