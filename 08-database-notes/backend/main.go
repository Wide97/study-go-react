package main

import (
	"fmt"
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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// health è l'endpoint minimo per verificare che il backend sia acceso.
// Prima di collegare il database, ci serve una base HTTP funzionante.
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}
