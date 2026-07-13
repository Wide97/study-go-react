package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Order struct {
	ID       int     `json:"id"`
	Customer string  `json:"customer"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
}

type OrdersResponse struct {
	Items []Order `json:"items"`
	Total int     `json:"total"`
}

var orders = []Order{
	{ID: 1, Customer: "Alice", Status: "pending", Total: 150.00},
	{ID: 2, Customer: "Bob", Status: "shipped", Total: 200.00},
	{ID: 3, Customer: "Charlie", Status: "delivered", Total: 300.00},
	{ID: 4, Customer: "David", Status: "pending", Total: 120.00},
	{ID: 5, Customer: "Eve", Status: "shipped", Total: 250.00},
	{ID: 6, Customer: "Frank", Status: "delivered", Total: 400.00},
	{ID: 7, Customer: "Grace", Status: "pending", Total: 180.00},
	{ID: 8, Customer: "Heidi", Status: "shipped", Total: 220.00},
	{ID: 9, Customer: "Ivan", Status: "delivered", Total: 350.00},
	{ID: 10, Customer: "Judy", Status: "pending", Total: 100.00},
}

// frontendOrigin è l'origin del dev server Vite.
// Origin = protocollo + host + porta. Quindi localhost:5173 e localhost:8080
// sono origin diversi anche se girano sulla stessa macchina.
const frontendOrigin = "http://localhost:5173"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /orders", ordersHandler)

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

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	filteredOrders := orders
	filteredOrders = filterOrdersByStatus(filteredOrders, status)
	filteredOrders = filterOrdersBySearch(filteredOrders, search)

	json.NewEncoder(w).Encode(OrdersResponse{
		Items: filteredOrders,
		Total: len(filteredOrders),
	})
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

func filterOrdersByStatus(orders []Order, status string) []Order {
	if status == "" {
		return orders
	}

	var filtered []Order
	for _, order := range orders {
		if order.Status == status {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func filterOrdersBySearch(orders []Order, search string) []Order {
	if search == "" {
		return orders
	}

	var filtered []Order
	for _, order := range orders {
		if containsIgnoreCase(order.Customer, search) {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func containsIgnoreCase(str, substr string) bool {
	str = strings.ToLower(str)
	substr = strings.ToLower(substr)
	return strings.Contains(str, substr)
}
