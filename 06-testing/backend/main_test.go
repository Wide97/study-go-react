package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// In Go i file che finiscono con _test.go vengono compilati solo durante i test.
// Le funzioni di test devono iniziare con Test e ricevere *testing.T.

func TestFilterOrdersByStatus(t *testing.T) {
	// Arrange: prepariamo i dati di partenza.
	// Nei test conviene usare dati piccoli, così è immediato capire cosa ci
	// aspettiamo.
	input := []Order{
		{ID: 1, Customer: "Alice", Status: "pending", Total: 150},
		{ID: 2, Customer: "Bob", Status: "shipped", Total: 200},
		{ID: 3, Customer: "David", Status: "pending", Total: 120},
	}

	// Act: chiamiamo la funzione che vogliamo verificare.
	result := filterOrdersByStatus(input, "pending")

	// Assert: controlliamo il comportamento osservabile.
	// Qui non ci interessa come la funzione filtra internamente, ci interessa
	// che il risultato contenga solo gli ordini pending.
	if len(result) != 2 {
		t.Fatalf("expected 2 pending orders, got %d", len(result))
	}

	for _, order := range result {
		if order.Status != "pending" {
			t.Fatalf("expected only pending orders, got status %q", order.Status)
		}
	}
}

func TestFilterOrdersByStatusReturnsAllWhenStatusIsEmpty(t *testing.T) {
	input := []Order{
		{ID: 1, Customer: "Alice", Status: "pending", Total: 150},
		{ID: 2, Customer: "Bob", Status: "shipped", Total: 200},
	}

	result := filterOrdersByStatus(input, "")

	if len(result) != len(input) {
		t.Fatalf("expected all orders when status is empty, got %d", len(result))
	}
}

func TestFilterOrdersBySearchIsCaseInsensitive(t *testing.T) {
	input := []Order{
		{ID: 1, Customer: "Alice", Status: "pending", Total: 150},
		{ID: 2, Customer: "Bob", Status: "shipped", Total: 200},
	}

	// Cerchiamo "ALI" maiuscolo, ma l'ordine contiene "Alice".
	// Il test documenta che la ricerca non deve distinguere maiuscole/minuscole.
	result := filterOrdersBySearch(input, "ALI")

	if len(result) != 1 {
		t.Fatalf("expected 1 matching order, got %d", len(result))
	}

	if result[0].Customer != "Alice" {
		t.Fatalf("expected Alice, got %q", result[0].Customer)
	}
}

func TestPaginateOrders(t *testing.T) {
	input := []Order{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
	}

	// page=2 e pageSize=2 devono prendere gli elementi con indice 2 e 3,
	// quindi gli ordini con ID 3 e 4.
	result := paginateOrders(input, 2, 2)

	if len(result) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result))
	}

	if result[0].ID != 3 || result[1].ID != 4 {
		t.Fatalf("expected order IDs 3 and 4, got %d and %d", result[0].ID, result[1].ID)
	}
}

func TestPaginateOrdersReturnsEmptySliceForPageOutOfRange(t *testing.T) {
	input := []Order{
		{ID: 1}, {ID: 2}, {ID: 3},
	}

	result := paginateOrders(input, 99, 2)

	if len(result) != 0 {
		t.Fatalf("expected empty result for page out of range, got %d orders", len(result))
	}
}

func TestGetIntQueryParam(t *testing.T) {
	// httptest.NewRequest crea una request finta senza avviare un server reale.
	// È utile per testare funzioni che leggono dati da *http.Request.
	request := httptest.NewRequest(http.MethodGet, "/orders?page=2&pageSize=3", nil)

	page := getIntQueryParam(request, "page", 1)
	pageSize := getIntQueryParam(request, "pageSize", 5)

	if page != 2 {
		t.Fatalf("expected page 2, got %d", page)
	}

	if pageSize != 3 {
		t.Fatalf("expected pageSize 3, got %d", pageSize)
	}
}

func TestGetIntQueryParamFallsBackToDefault(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/orders?page=abc&pageSize=-1", nil)

	page := getIntQueryParam(request, "page", 1)
	pageSize := getIntQueryParam(request, "pageSize", 5)

	if page != 1 {
		t.Fatalf("expected default page 1, got %d", page)
	}

	if pageSize != 5 {
		t.Fatalf("expected default pageSize 5, got %d", pageSize)
	}
}

func TestOrdersHandlerReturnsFilteredAndPaginatedResponse(t *testing.T) {
	// Questo è un test di handler HTTP.
	// Non avviamo il server su :8080: chiamiamo direttamente la funzione handler
	// usando request/recorder finti.
	request := httptest.NewRequest(
		http.MethodGet,
		"/orders?status=pending&search=a&page=1&pageSize=2",
		nil,
	)
	responseRecorder := httptest.NewRecorder()

	ordersHandler(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", responseRecorder.Code)
	}

	var response OrdersResponse
	err := json.NewDecoder(responseRecorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	// Gli ordini pending con "a" nel nome sono Alice, David e Grace.
	// pageSize=2 significa: Items contiene i primi 2, Total resta 3.
	if response.Total != 3 {
		t.Fatalf("expected total 3, got %d", response.Total)
	}

	if len(response.Items) != 2 {
		t.Fatalf("expected 2 paginated items, got %d", len(response.Items))
	}

	if response.Page != 1 || response.PageSize != 2 {
		t.Fatalf("expected page=1 and pageSize=2, got page=%d pageSize=%d", response.Page, response.PageSize)
	}
}
