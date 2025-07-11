package repository

import (
	"context"
	"database/sql"
	"log"
	"newblog/internal/model"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitSQLite(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// 数据表初始化
	err = InitTable(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitNeo4j(conf model.Neo4jConfig) *neo4j.DriverWithContext {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		conf.Host,
		neo4j.BasicAuth(conf.User, conf.Password, ""),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &driver
}
