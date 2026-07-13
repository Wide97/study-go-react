# 05 — API con paginazione, filtri e ricerca

Scopo: costruire una piccola app full-stack con backend Go che espone una lista interrogabile tramite query params e frontend React che mostra una tabella filtrabile.

## Idea del progetto

Costruiamo un catalogo ordini/clienti finti.

Il backend espone una lista con:

- ricerca testuale
- filtro per stato
- paginazione
- totale risultati

Il frontend mostra:

- input ricerca
- filtro stato
- tabella risultati
- controlli pagina

## Concetti da coprire

### Backend

- server HTTP base
- endpoint `/health`
- endpoint `GET /orders`
- query params con `r.URL.Query()`
- filtro testuale
- filtro per stato
- paginazione con `page` e `pageSize`
- response JSON con `items`, `total`, `page`, `pageSize`
- CORS per frontend Vite

### Frontend

- fetch con query string
- stato per search/status/page
- tabella dati
- debounce sulla ricerca
- gestione loading/error
- componenti riutilizzabili per filtri e tabella

## Stato: da iniziare

- [x] Backend Go inizializzato
- [x] Backend: `GET /health`
- [x] Backend: dati finti ordini
- [x] Backend: `GET /orders`
- [x] Backend: filtro ricerca
- [x] Backend: filtro stato
- [x] Backend: paginazione
- [x] Backend: CORS per frontend Vite
- [x] Frontend React inizializzato
- [ ] Frontend: fetch ordini
- [ ] Frontend: tabella ordini
- [ ] Frontend: filtri ricerca/stato
- [ ] Frontend: paginazione
- [ ] Frontend: debounce ricerca

## Metodo

Seguiamo il metodo del percorso: teoria breve, step piccolo, codice scritto da te, revisione insieme.

Non partiamo dalla soluzione completa.
