package repository

import (
	"database/sql"
	"os"
)

func InitTable(db *sql.DB) error {
	content, err := os.ReadFile("./db/init.sql")
	if err != nil {
		return err
	}

	q := string(content)

	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
