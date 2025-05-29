package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`

	Auth struct {
		Id        int64  `mapstructure:"id"`
		User      string `mapstructure:"user"`
		Password  string `mapstructure:"password"`
		SignKey   string `mapstructure:"sign_key"`
		LocalPath string `mapstructure:"local_path"`
		Issuer    string `mapstructure:"issuer"`
	} `mapstructure:"auth"`

	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}

var Global *Config

func InitConfig() {
	defaultPath := "./internal/config"
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load(defaultPath + "/.env")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultPath)

	// 支持通过环境变量覆盖配置
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	Global = &c

	log.Println("配置初始化完成")
}
