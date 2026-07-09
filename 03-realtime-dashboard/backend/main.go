package main

import (
	// fmt serve qui solo per scrivere "Ok" nella risposta di /health.
	"fmt"

	// log serve per stampare messaggi sul terminale e per chiudere il programma
	// se il server non riesce ad avviarsi.
	"log"

	// net/http è il pacchetto standard per creare server HTTP, router e handler.
	"net/http"

	// time serve per creare un ticker: un timer che produce un evento a intervalli
	// regolari. Lo useremo per inviare un dato ogni secondo.
	"time"

	// gorilla/websocket aggiunge il supporto WebSocket sopra net/http.
	// net/http gestisce la richiesta HTTP iniziale; Gorilla fa l'upgrade
	// della connessione e offre metodi comodi come WriteJSON.
	"github.com/gorilla/websocket"
)

// SensorMessage è il messaggio JSON che il backend invia al frontend.
//
// Il tag json:"value" dice che il campo Go Value diventerà "value" nel JSON:
// SensorMessage{Value: 3} -> {"value":3}
type SensorMessage struct {
	Value int `json:"value"`
}

// Upgrader è l'oggetto che trasforma una richiesta HTTP normale in una
// connessione WebSocket.
//
// Il client arriva con una richiesta GET /ws che contiene header speciali
// ("Upgrade: websocket"). Se tutto è valido, Upgrade risponde con uno status
// 101 Switching Protocols e da quel momento la connessione resta aperta.
var upgrader = websocket.Upgrader{}

func main() {
	// ServeMux è il router: associa metodo+path agli handler.
	mux := http.NewServeMux()

	// Endpoint HTTP classico, utile per verificare che il server sia acceso.
	mux.HandleFunc("GET /health", health)

	// Endpoint WebSocket: il client React si collegherà a ws://localhost:8080/ws.
	mux.HandleFunc("GET /ws", websocketHandler)

	log.Println("Server running on :8080")

	// ListenAndServe resta in ascolto sulla porta 8080.
	// log.Fatal stampa eventuali errori fatali, per esempio porta già occupata.
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Un handler ha sempre questa firma: riceve dove scrivere la risposta (w)
// e i dati della richiesta in arrivo (r).
func health(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter è "la penna" con cui scriviamo la risposta HTTP.
	// Qui non impostiamo Content-Type perché è solo testo semplice.
	fmt.Fprint(w, "Ok")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade prova a trasformare la richiesta HTTP in una connessione WebSocket.
	// Dopo l'upgrade non usiamo più w per scrivere risposte HTTP classiche:
	// usiamo conn per leggere/scrivere messaggi WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	// Chiudiamo la connessione quando l'handler termina.
	// L'handler terminerà quando WriteJSON fallisce, per esempio perché
	// il client ha chiuso la pagina o perso la connessione.
	defer conn.Close()

	// NewTicker crea un evento periodico: qui uno ogni secondo.
	// ticker.C è un canale: riceve un valore ad ogni tick.
	ticker := time.NewTicker(time.Second)

	// Fermare il ticker evita di lasciare risorse attive quando l'handler finisce.
	defer ticker.Stop()

	// Contatore finto: per ora simula un dato che cambia nel tempo.
	count := 0

	// Loop infinito: resta attivo finché la connessione WebSocket è valida.
	for {
		// <-ticker.C blocca il ciclo finché non arriva il prossimo tick.
		// In pratica: aspetta un secondo.
		<-ticker.C

		count++

		message := SensorMessage{
			Value: count,
		}

		// WriteJSON serializza la struct in JSON e la invia sul WebSocket.
		// Il frontend riceverà messaggi come {"value":1}, {"value":2}, ...
		err := conn.WriteJSON(message)
		if err != nil {
			// Se il client chiude la connessione, la prossima WriteJSON fallisce.
			// A quel punto usciamo dall'handler: defer chiuderà conn e ticker.
			log.Println("write error:", err)
			return
		}
	}
}
