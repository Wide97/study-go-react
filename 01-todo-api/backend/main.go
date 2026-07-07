package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Storage in memoria: uno slice condiviso da tutti gli handler.
// Non c'è un database, quindi i dati si perdono ad ogni riavvio del server.
var todos = []Todo{
	{ID: 1, Title: "Learn Go", Done: false},
	{ID: 2, Title: "Build a REST API", Done: false},
	{ID: 3, Title: "Deploy to production", Done: false},
}

// Le struct rappresentano dati strutturati in Go.
// I tag `json:"..."` dicono a encoding/json come chiamare i campi
// quando la struct viene convertita da/verso JSON (i nomi Go sono
// maiuscoli per essere esportati, il JSON può usarne altri).
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	// mux (ServeMux) è il router: decide quale handler chiamare in base
	// a metodo+path. Prima usavamo http.HandleFunc, che registra sul mux
	// "di default" globale — qui ne creiamo uno esplicito per poterlo
	// avvolgere tutto insieme con withCORS in un colpo solo.
	// Il pattern può includere il metodo HTTP ("PUT ...") e parametri
	// nel path ("{id}"), letti poi con r.PathValue("id") dentro l'handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/todos", todosHandler)
	mux.HandleFunc("PUT /todos/{id}", toggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodo)

	log.Println("Server is running on port 8080")
	// ListenAndServe blocca ed è in ascolto finché non c'è un errore fatale
	// (es. porta già occupata). log.Fatal stampa l'errore e chiude il programma:
	// senza, un errore all'avvio passerebbe inosservato.
	// withCORS avvolge l'intero mux (non più i singoli handler): così gira
	// per QUALSIASI richiesta, incluse le OPTIONS di preflight — se avvolgessimo
	// i singoli handler, un'OPTIONS su "/todos/{id}" non troverebbe nessun
	// pattern registrato per quel metodo e ServeMux risponderebbe 405
	// prima ancora di raggiungere withCORS, facendo fallire il preflight.
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

// Un handler ha sempre questa firma: riceve dove scrivere la risposta (w)
// e i dati della richiesta in arrivo (r).
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

// Un solo handler per /todos che si comporta diversamente in base al metodo:
// GET restituisce la lista, POST ne crea una nuova.
func todosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newTodo Todo
		// json.NewDecoder(r.Body).Decode legge il body della richiesta (JSON)
		// e lo scrive dentro newTodo. Serve &newTodo (l'indirizzo) perché
		// Decode deve poter modificare la variabile.
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			// Body non decodificabile (JSON malformato) => colpa del client, 400.
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// ID finto/semplice, va bene per un progetto di studio in memoria.
		newTodo.ID = len(todos) + 1
		// append aggiunge in coda e restituisce un nuovo slice: va riassegnato.
		todos = append(todos, newTodo)

		w.Header().Set("Content-Type", "application/json")
		// L'ordine conta: header e status vanno impostati PRIMA di scrivere il body.
		w.WriteHeader(http.StatusCreated) // 201: risorsa creata
		json.NewEncoder(w).Encode(newTodo)

	} else {
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode è il simmetrico di Decode: converte
		// una struct/slice Go in JSON e lo scrive nella risposta.
		json.NewEncoder(w).Encode(todos)
	}
}

// PUT /todos/{id}: inverte lo stato Done della todo con quell'id.
func toggleTodo(w http.ResponseWriter, r *http.Request) {
	// r.PathValue legge il segmento {id} catturato dal pattern di route.
	idStr := r.PathValue("id")
	// L'id nell'URL è sempre una stringa: va convertito nel tipo int
	// usato dalla struct Todo. La conversione può fallire (es. "/todos/abc"),
	// quindi va sempre controllato l'errore.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Scorriamo lo slice cercando l'elemento con id corrispondente.
	// i è l'indice (serve per modificare todos[i] direttamente),
	// todo è una COPIA dell'elemento (utile solo per leggerne i campi).
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Done = !todos[i].Done
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	// Ciclo finito senza return => nessun elemento trovato con quell'id.
	http.Error(w, "Todo not found", http.StatusNotFound)
}

// DELETE /todos/{id}: rimuove la todo con quell'id dallo slice.
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			// Go non ha un "remove" sugli slice: si concatena la parte
			// prima dell'indice (todos[:i]) con quella dopo (todos[i+1:]).
			// "..." spacchetta il secondo slice come argomenti singoli di append,
			// invece di aggiungerlo come un unico elemento annidato.
			todos = append(todos[:i], todos[i+1:]...)
			// 204 No Content: operazione riuscita, niente da restituire nel body.
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

// withCORS è un middleware: una funzione che prende un handler (next)
// e ne restituisce uno nuovo che fa qualcosa in più prima di eseguirlo.
// Serve perché il browser blocca di default le richieste fetch verso
// un origin diverso (protocollo+dominio+porta) da quello della pagina:
// il frontend Vite gira su :5173, il backend Go su :8080, quindi per
// il browser sono origin diversi anche se sono entrambi "localhost".
// Gli header sotto dicono esplicitamente al browser "fidati di :5173".
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Per metodi come PUT/DELETE (o richieste con Content-Type custom),
		// il browser manda prima una OPTIONS "di controllo" (preflight) per
		// chiedere il permesso, e manda la richiesta vera solo se la risposta
		// a questa OPTIONS è positiva. Qui rispondiamo subito, senza inoltrarla
		// a next: gli header CORS impostati sopra bastano per farla accettare.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204: preflight ok, nessun body
			return
		}

		next.ServeHTTP(w, r) // passa la richiesta all'handler originale
	})
}
