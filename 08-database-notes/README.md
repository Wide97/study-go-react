# 08 — Database Notes

Scopo: costruire una piccola app full-stack con dati persistenti, usando Go, React e SQLite.

## Idea del progetto

Creiamo una mini app note.

Fino a ora molti progetti hanno usato dati finti in memoria. Qui il punto è introdurre un database vero: le note devono restare salvate anche se il backend viene riavviato.

## API previste

```text
GET    /health
GET    /notes
POST   /notes
PUT    /notes/{id}
DELETE /notes/{id}
```

## Target struttura backend

```text
backend/
  main.go
  models.go
  db.go
  repository.go
  handlers.go
  cors.go
```

## Target struttura frontend

```text
frontend/src/
  App.tsx
  api.ts
  types.ts
  components/
    NoteForm.tsx
    NotesList.tsx
```

## Concetti da coprire

### Backend

- `database/sql`
- driver SQLite
- apertura connessione DB
- `CREATE TABLE IF NOT EXISTS`
- `Query`
- `QueryRow`
- `Exec`
- `Scan`
- repository layer
- gestione errori DB

### Frontend

- fetch CRUD
- form controllati
- create/update/delete
- loading/error
- stato sincronizzato con backend

## Stato: da iniziare

- [x] Scheletro progetto creato
- [x] Backend Go inizializzato
- [x] Backend: `GET /health`
- [x] Backend: collegare SQLite
- [x] Backend: creare tabella `notes`
- [ ] Backend: modello `Note`
- [ ] Backend: repository `ListNotes`
- [ ] Backend: repository `CreateNote`
- [ ] Backend: repository `UpdateNote`
- [ ] Backend: repository `DeleteNote`
- [ ] Backend: handler `GET /notes`
- [ ] Backend: handler `POST /notes`
- [ ] Backend: handler `PUT /notes/{id}`
- [ ] Backend: handler `DELETE /notes/{id}`
- [ ] Backend: CORS
- [x] Frontend React inizializzato
- [ ] Frontend: types
- [ ] Frontend: api client
- [ ] Frontend: fetch note
- [ ] Frontend: form creazione
- [ ] Frontend: lista note
- [ ] Frontend: delete nota
- [ ] Frontend: edit nota

## Metodo

Qui procediamo con calma: prima backend e database, poi frontend.

La regola principale è distinguere bene:

- handler HTTP: legge request e scrive response;
- repository: parla col database;
- model: descrive i dati.
