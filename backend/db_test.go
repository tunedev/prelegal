package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDB_CreatesUsersTable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	db, err := initDB(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users'`).Scan(&name)
	if err != nil {
		t.Fatalf("expected users table to exist, got error: %v", err)
	}
	if name != "users" {
		t.Errorf("expected table name 'users', got %q", name)
	}
}

func TestInitDB_RecreatesFromScratch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	db, err := initDB(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO users (email, password_hash) VALUES (?, ?)`, "a@example.com", "hash"); err != nil {
		t.Fatalf("unexpected error inserting row: %v", err)
	}
	db.Close()

	db2, err := initDB(path)
	if err != nil {
		t.Fatalf("unexpected error on second init: %v", err)
	}
	defer db2.Close()

	var count int
	if err := db2.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		t.Fatalf("unexpected error counting rows: %v", err)
	}
	if count != 0 {
		t.Errorf("expected users table to be empty after re-init, got %d rows", count)
	}
}

func TestInitDB_CreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "data", "app.db")

	db, err := initDB(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("expected valid db after init, ping failed: %v", err)
	}
}

func TestInitDB_RemovesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	if err := os.WriteFile(path, []byte("not a real sqlite file"), 0o644); err != nil {
		t.Fatalf("unexpected error writing stub file: %v", err)
	}

	db, err := initDB(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("expected valid db after init, ping failed: %v", err)
	}
}
