package main

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
