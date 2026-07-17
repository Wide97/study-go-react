package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

// una funzione anonima definita dentro un'altra e che "cattura" le variabili
// della funzione esterna si chiama closure: qui la funzione restituita si ricorda di "db"
func servicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := listServices(db)

		if err != nil {
			http.Error(w, "Failed to load Services", http.StatusInternalServerError)
			return
		}
		//dice al fe che quello che riceve è un formato json
		w.Header().Set("Content-Type", "application/json")
		//converte la struct di go in un json
		json.NewEncoder(w).Encode(services)
	}

}

func createServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ServiceRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.URL == "" || req.IntervalSeconds <= 0 {
			http.Error(w, "Name and URL are required. Interval seconds must be grater than 0.", http.StatusBadRequest)
			return
		}

		service, err := createService(db, req)
		if err != nil {
			http.Error(w, "Failed to create Services", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		//scrivo lo stato di creazione corretto (201)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(service)

	}
}
