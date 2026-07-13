package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

// createNoteHandler gestisce POST /notes.
//
// Qui siamo nello strato HTTP:
// - leggiamo il JSON ricevuto nel body;
// - lo convertiamo in NoteRequest;
// - facciamo una validazione minima;
// - deleghiamo al repository la scrittura su database.
func createNoteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NoteRequest

		// Il body della richiesta contiene JSON.
		// Decode legge quel JSON e riempie req.
		// Passiamo &req perché il decoder deve modificare la variabile originale.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Il database ha title/content NOT NULL, ma conviene validare prima
		// nell'handler: così rispondiamo con 400, cioè errore del client.
		if req.Title == "" || req.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		note, err := createNote(db, req)
		if err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// 201 Created è lo status corretto quando una POST crea una nuova risorsa.
		// WriteHeader va chiamato prima di scrivere il JSON nel body.
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

// updateNoteHandler gestisce PUT /notes/{id}.
//
// Questo handler è il ponte tra HTTP e repository:
// - legge l'id dalla URL;
// - legge il JSON dal body;
// - valida i dati ricevuti;
// - chiama updateNote, che contiene la logica SQL.
func updateNoteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Con il pattern "PUT /notes/{id}", net/http salva il pezzo dinamico
		// della URL dentro r.PathValue("id").
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Invalid note id", http.StatusBadRequest)
			return
		}

		var req NoteRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		note, err := updateNote(db, id, req)
		if err != nil {
			// sql.ErrNoRows è il segnale standard che la riga richiesta non esiste.
			// In HTTP lo traduciamo in 404 Not Found.
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Note not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to update note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}

// health è l'endpoint minimo per verificare che il backend sia acceso.
// Prima di collegare il database, ci serve una base HTTP funzionante.
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}
