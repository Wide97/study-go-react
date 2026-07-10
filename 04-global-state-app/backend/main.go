package main

import (
	// encoding/json converte struct e slice Go in JSON.
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Product rappresenta un prodotto disponibile nel catalogo.
// I tag json controllano i nomi dei campi nella risposta HTTP.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// products è uno storage in memoria.
// Non c'è database: i dati sono fissi e vivono finché il server è acceso.
var products = []Product{
	{ID: 1, Name: "Keyboard", Price: 79.90},
	{ID: 2, Name: "Mouse", Price: 39.90},
	{ID: 3, Name: "Monitor", Price: 249.90},
}

// frontendOrigin è l'origin del dev server Vite.
// Origin = protocollo + host + porta. Quindi localhost:5173 e localhost:8080
// sono origin diversi anche se girano sulla stessa macchina.
const frontendOrigin = "http://localhost:5173"

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)

	mux.HandleFunc("GET /products", productsHandler)

	log.Println("Server running on :8080")

	// ListenAndServe resta in ascolto sulla porta 8080.
	// log.Fatal stampa eventuali errori fatali, per esempio porta già occupata.
	// withCORS avvolge tutto il router, così eventuali richieste HTTP dal
	// frontend Vite ricevono gli header corretti.
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))

}

// Un handler ha sempre questa firma: riceve dove scrivere la risposta (w)
// e i dati della richiesta in arrivo (r).
func health(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter è "la penna" con cui scriviamo la risposta HTTP.
	// Qui non impostiamo Content-Type perché è solo testo semplice.
	fmt.Fprint(w, "Ok")

}

// withCORS permette al frontend Vite di chiamare endpoint HTTP del backend.
// Serve perché frontend (:5173) e backend (:8080) hanno origin diversi.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// productsHandler restituisce la lista prodotti in formato JSON.
func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
