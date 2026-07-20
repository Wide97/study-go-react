# 10 — Uptime Monitor

Applicativo reale: monitora periodicamente una lista di servizi/URL (es. quelli sul tuo NAS),
tiene uno storico su SQLite, calcola uptime % e tempi di risposta, e notifica quando un servizio
cambia stato (su→giù, giù→su).

## Regola di questo progetto

**Nessun codice scritto da me.** Per ogni file che scriverai trovi un `.md` con lo stesso nome
accanto (es. `models.go` ↔ `models.md`) con teoria e contesto — il contenuto del file lo scrivi
sempre tu. Io scrivo solo file `.md`, il `README.md`, e mi occupo dello scaffolding iniziale
(cartelle, `go mod init`, `npm create vite`, dipendenze) — mai la logica applicativa.

## Perché è "serio" rispetto ai progetti 01-09

- Non è un CRUD passivo: c'è un **processo in background** che agisce da solo nel tempo (scheduler
  che esegue controlli periodici), non solo in risposta a richieste HTTP.
- Calcolo reale di metriche nel tempo (uptime %, tempo di risposta medio) da dati storici.
- Integrazione con un servizio esterno per notifiche (webhook in uscita), primo assaggio di
  "il mio backend chiama un altro servizio", non solo "un frontend chiama il mio backend".
- Pensato per girare davvero sul NAS e monitorare servizi reali.

## Architettura, in breve

```
scheduler (goroutine + time.Ticker)
    │
    ├─▶ per ogni servizio configurato: HTTP GET con timeout, misura latenza e stato
    │
    ▼
SQLite: tabella "checks" (storico) + tabella "services" (configurazione)
    │
    ├─▶ se lo stato è cambiato rispetto all'ultimo check → notifica (webhook in uscita)
    │
    ▼
API REST (stessa struttura vista in 08-database-notes) → frontend React (dashboard stato + storico)
```

## Milestone

- [x] **M1 — Dati e CRUD servizi**: modello `Service`, SQLite, CRUD via API REST
      (`GET/POST/PUT/DELETE /services`). Nessun controllo reale ancora, solo configurazione.
- [ ] **M2 — Lo scheduler**: goroutine che, ogni tot secondi, esegue un HTTP check su ogni
      servizio configurato e salva il risultato (tabella `checks`).
- [ ] **M3 — Statistiche**: endpoint che calcola uptime % e tempo di risposta medio/storico per
      un servizio, a partire dai dati grezzi in `checks`.
- [ ] **M4 — Notifiche**: rilevare transizioni di stato (su→giù, giù→su) e mandare un webhook in
      uscita (es. verso ntfy.sh, Discord, o un endpoint generico a scelta tua).
- [ ] **M5 — Frontend**: dashboard con stato attuale di ogni servizio, storico/uptime, form per
      aggiungere/modificare servizi monitorati.
- [ ] **M6 — Deploy**: Dockerfile + docker-compose, come nel progetto 09.

## Stato

- [x] Backend Go inizializzato (`backend/go.mod`)
- [x] M1 — Dati e CRUD servizi (`db.go`, `models.go`, `repository.go`, `handlers.go`, `main.go`
      completi; `go build`/`go vet` puliti)
- [ ] M2 — Lo scheduler (**prossimo passo**)
