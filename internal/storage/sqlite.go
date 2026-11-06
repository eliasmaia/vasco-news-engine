package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	dbUrl := fmt.Sprintf("%s?authToken=%s", url, token)

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS sent_news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			link TEXT UNIQUE,
			title TEXT,
			source TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	return &DB{conn: db}, err
}

func (d *DB) IsNew(link string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM sent_news WHERE link = ?)"
	err := d.conn.QueryRow(query, link).Scan(&exists)
	if err != nil {
		return false
	}
	return !exists
}

func (d *DB) Save(n string, title string, source string) error {
	query := "INSERT OR IGNORE INTO sent_news (link, title, source) VALUES (?, ?, ?)"
	_, err := d.conn.Exec(query, n, title, source)
	return err
}
