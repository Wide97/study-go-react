package main

import (
	// encoding/json serve per trasformare struct Go in JSON.
	// Qui lo usiamo per rispondere a GET /orders con una lista strutturata.
	"encoding/json"

	// fmt contiene funzioni per scrivere/formattare testo.
	// Lo usiamo solo nell'endpoint /health per scrivere "Ok".
	"fmt"

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

	// strings contiene funzioni per lavorare con stringhe.
	// Qui lo usiamo per fare ricerca case-insensitive sul nome cliente.
	"strings"
)

// Order rappresenta un ordine restituito dal backend.
// I tag json dicono a encoding/json i nomi dei campi nella risposta HTTP.
type Order struct {
	ID       int     `json:"id"`
	Customer string  `json:"customer"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
}

// OrdersResponse è la forma della risposta di GET /orders.
//
// Items contiene solo la pagina corrente.
// Total contiene il numero totale di risultati dopo i filtri, ma prima della
// paginazione: al frontend serve per capire quante pagine esistono.
// Page e PageSize vengono rimandati al client per rendere esplicita la pagina
// che il backend ha usato.
type OrdersResponse struct {
	Items    []Order `json:"items"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

// orders è un dataset finto in memoria.
// Per ora non usiamo database: ci basta una slice per studiare query params,
// filtri, ricerca e paginazione.
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
	// Se il client non manda status, il filtro non deve fare nulla.
	// Ritorniamo direttamente la lista ricevuta.
	if status == "" {
		return orders
	}

	// Creiamo una nuova slice invece di modificare quella originale.
	// append aggiunge solo gli ordini che passano il controllo.
	var filtered []Order
	for _, order := range orders {
		if order.Status == status {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func filterOrdersBySearch(orders []Order, search string) []Order {
	// Se search è vuoto, non filtriamo.
	// Nota: strings.Contains("Alice", "") sarebbe true, ma qui rendiamo
	// esplicita l'intenzione.
	if search == "" {
		return orders
	}

	var filtered []Order
	for _, order := range orders {
		// Cerchiamo solo nel nome cliente.
		// containsIgnoreCase evita differenze tra "Alice", "alice" e "ALI".
		if containsIgnoreCase(order.Customer, search) {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func containsIgnoreCase(str, substr string) bool {
	// Per fare una ricerca case-insensitive normalizziamo entrambe le stringhe:
	// tutto minuscolo, poi usiamo strings.Contains.
	str = strings.ToLower(str)
	substr = strings.ToLower(substr)
	return strings.Contains(str, substr)
}

// getIntQueryParam legge un query param numerico.
//
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
