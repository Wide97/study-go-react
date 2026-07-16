package main

import (
	// Blank import: non usiamo direttamente il package, ma il suo init()
	// registra il driver "sqlite" dentro database/sql.
	"database/sql"

	_ "modernc.org/sqlite"
)

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "uptime.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS services (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			interval_seconds INTEGER NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil

}
