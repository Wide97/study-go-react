package main

import (
	"database/sql"
	"net/http"
)

func startScheduler(db *sql.DB, services []Service) {

}

func checkService(client *http.Client, service Service) (status string, responseTimeMs int) {

}
