package main

import (
	// context serve per portare dati legati a una singola richiesta HTTP.
	// Il middleware lo usa per passare l'email autenticata all'handler finale.
	"context"

	// encoding/json serve per trasformare JSON <-> struct Go.
	// Nel login lo usiamo in ingresso: leggiamo il body della richiesta HTTP,
	// che arriva come JSON, e lo decodifichiamo dentro LoginRequest.
	"encoding/json"

	// strings contiene funzioni utili per lavorare con stringhe.
	// Qui lo usiamo per controllare ed eliminare il prefisso "Bearer "
	// dall'header Authorization.
	"strings"

	// time serve per lavorare con date e durate.
	// Nel JWT lo usiamo per calcolare la scadenza del token: "adesso + 1 ora".
	"time"

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

	// jwt/v5 è la libreria che usiamo per creare, firmare e verificare
	// JSON Web Token.
	"github.com/golang-jwt/jwt/v5"

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

type contextKey string

const emailContextKey contextKey = "email"

// jwtSecret è la chiave usata per firmare il token.
// Con HS256 la firma è "simmetrica": la stessa chiave serve sia per generare
// il token sia per verificarlo quando il client lo rimanderà al backend.
//
// In un progetto reale questa stringa NON starebbe hardcoded nel codice:
// verrebbe letta da variabile d'ambiente o secret manager.
var jwtSecret = []byte("dev-secret")

// MeResponse descrive il JSON restituito da GET /me.
// È una rotta protetta: risponde solo se il client manda un JWT valido.
type MeResponse struct {
	Email string `json:"email"`
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

// LoginResponse descrive il JSON che il backend manda quando il login riesce:
// un token JWT che il frontend dovrà conservare e rimandare nelle richieste
// protette con header Authorization.
type LoginResponse struct {
	Token string `json:"token"`
}

func main() {
	// ServeMux è il router standard di net/http.
	// Il suo compito è associare una coppia metodo+path a una funzione handler.
	mux := http.NewServeMux()

	// Da Go 1.22 il pattern può includere anche il metodo HTTP.
	// "GET /health" significa: questa funzione risponde solo a GET /health.
	// Se provi POST /health, il router può rispondere 405 Method Not Allowed.
	mux.HandleFunc("GET /health", health)

	// GET /me è la prima rotta protetta.
	// Non basta chiamarla: il client deve mandare un header Authorization
	// con un token valido ricevuto da POST /login.
	mux.Handle("GET /me", withAuth(http.HandlerFunc(me)))

	// POST /login riceve credenziali nel body.
	// Usiamo POST, non GET, perché stiamo inviando dati sensibili e perché il
	// login è un'azione, non una semplice lettura di risorsa.
	mux.HandleFunc("POST /login", login)

	log.Println("Server running on :8080")

	// ListenAndServe mette il server in ascolto sulla porta 8080.
	// È una chiamata bloccante: il programma resta qui finché il server gira.
	// Se qualcosa va storto, per esempio porta occupata, restituisce un errore.
	// log.Fatal stampa quell'errore e chiude il programma.
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
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
// 4. generare un JWT;
// 5. rispondere JSON con il token se il login è valido.
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
	// Ora generiamo un token firmato: sarà la "prova" che il client ha fatto
	// login correttamente.
	token, err := generateToken(demoUser.Email)
	if err != nil {
		// Se non riusciamo a generare il token, non è colpa del client:
		// il problema è interno al backend, quindi 500.
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: token,
	}

	// Prima di scrivere il body, dichiariamo che la risposta è JSON.
	// Gli header HTTP vanno impostati prima di scrivere il body.
	w.Header().Set("Content-Type", "application/json")

	// Encode converte la struct Go in JSON e la scrive nella risposta.
	// LoginResponse{Token: "..."} diventa {"token":"..."}.
	json.NewEncoder(w).Encode(response)
}

// generateToken crea un JWT per l'utente autenticato.
//
// Il valore di ritorno è una stringa già firmata, pronta da mandare al frontend.
// Il secondo valore è un error perché la firma può fallire.
func generateToken(userEmail string) (string, error) {
	// MapClaims è una mappa libera di dati da mettere nel token.
	// Usiamo due claim standard:
	// - sub: subject, cioè "di chi parla" il token. Qui è l'email utente.
	// - exp: expiration time, cioè quando il token scade.
	//
	// exp nei JWT è un timestamp Unix: numero di secondi dal 1 gennaio 1970.
	claims := jwt.MapClaims{
		"sub": userEmail,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// NewWithClaims crea l'oggetto token ma non lo firma ancora.
	// SigningMethodHS256 indica l'algoritmo di firma: HMAC + SHA-256.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignedString firma il token usando jwtSecret e restituisce la stringa
	// finale composta da tre parti separate da punti:
	// header.payload.signature
	return token.SignedString(jwtSecret)
}

// me gestisce GET /me.
// È una rotta protetta: prima valida il JWT, poi restituisce i dati dell'utente.
//
// L'handler non legge direttamente l'header Authorization: ci pensa withAuth.
// Qui leggiamo solo l'email che il middleware ha salvato nel context.
func me(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(emailContextKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := MeResponse{
		Email: email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getEmailFromToken contiene la logica riusabile di autenticazione:
// prende una richiesta HTTP, cerca il JWT nell'header Authorization,
// lo valida e restituisce l'email dentro il claim sub.
func getEmailFromToken(r *http.Request) (string, error) {
	// Il token arriva in un header HTTP, non nel body.
	// Lo standard pratico più comune è:
	// Authorization: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	// Controlliamo il prefisso perché Authorization può contenere schemi
	// diversi. Qui accettiamo solo Bearer token.
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header")
	}

	// Tolto "Bearer ", resta solo la stringa JWT: header.payload.signature.
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// jwt.Parse legge il token, verifica la firma e valida i claim standard
	// supportati dalla libreria, come exp.
	//
	// La funzione callback serve alla libreria per sapere quale chiave usare
	// per verificare la firma. Prima controlliamo anche che l'algoritmo sia
	// HMAC, cioè la famiglia di HS256 usata da generateToken.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	// Se Parse ha dato errore, oppure il token esiste ma non è valido, la
	// richiesta non è autenticata.
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// I claims sono i dati dentro il token.
	// Siccome in generateToken abbiamo usato jwt.MapClaims, qui ci aspettiamo
	// di poter leggere token.Claims come jwt.MapClaims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// sub è il subject del token: nel nostro caso l'email dell'utente.
	// Il valore dentro una mappa interface{} va convertito al tipo atteso,
	// quindi controlliamo che sia davvero una stringa.
	email, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token subject")
	}
	return email, nil
}

// withAuth è un middleware di autenticazione.
// Prende un handler "protetto" (next) e restituisce un nuovo handler che:
// 1. valida il token JWT;
// 2. salva l'email utente nel context della request;
// 3. chiama l'handler finale solo se l'autenticazione è valida.
func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, err := getEmailFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// context.WithValue non modifica il context originale: ne crea uno nuovo
		// con un valore in più. r.WithContext crea poi una copia della request
		// che porta quel context aggiornato.
		ctx := context.WithValue(r.Context(), emailContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withCORS è un middleware che permette al frontend Vite di chiamare il backend.
// Browser e server considerano "origin" diversi se cambia protocollo, dominio
// o porta: http://localhost:5173 e http://localhost:8080 sono quindi diversi.
//
// Senza questi header, il browser bloccherebbe le fetch del frontend prima
// ancora di consegnare la risposta al codice React.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permettiamo solo il frontend locale usato da Vite.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// Metodi che il browser può usare verso questa API.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Header custom ammessi. Authorization serve per mandare il Bearer token.
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Prima di alcune richieste, il browser manda una OPTIONS di preflight:
		// chiede al backend se la richiesta vera sarà permessa.
		// Rispondiamo subito 204, senza passare al router.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
