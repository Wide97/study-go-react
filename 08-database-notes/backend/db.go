package main

import (
	"database/sql"

	// Blank import: non usiamo direttamente il package, ma il suo init()
	// registra il driver "sqlite" dentro database/sql.
	_ "modernc.org/sqlite"
)

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "notes.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
