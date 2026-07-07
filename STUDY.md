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

1. **Teoria/concetti** — breve blocco sui concetti coinvolti (nuovo hook React, pattern Go, protocollo, ecc.), giusto quanto basta per capire il *perché* prima del *come*.
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

- [ ] **To-do CRUD** — API REST in Go (net/http o chi/echo/fiber) + frontend React con fetch/axios, gestione stato locale, form. Concetti: routing REST, hook `useState`/`useEffect`, gestione form controllati.
- [ ] **Auth semplice (JWT)** — login/logout, rotte protette lato React (context o router guard), middleware di auth in Go. Concetti: Context API, protected routes, gestione token/refresh.
- [ ] **Dashboard realtime** — backend Go che espone dati via WebSocket (es. contatore, dati finti tipo sensori), frontend React che li visualizza live con grafici. Concetti: WebSocket lato client, gestione stato che cambia in tempo reale, librerie di charting (recharts/visx).
- [ ] **App con stato globale complesso** — es. carrello o gestione multi-entità con relazioni. Concetti: Context API avanzato o libreria di state management (Zustand/Redux Toolkit), custom hooks.
- [ ] **API con paginazione/filtri/ricerca** — backend Go con query params, frontend con UI di filtro e tabella. Concetti: debounce, gestione query string, componenti tabella riutilizzabili.
- [ ] **Testing** — aggiungere test a uno dei progetti precedenti: unit test Go (`testing` package) e test React (Testing Library). Concetti: mocking, test di componenti, test di handler HTTP.
- [ ] **Deploy/containerizzazione** — Dockerfile per backend Go + frontend React, docker-compose per farli girare insieme. Concetti: multi-stage build, reverse proxy (nginx) per servire frontend + proxare API.

Nuove idee si possono aggiungere liberamente man mano che vengono in mente, semplicemente spuntando/aggiungendo righe qui sopra.
