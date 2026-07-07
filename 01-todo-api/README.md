# 01 — Todo API

Scopo: primo progetto del percorso di studio, per prendere confidenza con Go e React da zero, collegando backend e frontend.

## Concetti coperti (backend)

- `net/http`: handler, `ListenAndServe`, gestione errori con `log.Fatal`
- struct e JSON tag, `encoding/json` (`Encode`/`Decode`)
- routing con metodo + path parameter (`"PUT /todos/{id}"`, `r.PathValue`), feature di `net/http` da Go 1.22+
- conversione stringa→numero con `strconv.Atoi` e gestione errore
- CRUD su slice in memoria (append, ricerca con `for range`, rimozione con `append(s[:i], s[i+1:]...)`)
- status code HTTP (200, 201, 204, 400, 404, 405)
- CORS: origin, preflight `OPTIONS`, middleware che avvolge l'intero router

## Concetti coperti (frontend)

- `useState`/`useEffect`, form controllati, `fetch` (GET/POST/PUT/DELETE)
- interfacce TypeScript per i dati che arrivano dal backend
- aggiornare uno stato-array senza mutarlo (`map`/`filter`/spread)
- styling con Bootstrap (`list-group`, `input-group`, badge, bottoni)

## Stato: completo

- [x] Backend Go (`backend/`): CRUD completo in memoria + CORS
  - `GET /health`
  - `GET /todos`
  - `POST /todos`
  - `PUT /todos/{id}` (toggle done)
  - `DELETE /todos/{id}`
- [x] Frontend React (`frontend/`): crea/leggi/aggiorna/cancella collegato all'API, stile Bootstrap

## Come lanciare

Due terminali:

```
go -C backend run .          # http://localhost:8080
```

```
cd frontend && npm run dev   # http://localhost:5173
```
