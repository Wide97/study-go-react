# db.go — teoria

## Cosa creare, esattamente

1. Crea il file `backend/db.go`, stesso package (`package main`).
2. Import necessari: `"database/sql"` (pacchetto standard per parlare con database SQL) e
   `_ "modernc.org/sqlite"` — nota il `_` davanti: è un **blank import**, spiegato sotto.
3. Scrivi una funzione `openDatabase() (*sql.DB, error)` che:
   - apre la connessione al file SQLite (nome a tua scelta, es. `"uptime.db"`),
   - verifica che la connessione funzioni davvero,
   - crea la tabella `services` se non esiste già,
   - restituisce la connessione pronta all'uso (o un errore, se qualcosa fallisce).

Stessa struttura di `openDatabase` che hai già scritto in `08-database-notes/backend/db.go` —
qui cambia solo la tabella creata (`services` invece di `notes`), con colonne coerenti con la
struct `Service` di `models.go`.

## Scopo di questo file

Centralizza in un unico posto **come si apre la connessione al database** e **come è fatta la
tabella** — nessun altro file (repository, handler) deve occuparsi di questi dettagli, li
chiederà a `openDatabase`.

## Concetti (già visti in 08-database-notes, qui solo un ripasso mirato)

- **`database/sql`** è il pacchetto standard Go per lavorare con database relazionali. Non parla
  direttamente con SQLite: ha bisogno di un **driver** specifico registrato.
- **Blank import** (`_ "modernc.org/sqlite"`): non usi mai direttamente niente di quel pacchetto
  nel tuo codice (infatti non ha un nome davanti all'import), ma importarlo esegue comunque la sua
  funzione `init()`, che si registra come driver disponibile per `database/sql`. Senza questo
  import, `sql.Open("sqlite", ...)` non saprebbe cosa fare con il nome `"sqlite"`.
- **`sql.Open`** non apre subito una connessione reale — prepara solo l'oggetto. Per verificare
  che il database sia davvero raggiungibile/apribile serve chiamare **`db.Ping()`** subito dopo
  (stesso ordine logico usato in `08-database-notes`).
- **`CREATE TABLE IF NOT EXISTS`**: eseguito ad ogni avvio del programma. La clausola
  `IF NOT EXISTS` evita errori se il programma riparte e la tabella esiste già — è il modo più
  semplice di gestire "migrazioni" quando lo schema è piccolo e stabile.

## Schema della tabella `services`

Deve rispecchiare i campi di `Service` in `models.go`:

| Colonna | Tipo SQLite | Note |
|---|---|---|
| `id` | `INTEGER PRIMARY KEY AUTOINCREMENT` | Generato dal database, come `notes.id` in 08 |
| `name` | `TEXT NOT NULL` | |
| `url` | `TEXT NOT NULL` | |
| `interval_seconds` | `INTEGER NOT NULL` | |
| `created_at` | `TEXT NOT NULL` | stringa RFC3339, come `notes.created_at`/`updated_at` in 08 |

## Da tenere a mente

`db.Exec` (per `CREATE TABLE`) restituisce un errore da controllare, come sempre — se qualunque
passaggio fallisce (`sql.Open`, `db.Ping`, `db.Exec`), `openDatabase` deve restituire quell'errore
al chiamante (sarà `main.go`, nel prossimo step, a decidere cosa farne — probabilmente `log.Fatal`,
perché senza database l'applicazione non ha senso che parta).
