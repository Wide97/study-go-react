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
	mux.HandleFunc("GET /services", servicesHandler(db))
	mux.HandleFunc("POST /services", createServiceHandler(db))
	mux.HandleFunc("PUT /services/{id}", updateServiceHandler(db))
	mux.HandleFunc("DELETE /services/{id}", deleteServiceHandler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))

}
