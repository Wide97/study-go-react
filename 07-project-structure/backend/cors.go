package main

import "net/http"

// frontendOrigin è l'origin del dev server Vite.
// Origin = protocollo + host + porta. Quindi localhost:5173 e localhost:8080
// sono origin diversi anche se girano sulla stessa macchina.
const frontendOrigin = "http://localhost:5173"

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
