package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const usersSchema = `
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

// initDB removes any existing database file at path and creates a fresh
// one with the users table, so each container start gets a clean schema.
func initDB(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating db directory: %w", err)
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("removing existing db file: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	if _, err := db.Exec(usersSchema); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating users table: %w", err)
	}

	return db, nil
}
