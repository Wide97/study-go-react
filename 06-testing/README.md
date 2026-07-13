# 06 — Testing Go + React

Scopo: prendere una piccola app full-stack già nota e aggiungere test automatici, senza introdurre nuova logica applicativa.

## Idea del progetto

La base è simile al progetto `05-filtered-api-table`: backend Go con `GET /orders` e frontend React con tabella, filtri, paginazione e debounce.

Qui il focus non è costruire nuove feature, ma imparare a verificare quelle esistenti.

## Concetti coperti

### Backend Go

- file `*_test.go`
- package `testing`
- pattern arrange / act / assert
- test di funzioni pure
- test di query param
- test di handler HTTP con `httptest`
- decoding JSON della response

### Frontend React

- Vitest
- Testing Library
- setup test con `jest-dom`
- mock di `fetch`
- test asincroni con `findBy...` e `waitFor`
- test di stato vuoto
- test di filtri e debounce

## Stato: completato

- [x] Backend: test funzioni filtro
- [x] Backend: test paginazione
- [x] Backend: test parsing query params
- [x] Backend: test handler `GET /orders`
- [x] Frontend: setup Vitest + Testing Library
- [x] Frontend: test rendering dati
- [x] Frontend: test stato vuoto
- [x] Frontend: test filtro status
- [x] Frontend: test debounce ricerca
- [x] Frontend: test errore HTTP

## Comandi

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm test
npm run build
npm run lint
```

Nota: in questa sessione le dipendenze frontend di test sono state riusate da un `node_modules` locale già presente, perché l'installazione via rete non era disponibile.
