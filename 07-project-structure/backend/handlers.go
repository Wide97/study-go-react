package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Un handler ha sempre questa firma: riceve dove scrivere la risposta (w)
// e i dati della richiesta in arrivo (r).
func health(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter è "la penna" con cui scriviamo la risposta HTTP.
	// Qui non impostiamo Content-Type perché è solo testo semplice.
	fmt.Fprint(w, "Ok")
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	// Dichiariamo subito che la risposta sarà JSON.
	// Gli header vanno impostati prima di scrivere il body.
	w.Header().Set("Content-Type", "application/json")

	// Query() legge i parametri dopo il ? nell'URL.
	// Esempio: /orders?status=pending&search=ali
	// status vale "pending", search vale "ali".
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	// page e pageSize arrivano come stringhe.
	// getIntQueryParam incapsula conversione, default e validazione.
	page := getIntQueryParam(r, "page", 1)
	pageSize := getIntQueryParam(r, "pageSize", 5)

	// Pipeline dei filtri:
	// partiamo dalla lista completa, poi ogni funzione riceve una lista e
	// restituisce una nuova lista filtrata.
	//
	// Questo rende semplice combinare i filtri:
	// prima status, poi search, poi paginazione.
	filteredOrders := orders
	filteredOrders = filterOrdersByStatus(filteredOrders, status)
	filteredOrders = filterOrdersBySearch(filteredOrders, search)

	// total va calcolato DOPO i filtri ma PRIMA della paginazione.
	// Se ci sono 10 risultati filtrati e pageSize=5, il frontend deve sapere
	// che total è 10 anche se Items contiene solo 5 elementi.
	total := len(filteredOrders)

	// paginateOrders taglia la lista filtrata e restituisce solo gli elementi
	// della pagina richiesta.
	paginatedOrders := paginateOrders(filteredOrders, page, pageSize)

	json.NewEncoder(w).Encode(OrdersResponse{
		Items:    paginatedOrders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
