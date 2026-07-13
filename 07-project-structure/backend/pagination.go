package main

import (
	"net/http"
	"strconv"
)

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
