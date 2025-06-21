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
)
