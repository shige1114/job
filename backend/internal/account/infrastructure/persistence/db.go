package persistence

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// テーブルの初期化
	schema := `
	CREATE TABLE IF NOT EXISTS applications (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		code TEXT NOT NULL,
		password TEXT,
		expires_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_applications_email ON applications(email);

	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}
