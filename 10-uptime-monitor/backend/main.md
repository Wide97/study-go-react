# main.go — teoria

## Scopo di questo file

L'entry point: apre il database, registra le route (collega ogni URL/metodo HTTP alla funzione
handler giusta) e avvia il server HTTP. Nessuna logica applicativa qui — solo *wiring*, tutto il
resto (SQL, validazione, JSON) resta negli altri file.

## Riferimento diretto

Apri **`08-database-notes/backend/main.go`**: struttura identica, cambiano solo i nomi
(`notes` → `services`) e — per ora — l'assenza di CORS (vedi sotto).

## Cosa creare, esattamente

Una funzione `func main()` in `backend/main.go` (`package main`) che, in ordine:

1. Chiama `openDatabase()`. Se ritorna un errore, `log.Fatal(err)` — senza database
   l'applicazione non ha senso che parta (stesso ragionamento già scritto in `db.md`).
2. `defer db.Close()` — assicura la chiusura della connessione quando `main` termina.
3. Crea un router con `http.NewServeMux()`.
4. Registra le route:

| Metodo + path | Handler |
|---|---|
| `GET /health` | `health` |
| `GET /services` | `servicesHandler(db)` |
| `POST /services` | `createServiceHandler(db)` |
| `PUT /services/{id}` | `updateServiceHandler(db)` |
| `DELETE /services/{id}` | `deleteServiceHandler(db)` |

5. Avvia il server con `http.ListenAndServe` su una porta a tua scelta (`08-database-notes` usa
   `:8080`), passando il mux. Avvolgi la chiamata in `log.Fatal(...)`: `ListenAndServe` ritorna solo
   se il server si ferma per un errore, e in quel caso va loggato.

## Concetti

- **`http.NewServeMux()`** — il router di Go: smista ogni richiesta in arrivo verso la funzione
  registrata per quel metodo+path. Già usato uguale in `01-todo-api` e `08-database-notes`.
- **Pattern `GET /services/{id}`** — la sintassi `{id}` nel path è supportata nativamente da
  `ServeMux` da Go 1.22 in poi (il modulo qui è su Go 1.25, quindi disponibile). È lo stesso motivo
  per cui negli handler leggi l'id con `r.PathValue("id")`, non con `mux.Vars` o simili (quelle
  servono con router esterni tipo `gorilla/mux`, qui non ti serve nessuna dipendenza in più).
- **Perché `servicesHandler(db)` e non `servicesHandler`** — qui richiami il concetto già visto in
  `handlers.md`: quattro degli handler sono funzioni che *restituiscono* un `http.HandlerFunc` via
  closure su `db`. Le chiami passando `db` **una sola volta, qui in `main`**, e il risultato è quel
  che registri nel mux. `health` invece si registra così com'è, senza chiamarla e senza parentesi
  extra: non ha bisogno di `db`.
- **Ordine di dichiarazione delle route** — con `ServeMux`, `GET /services` e
  `GET /services/{id}` (se esistesse) non andrebbero in conflitto: qui però, come già notato in
  `handlers.md`, non esiste un handler per "leggi un singolo servizio" esposto via HTTP, quindi il
  problema non si pone.

## CORS: perché qui (per ora) non serve

In `08-database-notes` il mux viene avvolto in `withCORS(mux)` perché esisteva già un frontend
React che chiamava l'API da un'origine diversa (porta diversa in dev). Qui il frontend arriva solo
con **M5**: fino ad allora puoi testare l'API con `curl`/Postman senza bisogno di CORS. Quando
arriverai al frontend, potrai riprendere `cors.go` da `08-database-notes` allo stesso modo in cui
hai riportato `openDatabase`/gli handler — te lo ricorderò a quel punto se serve.

## Estensione per M2: avviare lo scheduler

Dopo aver aperto il database (`openDatabase`) e **prima** di registrare le route del mux, aggiungi
in `main`:

1. Una chiamata a `listServices(db)` per leggere i servizi già configurati (la stessa funzione già
   usata da `servicesHandler` — qui la chiami direttamente, non dentro un handler HTTP).
2. Una chiamata a `startScheduler(db, services)` (vedi `scheduler.md`).

`startScheduler` non blocca: avvia le sue goroutine e ritorna subito, quindi `main` prosegue
normalmente verso `mux := http.NewServeMux()` e poi `http.ListenAndServe` come già scritto — i due
"mondi" (server HTTP che risponde alle richieste, scheduler che controlla i servizi in background)
girano insieme da quel momento in poi, senza che uno blocchi l'altro.

Nota: se `listServices` ritorna un errore qui, decidi tu come gestirlo — dato che senza servizi
configurati lo scheduler semplicemente non avvia nessuna goroutine (non è un errore fatale come la
mancata apertura del database), potresti anche solo loggarlo e continuare, invece di `log.Fatal`.
