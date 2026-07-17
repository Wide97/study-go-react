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
