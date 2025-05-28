package global

import (
	"database/sql"
	"os"
	"strconv"
)

type AdminInfo struct {
	User     string
	Password string
}

type TokenInfo struct {
	Key  string
	Path string
}

var (
	Port, _ = strconv.Atoi(os.Getenv("PORT"))

	Admin = &AdminInfo{
		User:     os.Getenv("ADMIN_NAME"),
		Password: os.Getenv("PASSWORD"),
	}
	Token = &TokenInfo{
		Key:  os.Getenv("JWT_SIGN_KEY"),
		Path: os.Getenv("LOCAL_TOKEN_PATH"),
	}

	DbURL      = os.Getenv("BLUEPRINT_DB_URL")
	DbInstance *sql.DB
)
