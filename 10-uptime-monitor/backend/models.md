# models.go — teoria

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

## Concetto nuovo: `time.Time` per i timestamp

Finora non hai mai avuto un campo data/ora persistito. Il pacchetto standard `time` di Go offre il
tipo `time.Time`, che rappresenta un istante preciso (data + ora + timezone). Alcune cose da sapere:

- Il driver SQLite che userai (lo stesso di `08-database-notes`) sa convertire `time.Time` da/verso
  il database automaticamente — non devi occupartene manualmente in questo file, è compito di
  `db.go`/`repository.go`. Qui ti basta *dichiarare* il campo con tipo `time.Time`.
- Quando questo campo viene serializzato in JSON (per il frontend), Go lo converte automaticamente
  in una stringa in formato RFC3339 (es. `"2026-07-16T10:30:00Z"`) — è lo standard che userai anche
  lato frontend per interpretarlo.
- Ti servirà importare il pacchetto `"time"` in questo file per poter usare il tipo.

## Cosa deve contenere `Service`

Una struct con questi campi (nome del campo Go, tipo, e perché serve):

| Campo | Tipo | Perché |
|---|---|---|
| `ID` | intero | Identificativo univoco, generato dal database (come in 08-database-notes) |
| `Name` | stringa | Etichetta leggibile per l'utente (es. "NAS - Portainer") |
| `URL` | stringa | L'indirizzo che lo scheduler dovrà controllare (es. `http://192.168.1.50:9000`) |
| `IntervalSeconds` | intero | Ogni quanti secondi controllare questo servizio — servizi diversi possono avere frequenze diverse (un servizio critico ogni 30s, uno secondario ogni 5 minuti) |
| `CreatedAt` | `time.Time` | Quando il servizio è stato aggiunto alla configurazione |

Ricordati i tag JSON per ogni campo (stesso pattern già usato: nome del campo Go in maiuscolo,
tag JSON in minuscolo/snake_case o camelCase a tua scelta, basta essere coerente).

## Cosa NON deve contenere (ancora)

Nessun campo relativo allo stato attuale del servizio (su/giù, ultima risposta, ecc.) — quello
è lo scopo della struct `Check`, che vedremo in M2. Tenere `Service` (configurazione) separato da
`Check` (risultati nel tempo) è una scelta di design: la configurazione cambia raramente, i
risultati si accumulano continuamente — sono due concetti diversi anche se collegati.
