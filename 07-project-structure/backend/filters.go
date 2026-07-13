package main

import "strings"

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
