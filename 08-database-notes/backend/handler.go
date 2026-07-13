package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func notesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := listNotes(db)
		if err != nil {
			http.Error(w, "Failed to list notes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}

// health è l'endpoint minimo per verificare che il backend sia acceso.
// Prima di collegare il database, ci serve una base HTTP funzionante.
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}
