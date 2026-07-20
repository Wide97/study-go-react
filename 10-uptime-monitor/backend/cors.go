package main

import "net/http"

// frontendOrigin è l'origin del dev server Vite.
//
// Origin significa: protocollo + host + porta.
// Quindi:
// - http://localhost:5173 è l'origin del frontend Vite;
// - http://localhost:8080 è l'origin del backend Go.
//
// Per il browser sono due origin diversi, anche se girano entrambi sulla
// stessa macchina. Per questo serve CORS.
const frontendOrigin = "http://localhost:5173"

// withCORS è un middleware.
//
// Un middleware prende un handler e restituisce un altro handler.
// In pratica avvolge il router principale e aggiunge un comportamento comune
// a tutte le richieste: in questo caso gli header CORS.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Access-Control-Allow-Origin dice al browser quale frontend può leggere
		// le risposte del backend.
		w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)

		// Access-Control-Allow-Methods elenca i metodi HTTP permessi quando la
		// richiesta arriva dal frontend.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Access-Control-Allow-Headers elenca gli header che il frontend può
		// mandare. Content-Type serve perché POST e PUT inviano JSON.
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Prima di certe richieste "non semplici", il browser manda una richiesta
		// OPTIONS detta preflight. Non è la vera richiesta applicativa: serve solo
		// a chiedere al backend se il browser è autorizzato a procedere.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
