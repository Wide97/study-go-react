# handlers.go — teoria

## Scopo di questo file

Lo strato HTTP: legge/scrive richieste e risposte (JSON, status code, path params), ma **non
contiene mai SQL** — delega sempre al repository (`listServices`, `createService`, ecc.). Stessa
separazione netta già vista in `08-database-notes`: repository = SQL, handler = HTTP.

## Riferimento diretto

Apri **`08-database-notes/backend/handler.go`** e **`08-database-notes/backend/main.go`** side by
side: la struttura è identica, cambiano solo i nomi (`notes` → `services`) e due dettagli di
validazione spiegati sotto. Copia la struttura, adatta i nomi.

## Cosa creare, esattamente

Cinque funzioni in `backend/handlers.go` (`package main`):

| Funzione | Firma | Gestisce |
|---|---|---|
| `health` | `func(w http.ResponseWriter, r *http.Request)` | `GET /health` |
| `servicesHandler` | `func(db *sql.DB) http.HandlerFunc` | `GET /services` |
| `createServiceHandler` | `func(db *sql.DB) http.HandlerFunc` | `POST /services` |
| `updateServiceHandler` | `func(db *sql.DB) http.HandlerFunc` | `PUT /services/{id}` |
| `deleteServiceHandler` | `func(db *sql.DB) http.HandlerFunc` | `DELETE /services/{id}` |

Nota: **non serve** un handler per "leggi un singolo servizio" (`GET /services/{id}`) — esattamente
come in `08-database-notes`, `getServiceByID`/`getNoteByID` esistono nel repository ma vengono
usate solo **internamente** (da `updateService`), non esposte come endpoint a sé.

## Concetto: perché quattro di queste funzioni "restituiscono una funzione"

Guarda la firma di `servicesHandler`: `func(db *sql.DB) http.HandlerFunc` — non è l'handler vero,
è una funzione che **costruisce e restituisce** l'handler vero. La chiami una volta sola in
`main.go` (`servicesHandler(db)`), passandole `db`; lei ti restituisce una funzione già "pronta"
che ha `db` disponibile ogni volta che viene invocata da una richiesta HTTP.

**L'hai già scritto tu stesso**, in un altro progetto: `withCORS` in `01-todo-api/backend/main.go`
ha esattamente questa forma — prende qualcosa (lì `next http.Handler`, qui `db *sql.DB`) e
restituisce un `http.HandlerFunc` che lo "ricorda" per sempre grazie alla closure. Stesso identico
meccanismo, cambia solo cosa viene "ricordato" e perché: qui serve perché `db` non è una variabile
globale (a differenza di `todos` in `01-todo-api`) — è locale a `main`, quindi ogni handler che ne
ha bisogno se lo deve far "iniettare" così.

`health`, invece, non tocca mai il database — resta una funzione handler semplice, senza wrapping,
come già in `01-todo-api` e `08-database-notes`.

## Le due differenze di validazione rispetto a `08-database-notes`

In `createNoteHandler`/`updateNoteHandler`, la validazione dopo il `Decode` era
`if req.Title == "" || req.Content == ""`. Qui i campi sono diversi:

- **`req.Name == ""`** e **`req.URL == ""`** → stesso tipo di controllo, stesso motivo (evitare di
  arrivare al database con dati chiaramente invalidi, rispondendo `400` prima).
- **In più**, un controllo che in `08-database-notes` non serviva: **`req.IntervalSeconds <= 0`**.
  Un intervallo zero o negativo non ha senso per uno scheduler che dovrà controllare il servizio
  "ogni tot secondi" (lo scheduler arriva nella prossima milestone, ma la validazione va fatta ora,
  quando accettiamo il dato) — se non lo blocchi qui, più avanti dovresti gestire un caso limite
  strano nel codice dello scheduler stesso.

Per il resto (decodifica del body, `errors.Is(err, sql.ErrNoRows)` per tradurre "non trovato" in
`404`, `strconv.Atoi(r.PathValue("id"))` per leggere l'id dal path, status code da usare) è
identico a `08-database-notes` — stessi pattern, stessi errori da gestire allo stesso modo.
