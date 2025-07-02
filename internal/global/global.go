package global

import (
	"database/sql"
	"log/slog"
	"newblog/internal/util"
)

var (
	DbInstance *sql.DB
	JwtService util.JwtService
	Logger     *slog.Logger
	Visitors   *util.Visitors
)

func Init(db *sql.DB, jwt util.JwtService, logger *slog.Logger, visitors *util.Visitors) {
	DbInstance = db
	JwtService = jwt
	Logger = logger
	Visitors = visitors
}
