package main

import (
	// log serve per stampare messaggi sul terminale.
	// log.Fatal stampa l'errore e termina il programma se il server non parte.
	"log"

	// net/http è il pacchetto standard Go per creare server HTTP.
	// Qui usiamo router, handler, request, response e status code.
	"net/http"
)

// frontendOrigin è l'origin del dev server Vite.
// Origin = protocollo + host + porta. Quindi localhost:5173 e localhost:8080
// sono origin diversi anche se girano sulla stessa macchina.
const frontendOrigin = "http://localhost:5173"

func main() {
	// ServeMux è il router standard di net/http.
	// Associa metodo+path a una funzione handler.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)

	// GET /orders è l'endpoint principale del progetto 05.
	// Supporta query params come:
	// /orders?status=pending&search=ali&page=1&pageSize=5
	mux.HandleFunc("GET /orders", ordersHandler)

	log.Println("Server running on :8080")

	// ListenAndServe resta in ascolto sulla porta 8080.
	// log.Fatal stampa eventuali errori fatali, per esempio porta già occupata.
	// withCORS avvolge tutto il router, così eventuali richieste HTTP dal
	// frontend Vite ricevono gli header corretti.
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
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
