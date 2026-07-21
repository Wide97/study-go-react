# scheduler.go — teoria

## Scopo di questo file

Il pezzo che rende questo progetto diverso da tutti i precedenti (vedi README, "perché è serio"):
un processo che gira **da solo nel tempo**, non in risposta a una richiesta HTTP. Controlla
periodicamente ogni servizio configurato e salva l'esito con `recordCheck` (già scritta).

## Un concetto che hai già usato, ma non così

In `03-realtime-dashboard` hai già scritto `time.NewTicker` e il pattern `<-ticker.C` dentro un
loop infinito (`websocketHandler`). Riusa esattamente quel pattern qui. La differenza è **dove**
gira quel loop:

- In `03-realtime-dashboard`, il loop viveva **dentro l'handler** di una richiesta HTTP (`/ws`):
  girava perché `net/http` mette automaticamente ogni richiesta in arrivo nella sua goroutine —
  non l'hai avviata tu esplicitamente, è già "concorrente" di suo.
- Qui, invece, il loop deve partire **una volta sola all'avvio del programma**, non in risposta a
  nessuna richiesta — nessuno "chiama" lo scheduler, deve semplicemente esserci, girare per tutta
  la vita del programma. Per questo serve **avviarlo esplicitamente come goroutine** con la keyword
  `go`, cosa che finora non hai ancora fatto in nessun progetto precedente.

## Cosa creare, esattamente

Due funzioni in `backend/scheduler.go` (`package main`):

### 1. `startScheduler(db *sql.DB, services []Service)`

Non restituisce nulla, non blocca: per **ogni** servizio in `services`, avvia una goroutine
indipendente con `go func() { ... }()` che contiene un loop infinito con lo stesso pattern già
visto in `03-realtime-dashboard`:

```
ticker := time.NewTicker(...)
for {
    <-ticker.C
    // esegui il check di questo servizio, salva con recordCheck
}
```

L'intervallo del ticker per ogni servizio è `time.Duration(service.IntervalSeconds) * time.Second`
— non un valore fisso: servizi diversi hanno `IntervalSeconds` diversi (è il motivo per cui quel
campo esiste in `Service`).

**Perché una goroutine per servizio e non un ticker solo per tutti**: se un servizio va controllato
ogni 30s e un altro ogni 300s, un singolo ticker condiviso non può rispettare cadenze diverse —
ogni servizio ha bisogno del proprio ticker indipendente, quindi della propria goroutine.

### 2. `checkService(client *http.Client, service Service) (status string, responseTimeMs int)`

La funzione che fa il lavoro vero per un singolo controllo:
l'
1. Segna l'istante di partenza: `start := time.Now()`.
2. Esegue `client.Get(service.URL)`.
3. Calcola il tempo trascorso: `time.Since(start)`, convertito in millisecondi
   (`.Milliseconds()` sul risultato, che è un `time.Duration`).
4. Decide lo status: se `client.Get` ritorna un errore (timeout, connessione rifiutata, host
   irraggiungibile...) → `"down"`, `responseTimeMs` a `0` (vedi `models.md`: è la scelta già presa
   per M2). Se non c'è errore, guarda `resp.StatusCode`: sotto `400` → `"up"`, da `400` in su →
   `"down"` (un 500 vuol dire che il servizio ha risposto ma è in errore — non è "su" in senso utile).
5. **Importante**: se `client.Get` non ha dato errore, hai un `resp.Body` da chiudere con
   `defer resp.Body.Close()` — anche se non ti serve leggere il body, va chiuso comunque per non
   lasciare la connessione aperta. Novità rispetto a quanto hai fatto finora: qui è il **tuo codice
   a fare da client** HTTP verso un altro servizio, non da server che risponde — è la prima volta in
   questo percorso che il tuo backend chiama un altro servizio invece di essere chiamato.

## Il `client *http.Client` con timeout

Non usare `http.Get(url)` diretto (la funzione di pacchetto): quella usa un client di default
**senza timeout** — se un servizio non risponde mai, quella chiamata resterebbe bloccata per
sempre, bloccando anche il ticker di quel servizio. Crea invece un `http.Client{Timeout: ...}` (es.
5 secondi) e passalo come parametro a `checkService` — puoi crearlo una volta sola in
`startScheduler` e riusarlo per tutti i check di tutti i servizi (un `http.Client` è sicuro da usare
da più goroutine contemporaneamente, non serve un client per servizio).

## `*sql.DB` e concorrenza

Più goroutine (una per servizio) chiameranno `recordCheck(db, ...)` in momenti diversi,
potenzialmente in parallelo. Non serve nessuna precauzione in più da parte tua: `*sql.DB` in Go è
progettato per essere usato da più goroutine contemporaneamente senza corse critiche — gestisce lui
internamente un pool di connessioni. È lo stesso `db` che passa già `main.go` agli handler HTTP
(anche quelli, ricorda, girano ognuno nella propria goroutine per via di `net/http`).

## Nota su Go 1.22+ (questo modulo è su Go 1.25): niente da temere sulla variabile di loop

Se cerchi esempi online su "goroutine dentro un for", troverai spesso l'avviso di passare la
variabile di loop come parametro alla goroutine (`go func(s Service) {...}(service)`) invece di
usarla direttamente, per evitare un bug classico di Go (tutte le goroutine finivano per condividere
la stessa variabile, vedendo solo l'ultimo valore). **Da Go 1.22 in poi questo bug non esiste più**:
ogni iterazione del `for range` crea una variabile nuova, quindi puoi scrivere semplicemente
`go func() { ... service ... }()` dentro il loop, usando `service` direttamente, senza passarla come
parametro. Lo segnalo perché è un dettaglio che cambia a seconda della versione di Go e molte guide
in giro sono scritte per versioni precedenti.

## Dove viene chiamato

`startScheduler` non la richiami tu da nessuna parte in questo file — la richiamerà `main.go`,
prossimo (piccolo) aggiornamento: vedi la sezione aggiunta in `main.md`.
