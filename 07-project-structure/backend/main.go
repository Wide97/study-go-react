package main

import (
	// log serve per stampare messaggi sul terminale.
	// log.Fatal stampa l'errore e termina il programma se il server non parte.
	"log"

	// net/http è il pacchetto standard Go per creare server HTTP.
	// Qui usiamo router, handler, request, response e status code.
	"net/http"

	// strconv contiene conversioni tra stringhe e tipi primitivi.
	// I query param arrivano sempre come stringhe, quindi page/pageSize vanno
	// convertiti in int con strconv.Atoi.
	"strconv"
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

// getIntQueryParam legge un query param numerico.
// Serve perché r.URL.Query().Get(...) restituisce sempre stringhe.
// Se il parametro manca, non è un numero, oppure è <= 0, usiamo defaultValue.
func getIntQueryParam(r *http.Request, name string, defaultValue int) int {
	value := r.URL.Query().Get(name)

	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return defaultValue
	}

	return number
}

// paginateOrders restituisce solo la pagina richiesta.
//
// In Go uno slice si taglia con orders[start:end].
// start è incluso, end è escluso.
//
// Esempio: orders[0:3] prende gli elementi con indice 0, 1, 2.
func paginateOrders(orders []Order, page int, pageSize int) []Order {
	// page parte da 1 per il client, ma gli indici degli slice Go partono da 0.
	// Per questo page=1 produce start=0.
	start := (page - 1) * pageSize
	end := start + pageSize

	// Se start supera la lunghezza della lista, la pagina richiesta non esiste.
	// Ritorniamo una slice vuota invece di causare un panic per indice fuori range.
	if start >= len(orders) {
		return []Order{}
	}

	// Nell'ultima pagina end può andare oltre la lunghezza della lista.
	// Lo riportiamo al massimo valido: len(orders).
	if end > len(orders) {
		end = len(orders)
	}

	return orders[start:end]
}
