package main

import (
	"log"
	"net/http"
)

func main() {

	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /notes", notesHandler(db))
	mux.HandleFunc("POST /notes", createNoteHandler(db))
	mux.HandleFunc("PUT /notes/{id}", updateNoteHandler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
