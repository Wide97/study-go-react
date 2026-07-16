package main

import (
	"database/sql"
	"time"
)

func listServices(db *sql.DB) ([]Service, error) {
	rows, err := db.Query(`
		SELECT id, name, url, interval_seconds, created_at
		FROM services
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	// defer perchè altrimenti resterebbe aperta la connessione a db, serve per chiuderla.
	defer rows.Close()

	services := []Service{}

	for rows.Next() {
		var service Service

		// & prende l'indirizzo di ogni campo di "service": Scan scrive il valore
		// letto dalla colonna direttamente lì dentro, non su una copia.
		// Vedi repository.md per il dettaglio su &.
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.URL,
			&service.IntervalSeconds,
			&service.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, service)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil

}

func createService(db *sql.DB, req ServiceRequest) (Service, error) {

	now := time.Now().Format(time.RFC3339)

	result, err := db.Exec(`
		INSERT INTO services (name, url, interval_seconds, created_at)
		VALUES (?, ?, ?, ?)
		`, req.Name, req.URL, req.IntervalSeconds, now)
	//quanto sopra sono i parametri che passo agli ?
	if err != nil {
		return Service{}, err
	}

	//LastInsertId mi restituisce il valore dell' id messo in tabella
	id, err := result.LastInsertId()
	if err != nil {
		return Service{}, err
	}

	service := Service{
		ID:              int(id),
		Name:            req.Name,
		URL:             req.URL,
		IntervalSeconds: req.IntervalSeconds,
		CreatedAt:       now,
	}

	return service, nil
}
