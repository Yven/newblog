package global

import (
	"database/sql"
	"newblog/internal/util"
)

var (
	DbInstance *sql.DB
	JwtService util.JwtService
)
