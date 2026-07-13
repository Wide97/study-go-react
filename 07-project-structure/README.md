# 07 — Strutturazione progetto multi-file

Scopo: prendere una piccola app full-stack già conosciuta e organizzarla in file e responsabilità separate.

## Idea del progetto

La base è il progetto `05-filtered-api-table`: backend Go con ordini, filtri, ricerca e paginazione; frontend React con tabella, filtri e debounce.

In questo progetto non aggiungiamo feature nuove. L'obiettivo è imparare a riconoscere cosa deve stare insieme e cosa invece va separato.

## Target struttura backend

```text
backend/
  main.go          # avvio server e registrazione rotte
  models.go        # struct Order, OrdersResponse
  data.go          # dati finti in memoria
  handlers.go      # handler HTTP
  filters.go       # filtro status/search
  pagination.go    # query param numerici e paginazione
  cors.go          # middleware CORS
```

## Target struttura frontend

```text
frontend/src/
  App.tsx
  api.ts
  types.ts
  components/
    OrdersFilters.tsx
    OrdersTable.tsx
    PaginationControls.tsx
```

## Concetti da coprire

### Backend

- package `main` distribuito su più file
- separazione tra model, dati, handler e funzioni pure
- import tra file dello stesso package
- mantenere invariato il comportamento durante un refactor
- verificare con `go test ./...`

### Frontend

- componenti in file separati
- `export` e `import`
- types condivisi in `types.ts`
- API client separato in `api.ts`
- `App.tsx` come composizione di stato + componenti
- verificare con `npm run build` e `npm run lint`

## Stato: da iniziare

- [x] Base copiata dal progetto `05-filtered-api-table`
- [x] Backend: separare `models.go`
- [x] Backend: separare `data.go`
- [x] Backend: separare `handlers.go`
- [x] Backend: separare `filters.go`
- [x] Backend: separare `pagination.go`
- [x] Backend: separare `cors.go`
- [x] Frontend: creare `types.ts`
- [x] Frontend: creare `api.ts`
- [ ] Frontend: separare `OrdersFilters`
- [ ] Frontend: separare `OrdersTable`
- [ ] Frontend: separare `PaginationControls`
- [ ] Verifica finale backend/frontend

## Metodo

Qui il punto non è scrivere più codice, ma spostare codice senza rompere nulla.

Ogni step deve seguire questa regola:

1. spostiamo una sola responsabilità;
2. compiliamo;
3. verifichiamo;
4. solo dopo passiamo allo step successivo.
