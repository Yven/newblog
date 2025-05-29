package config

import (
	"log"
	"newblog/internal/model"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Global *model.Config

func InitConfig() {
	defaultPath := "./internal/config"
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load(defaultPath + "/.env")

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

	var c model.Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	Global = &c

	log.Println("配置初始化完成")
}
