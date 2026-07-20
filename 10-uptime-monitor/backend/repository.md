# repository.go — teoria

## Scopo di questo file

Lo strato che parla SQL. Nessun HTTP, nessun JSON qui dentro — riceve/restituisce solo dati Go
(`Service`, `ServiceRequest`, errori). Gli handler HTTP (prossimo file) chiameranno queste funzioni
invece di scrivere query direttamente — stessa separazione già vista in `08-database-notes`.

## Cosa creare, esattamente

Cinque funzioni, tutte nel file `backend/repository.go` (`package main`), tutte con `db *sql.DB`
come primo parametro (la connessione aperta da `openDatabase`, che qualcun altro — `main.go` —
passerà loro):

| Funzione | Parametri (oltre a `db`) | Cosa restituisce | Cosa fa |
|---|---|---|---|
| `listServices` | — | `([]Service, error)` | Legge tutti i servizi configurati |
| `createService` | `req ServiceRequest` | `(Service, error)` | Inserisce un nuovo servizio |
| `getServiceByID` | `id int` | `(Service, error)` | Legge un singolo servizio |
| `updateService` | `id int`, `req ServiceRequest` | `(Service, error)` | Modifica un servizio esistente |
| `deleteService` | `id int` | `error` | Cancella un servizio |

Stessa identica forma di `listNotes`/`createNote`/`getNoteByID`/`updateNote`/`deleteNote` in
`08-database-notes/backend/repository.go` — vai a rileggerlo come riferimento diretto per lo
stile (uso di `db.Query`/`db.QueryRow`/`db.Exec`, `rows.Scan`, `result.LastInsertId()`,
`result.RowsAffected()`), qui cambiano solo la tabella (`services`) e i campi.

## Le due differenze rispetto a `08-database-notes` a cui fare attenzione

**1. Un solo timestamp, non due.** `Note` aveva `CreatedAt` **e** `UpdatedAt` (perché una nota
modificata cambia anche `updated_at`). `Service` ha **solo** `CreatedAt` — non l'abbiamo previsto
in `models.md` come campo che cambia nel tempo, rappresenta solo "quando è stato aggiunto". Quindi:
- `createService` genera `created_at` con `time.Now().Format(time.RFC3339)`, come `createNote`
  faceva per entrambi i suoi campi data.
- `updateService` **non tocca `created_at`** — l'`UPDATE` SQL deve modificare solo `name`, `url`,
  `interval_seconds`, lasciando `created_at` invariato (semplicemente non lo includi nella `SET`).

**2. "Non trovato" in tre punti diversi.** Come in `08-database-notes`, il modo per segnalare
"questo id non esiste" è restituire `sql.ErrNoRows`:
- in `getServiceByID`, arriva da solo (è l'errore che `QueryRow(...).Scan(...)` restituisce quando
  non trova righe — non devi generarlo tu esplicitamente);
- in `updateService` e `deleteService`, invece, `db.Exec` **non fallisce** se l'id non esiste (un
  `UPDATE`/`DELETE` su un id inesistente "riesce" comunque, semplicemente non tocca nessuna riga) —
  per questo devi controllare tu `result.RowsAffected()`: se è `0`, restituisci esplicitamente
  `sql.ErrNoRows` (stesso pattern già scritto in `updateNote`/`deleteNote`).

## Un dettaglio da non perdere: `updateService` rilegge il record

Dopo un `UPDATE` riuscito, `updateService` non deve "inventarsi" la `Service` da restituire
assemblandola a mano dai campi di `req` — richiama `getServiceByID(db, id)` alla fine e restituisce
quello, così il valore tornato è **esattamente** quello che è finito nel database (stesso motivo/
pattern di `updateNote` in `08-database-notes`).

## Estensione per M2: `recordCheck`

Una sola funzione nuova in `repository.go`:

| Funzione | Parametri (oltre a `db`) | Cosa restituisce | Cosa fa |
|---|---|---|---|
| `recordCheck` | `serviceID int`, `status string`, `responseTimeMs int` | `(Check, error)` | Inserisce un nuovo check nella tabella `checks` |

Stessa forma di `createService`: genera `checked_at` con `time.Now().Format(time.RFC3339)`,
esegue un `INSERT INTO checks (...)`, recupera l'id con `result.LastInsertId()`, e costruisce/
restituisce il `Check` completo — non serve rileggerlo dal database con una `SELECT` separata (a
differenza di `updateService`), perché qui hai già in mano tutti i valori che hai appena inserito:
non c'è nulla che il database possa aver calcolato lui e che tu non conosca già (a differenza di un
`UPDATE`, dove vuoi essere sicuro che il record letto rifletta esattamente cosa è stato scritto).

Non serve nessuna funzione per **leggere** i check in questa milestone: chi la userà è lo
scheduler (prossimo file), che chiama solo `recordCheck` dopo ogni controllo HTTP. Le funzioni di
lettura (es. "tutti i check di un servizio", "l'ultimo check di un servizio") arriveranno quando
servono davvero — in M3 (statistiche) e M4 (notifiche), non prima. Evita di scriverle ora "perché
prima o poi serviranno": codice non usato è solo rumore da mantenere.

### Perché i parametri sono tre valori separati e non un `Check` già pronto

Potresti chiederti: perché non `recordCheck(db, check Check) (Check, error)`, passando un `Check`
già assemblato? Perché chi chiama questa funzione (lo scheduler) non conosce ancora `ID` né
`CheckedAt` — li decide il repository (come in `createService`, che riceve `ServiceRequest`, non
`Service`). Passare i tre valori "grezzi" misurati dallo scheduler (`serviceID`, `status`,
`responseTimeMs`) rende esplicito cosa lo scheduler sa davvero e cosa invece è responsabilità del
repository generare.
