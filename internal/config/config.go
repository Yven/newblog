package config

import (
	"log"
	"log/slog"
	"newblog/internal/model"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Global *model.Config

func InitConfig() {
	defaultPath := "./internal/config"
	// 加载 .env 文件（如果存在）（默认读根目录）
	_ = godotenv.Load("./.env")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultPath)

	// 支持嵌套结构的环境变量覆盖：server.port => SERVER_PORT
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 支持通过环境变量覆盖配置
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 将SERVER_ADDR环境变量按逗号分割为字符串数组
	if addr := viper.GetString("SERVER_ADDR"); addr != "" {
		viper.Set("server.addr", strings.Split(addr, ","))
	}

	// 将字符串日志级别映射到 slog.Level
	logLevel := strings.ToLower(viper.GetString("log.level"))
	switch logLevel {
	case "debug":
		viper.Set("log.level", slog.LevelDebug)
	case "info":
		viper.Set("log.level", slog.LevelInfo)
	case "warn":
		viper.Set("log.level", slog.LevelWarn)
	case "error":
		viper.Set("log.level", slog.LevelError)
	default:
		viper.Set("log.level", slog.LevelInfo)
	}

	reload := func() {
		var newConfig model.Config
		if err := viper.Unmarshal(&newConfig); err != nil {
			log.Printf("热加载配置失败: %v", err)
			return
		}

		Global = &newConfig
		log.Println("配置已热更新")
	}

	// 初始加载
	reload()

	// 启动监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("检测到配置文件变化: %s", e.Name)
		reload()
	})
}
