# Percorso di studio: Go + React

## Punto di partenza

- **Go**: lo usa lato lavoro (gateway Smart Eye) ma non è autonomo nello scrivere codice da zero. Non dare per scontate le basi.
- **React**: in fase di apprendimento, stesso discorso — non è autonomo a scrivere da zero.

In generale: **non presumere autonomia nello scrivere codice da zero, in nessuno dei due**. Il percorso serve anche a costruire quella capacità, non solo a "consolidare".

## Obiettivo

Nessun obiettivo finale rigido: una raccolta di **piccoli progetti indipendenti**, da affrontare con calma quando c'è voglia/tempo, senza scadenze o ritmo fisso. Ogni progetto è a sé stante — si può abbandonare, riprendere o saltare senza dipendenze dagli altri.

## Metodo: guida, non soluzioni

Regola principale: **non mostrare codice già scritto come prima cosa**. L'obiettivo è che il codice lo scriva lui, guidato passo passo — altrimenti non impara a scrivere da zero, solo a leggere/copiare.

Per ogni progetto/step:

1. **Teoria/concetti** — breve blocco sui concetti coinvolti (nuovo hook React, pattern Go, protocollo, ecc.), giusto quanto basta per capire il _perché_ prima del _come_.
2. **Guida passo passo** — spiegare cosa fare e perché, con domande o indicazioni ("ora ti serve una funzione che...", "prova a scrivere..."), lasciando scrivere il codice a lui. Mostrare frammenti di codice solo se serve davvero a sbloccare un punto specifico, non come soluzione completa.
3. **Revisione** — dopo che ha scritto qualcosa, guardarlo insieme, correggere, spiegare gli errori.
4. **Note finali** (facoltative) — cosa è emerso, dubbi rimasti, cose da approfondire dopo.

## Struttura cartelle

Un progetto = una sottocartella in `study/`, con un proprio README minimo (scopo, concetti coperti, stato: da iniziare / in corso / concluso).

```
study/
  STUDY.md          <- questo file
  01-todo-api/
  02-realtime-dashboard/
  ...
```

## Backlog progetti (idee, in ordine di complessità crescente)

Focus su full-stack Go (backend) + React (frontend) per far lavorare insieme i due lati, spingendo soprattutto sul fronte React dato che è la parte da consolidare.

- [x] **To-do CRUD** — completato in `01-todo-api/` (Go net/http + React/TS + Bootstrap). Vedi il README del progetto per dettagli e concetti coperti.
- [x] **Auth semplice (JWT)** — login/logout, rotte protette lato React (context o router guard), middleware di auth in Go. Concetti: Context API, protected routes, gestione token/refresh.
- [x] **Dashboard realtime** — completato in `03-realtime-dashboard/` (Go WebSocket + React/TS). Concetti: WebSocket lato client, gestione stato che cambia in tempo reale.
- [x] **App con stato globale complesso** — completato in `04-global-state-app/` (mini carrello full-stack con Go + React). Concetti: stato locale vs globale, Context API, provider, custom hook, derived state.
- [x] **API con paginazione/filtri/ricerca** — completato in `05-filtered-api-table/`. Backend Go con query params, frontend con UI di filtro e tabella. Concetti: debounce, gestione query string, componenti tabella riutilizzabili.
- [x] **Testing** — completato in `06-testing/`: unit test Go (`testing` package) e test React con Vitest/Testing Library. Concetti: mocking, test di componenti, test di handler HTTP.
- [x] **Strutturazione progetto multi-file** — completato in `07-project-structure/`: preso un progetto già fatto e separato `main.go` in file dedicati e `App.tsx` in componenti/file reali. Concetti: separazione responsabilità, props tra componenti, import/export, package Go, handler/service/model, organizzazione cartelle.
- [x] **Persistenza con database** — completato in `08-database-notes/`: mini app note full-stack con Go, React e SQLite. Concetti: `database/sql`, driver SQLite, connessione DB, `SELECT/INSERT/UPDATE/DELETE`, gestione errori, repository layer, CRUD frontend, form controllati, sincronizzazione stato/backend.
- [x] **Deploy/containerizzazione** — completato in `09-deploy-containerization/`: Dockerfile per backend Go + frontend React, Docker Compose e verifica runtime. Concetti: multi-stage build, reverse proxy nginx per servire frontend e inoltrare API, rete tra servizi e build riproducibile con `npm ci`.

Nuove idee si possono aggiungere liberamente man mano che vengono in mente, semplicemente spuntando/aggiungendo righe qui sopra.

## Dopo il backlog: progetti "seri"

Backlog sopra completato interamente. Dal progetto `10-uptime-monitor` in poi si passa ad
applicativi "seri" (uso reale, non esercizio di consolidamento): per questi vale una regola più
stretta di quella generale di questo file — nessun codice scritto da Claude, solo file `.md` di
teoria affiancati al codice scritto da Marco. Dettagli, motivazione e stato del progetto attivo in
`.ai-docs/conventions.md` e `.ai-docs/current-state.md` (non duplicati qui).
