package main

import (
	// encoding/json serve per trasformare JSON <-> struct Go.
	// Nel login lo usiamo in ingresso: leggiamo il body della richiesta HTTP,
	// che arriva come JSON, e lo decodifichiamo dentro LoginRequest.
	"encoding/json"

	// fmt contiene funzioni per scrivere/formattare testo.
	// Qui lo usiamo solo per rispondere "Ok" nell'endpoint /health.
	"fmt"

	// log serve per stampare messaggi sul terminale.
	// log.Fatal è utile con ListenAndServe perché stampa l'errore e termina
	// il programma se il server non riesce ad avviarsi.
	"log"

	// net/http è il pacchetto standard Go per server e client HTTP.
	// Qui usiamo server, router, handler, status code e funzioni di errore.
	"net/http"

	// bcrypt è un algoritmo pensato per salvare password in modo sicuro.
	// Non si confronta mai la password in chiaro salvata da qualche parte:
	// si salva un hash, poi al login si confronta la password ricevuta con
	// quell'hash tramite CompareHashAndPassword.
	"golang.org/x/crypto/bcrypt"
)

// demoUser è un utente finto tenuto in memoria.
// Per ora non abbiamo database, registrazione o più utenti: ci basta un solo
// utente per imparare il flusso di login.
//
// Nota importante: la password NON è salvata come "password123".
// PasswordHash contiene l'hash bcrypt di "password123".
// Questo significa che chi legge questa variabile non vede la password reale,
// ma il backend può comunque verificare se una password inserita combacia.
var demoUser = User{
	Email:        "demo@example.com",
	PasswordHash: "$2a$10$rSUOmWsnmo536ewsL1Si1O5qJUTHqGgxMY/cUwhmiW2FQWD02UC3u",
}

// User rappresenta l'utente interno al backend.
// Non ha tag json perché non viene letto direttamente da una richiesta
// e non viene mandato al frontend. È solo un modello interno.
type User struct {
	Email string

	// PasswordHash è una stringa generata da bcrypt.
	// Un hash bcrypt include dentro di sé anche il "salt" e il costo usato:
	// per questo la stessa password può produrre hash diversi ma restare
	// verificabile con CompareHashAndPassword.
	PasswordHash string
}

// LoginRequest descrive il JSON che il client deve mandare a POST /login.
//
// Esempio:
//
//	{
//	  "email": "demo@example.com",
//	  "password": "password123"
//	}
//
// I tag `json:"email"` e `json:"password"` dicono a encoding/json come
// collegare i campi JSON, minuscoli, ai campi Go, che devono iniziare con
// maiuscola per essere esportati e quindi valorizzabili dal decoder.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse descrive il JSON che il backend manda quando il login riesce.
// Per ora contiene solo un messaggio. Nel prossimo step lo sostituiremo con
// un token JWT, perché il frontend dovrà conservarlo e rimandarlo nelle
// richieste protette.
type LoginResponse struct {
	Message string `json:"message"`
}

func main() {
	// ServeMux è il router standard di net/http.
	// Il suo compito è associare una coppia metodo+path a una funzione handler.
	mux := http.NewServeMux()

	// Da Go 1.22 il pattern può includere anche il metodo HTTP.
	// "GET /health" significa: questa funzione risponde solo a GET /health.
	// Se provi POST /health, il router può rispondere 405 Method Not Allowed.
	mux.HandleFunc("GET /health", health)

	// POST /login riceve credenziali nel body.
	// Usiamo POST, non GET, perché stiamo inviando dati sensibili e perché il
	// login è un'azione, non una semplice lettura di risorsa.
	mux.HandleFunc("POST /login", login)

	log.Println("Server running on :8080")

	// ListenAndServe mette il server in ascolto sulla porta 8080.
	// È una chiamata bloccante: il programma resta qui finché il server gira.
	// Se qualcosa va storto, per esempio porta occupata, restituisce un errore.
	// log.Fatal stampa quell'errore e chiude il programma.
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Un handler ha sempre questa firma: riceve dove scrivere la risposta (w)
// e i dati della richiesta in arrivo (r).
func health(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter è "la penna" con cui scriviamo la risposta HTTP.
	// Qui non impostiamo Content-Type perché è solo testo semplice.
	fmt.Fprint(w, "Ok")
}

// login gestisce POST /login.
// Il suo lavoro è:
// 1. leggere il JSON mandato dal client;
// 2. verificare che l'email esista;
// 3. verificare che la password combaci con l'hash bcrypt;
// 4. rispondere JSON se il login è valido.
func login(w http.ResponseWriter, r *http.Request) {
	// req partirà vuota. Decode la riempirà leggendo il body della richiesta.
	// Serve passare &req, cioè l'indirizzo della variabile, perché Decode deve
	// poter modificare il valore originale.
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Se il JSON è malformato o non leggibile, l'errore è del client.
		// 400 Bad Request significa: la richiesta non ha una forma valida.
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Primo controllo: email.
	// Rispondiamo con un messaggio generico, non "email non trovata", perché
	// nei sistemi reali è meglio non aiutare un attaccante a capire quali email
	// sono registrate.
	if req.Email != demoUser.Email {
		// 401 Unauthorized significa: credenziali assenti o non valide.
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Secondo controllo: password.
	// CompareHashAndPassword non rigenera un hash uguale e poi confronta stringhe:
	// legge le informazioni dentro l'hash bcrypt e verifica se la password
	// ricevuta produce un risultato compatibile.
	//
	// La funzione vuole []byte, non stringhe, perché gli algoritmi crittografici
	// lavorano su sequenze di byte.
	err = bcrypt.CompareHashAndPassword([]byte(demoUser.PasswordHash), []byte(req.Password))
	if err != nil {
		// Qualsiasi errore qui significa password non valida.
		// Anche in questo caso teniamo lo stesso messaggio generico usato per
		// l'email sbagliata.
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Se siamo arrivati qui, email e password sono corrette.
	// Per ora rispondiamo con un messaggio; poi questa response conterrà il JWT.
	response := LoginResponse{
		Message: "Login successful",
	}

	// Prima di scrivere il body, dichiariamo che la risposta è JSON.
	// Gli header HTTP vanno impostati prima di scrivere il body.
	w.Header().Set("Content-Type", "application/json")

	// Encode converte la struct Go in JSON e la scrive nella risposta.
	// LoginResponse{Message: "..."} diventa {"message":"..."}.
	json.NewEncoder(w).Encode(response)
}
