package repository

import (
	"database/sql"
	"log"
	"newblog/internal/global"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func New() *sql.DB {
	// Reuse Connection
	if global.DbInstance != nil {
		return global.DbInstance
	}

	db, err := sql.Open("sqlite3", global.DbURL)
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
