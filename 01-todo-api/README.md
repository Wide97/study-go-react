# 01 — Todo API

Scopo: primo progetto del percorso di studio, per prendere confidenza con Go da zero e poi collegarci un frontend React.

## Concetti coperti (backend)

- `net/http`: handler, `ListenAndServe`, gestione errori con `log.Fatal`
- struct e JSON tag, `encoding/json` (`Encode`/`Decode`)
- routing con metodo + path parameter (`"PUT /todos/{id}"`, `r.PathValue`), feature di `net/http` da Go 1.22+
- conversione stringa→numero con `strconv.Atoi` e gestione errore
- CRUD su slice in memoria (append, ricerca con `for range`, rimozione con `append(s[:i], s[i+1:]...)`)
- status code HTTP (200, 201, 204, 400, 404, 405)

## Stato

- [x] Backend Go (`backend/`): CRUD completo in memoria
  - `GET /health`
  - `GET /todos`
  - `POST /todos`
  - `PUT /todos/{id}` (toggle done)
  - `DELETE /todos/{id}`
- [ ] Frontend React: da fare

## Come lanciare il backend

```
go -C backend run .
```

Server su `http://localhost:8080`.
