package main

import "database/sql"

// listNotes legge tutte le note dal database.
//
// Il repository è lo strato che conosce SQL.
// Gli handler HTTP useranno questa funzione invece di scrivere query direttamente.
func listNotes(db *sql.DB) ([]Note, error) {
	rows, err := db.Query(`
		SELECT id, title, content, created_at, updated_at
		FROM notes
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []Note{}

	// rows.Next avanza una riga alla volta.
	// Quando non ci sono più righe, il ciclo finisce.
	for rows.Next() {
		var note Note

		// Scan copia le colonne della riga corrente dentro i campi Go.
		// L'ordine deve combaciare con quello della SELECT.
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	// rows.Err controlla errori avvenuti durante l'iterazione.
	// È separato dagli errori di Scan.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
