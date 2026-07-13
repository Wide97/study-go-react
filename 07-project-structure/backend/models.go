package main

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
