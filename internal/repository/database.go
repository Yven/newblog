package repository

import (
	"database/sql"
	"log"
	"newblog/internal/config"
	"newblog/internal/global"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() *sql.DB {
	if global.DbInstance != nil {
		return global.DbInstance
	}

	db, err := sql.Open("sqlite3", config.Global.Database.Host)
	if err != nil {
		log.Fatal(err)
	}

	global.DbInstance = db

	// 数据表初始化
	err = InitTable(db)
	if err != nil {
		log.Fatal(err)
	}

	return global.DbInstance
}
