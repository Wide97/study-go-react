# models.go — teoria

## Cosa creare, esattamente

1. Crea il file `backend/models.go` (nella cartella dove hai già `go.mod`).
2. Prima riga: dichiarazione del package. In tutti i progetti precedenti (01, 08...) hai sempre
   usato `package main` per il backend — fai lo stesso qui, sarà lo stesso package per tutti i
   file di questo backend (`db.go`, `handlers.go`, ecc., quando li scriveremo).
3. Subito dopo, il blocco `import` con l'unico pacchetto che ti serve in questo file: `"time"`
   (vedi sotto il perché).
4. Definisci un **tipo struct** chiamato `Service`, con i campi elencati nella tabella più sotto
   ("Cosa deve contenere `Service`") — stesso pattern di struct con tag JSON che hai già scritto
   in `01-todo-api` (`type Todo struct { ... }`) e in `08-database-notes` (`type Note struct { ... }`):
   ogni campo ha un nome Go (maiuscolo) e un tag `` `json:"..."` `` con il nome in minuscolo.

Solo questo: la struct `Service`. Nessuna funzione, nessuna logica — è solo la definizione del
"tipo di dato", il resto (creare la tabella nel database, le funzioni CRUD, gli handler HTTP)
arriva nei prossimi file/step, ognuno con il proprio `.md`.

## Scopo di questo file

Contiene i tipi Go che rappresentano i dati dell'applicazione — in questa milestone (M1), solo
il **servizio monitorato** (`Service`): la configurazione di "cosa" controllare, non ancora i
risultati dei controlli (quelli arriveranno in M2, quando aggiungeremo `Check`).

## Concetti già noti (li hai già usati in 01-todo-api e 08-database-notes)

- Una **struct** Go per rappresentare un record.
- I **tag `json:"..."`** sui campi, per controllare come la struct viene serializzata/deserializzata
  quando parla con il frontend (`encoding/json`).
- Il campo `ID` come intero auto-incrementante, gestito dal database (come in `08-database-notes`),
  non generato a mano lato applicazione.

## Timestamp: stesso pattern di `08-database-notes`

**Correzione rispetto a una versione precedente di questa nota**: qui NON usiamo il tipo
`time.Time` per il campo data. In `08-database-notes` la scelta (che funziona bene ed è più
semplice da gestire col driver SQLite in uso) è stata salvare i timestamp come **stringa** in
formato RFC3339 (es. `"2026-07-16T10:30:00+02:00"`), generata con `time.Now().Format(time.RFC3339)`
nel repository al momento dell'inserimento — non come tipo `time.Time` nella struct.

Quindi: il campo `CreatedAt` in `Service` è di tipo **stringa**, non `time.Time`. Il pacchetto
`"time"` non serve in *questo* file (`models.go`) — servirà più avanti in `repository.go`, dove
genereremo davvero il timestamp. Se lo hai già importato in `models.go`, puoi toglierlo: un
import non usato è un errore di compilazione in Go.

## Cosa deve contenere `Service`

Una struct con questi campi (nome del campo Go, tipo, e perché serve):

| Campo | Tipo | Perché |
|---|---|---|
| `ID` | intero | Identificativo univoco, generato dal database (come in 08-database-notes) |
| `Name` | stringa | Etichetta leggibile per l'utente (es. "NAS - Portainer") |
| `URL` | stringa | L'indirizzo che lo scheduler dovrà controllare (es. `http://192.168.1.50:9000`) |
| `IntervalSeconds` | intero | Ogni quanti secondi controllare questo servizio — servizi diversi possono avere frequenze diverse (un servizio critico ogni 30s, uno secondario ogni 5 minuti) |
| `CreatedAt` | stringa | Quando il servizio è stato aggiunto alla configurazione (formato RFC3339, stesso pattern di `08-database-notes`) |

Ricordati i tag JSON per ogni campo (stesso pattern già usato: nome del campo Go in maiuscolo,
tag JSON in minuscolo/snake_case o camelCase a tua scelta, basta essere coerente).

## Un secondo tipo: `ServiceRequest`

Stesso pattern di `08-database-notes` (`Note` vs `NoteRequest`): oltre a `Service` (il record
completo, con `ID` e `CreatedAt` generati dal database), serve un tipo più piccolo che rappresenta
**solo i dati che un client può mandare** quando crea o modifica un servizio — senza `ID` né
`CreatedAt`, che non deve poter decidere lui.

Aggiungi in `models.go` un secondo tipo, `ServiceRequest`, con solo tre campi: `Name`, `URL`,
`IntervalSeconds` (stessi tipi e tag JSON già usati in `Service` per questi tre campi).

Lo userai nei prossimi file (repository, handler) come tipo di ingresso per `POST`/`PUT`, mentre
`Service` resterà il tipo usato per leggere/restituire i dati già salvati.

## Cosa NON deve contenere (ancora)

Nessun campo relativo allo stato attuale del servizio (su/giù, ultima risposta, ecc.) — quello
è lo scopo della struct `Check`, che vedremo in M2. Tenere `Service` (configurazione) separato da
`Check` (risultati nel tempo) è una scelta di design: la configurazione cambia raramente, i
risultati si accumulano continuamente — sono due concetti diversi anche se collegati.
