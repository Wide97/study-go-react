# 02 — Auth semplice con JWT

Scopo: costruire una piccola autenticazione full-stack con backend Go e frontend React, senza database, usando un utente finto in memoria.

## Concetti da coprire

### Backend

- handler HTTP per login e rotte protette
- parsing JSON di credenziali (`email`, `password`)
- hash e verifica password con `bcrypt`
- creazione e firma di un JWT
- header `Authorization: Bearer <token>`
- middleware di autenticazione
- gestione status code: `200`, `400`, `401`, `405`
- CORS per frontend Vite

### Frontend

- form controllato per login
- salvataggio token lato client
- stato di autenticazione
- chiamate `fetch` con header `Authorization`
- schermata protetta
- logout

## Stato: in corso

- [x] Backend Go inizializzato (`backend/go.mod`)
- [x] Backend: server HTTP base
- [x] Backend: `POST /login` con verifica credenziali e bcrypt
- [x] Backend: generazione JWT
- [x] Backend: middleware auth
- [x] Backend: endpoint protetto `GET /me`
- [x] Backend: CORS per frontend Vite
- [x] Frontend React inizializzato
- [x] Frontend: login/logout
- [x] Frontend: schermata protetta

## Come lavoriamo

Questo progetto segue la regola del percorso: prima concetti, poi piccoli step guidati, poi revisione del codice scritto.

Non partiremo dalla soluzione completa. Ogni blocco va scritto e capito prima di passare al successivo.
