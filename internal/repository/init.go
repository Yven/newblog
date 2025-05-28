package repository

import "database/sql"

func InitTable(db *sql.DB) error {
	q := `
CREATE TABLE IF NOT EXISTS article(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	cid INTEGER DEFAULT 0,
	create_time TEXT NOT NULL DEFAULT (datetime('now')),
	update_time TEXT NOT NULL DEFAULT (datetime('now')),
	delete_time TEXT
);
CREATE TABLE IF NOT EXISTS category(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS article_tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	aid INTEGER NOT NULL,
	tid INTEGER NOT NULL
);
`

	_, err := db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
