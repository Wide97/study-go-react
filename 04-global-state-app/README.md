# 04 — App con stato globale complesso

Scopo: costruire una piccola app full-stack con più entità collegate, per studiare stato globale lato React e API Go semplici.

## Idea del progetto

Costruiamo un mini carrello:

- prodotti disponibili
- carrello con quantità
- totale calcolato
- azioni di aggiunta/rimozione/modifica quantità

Il backend espone i dati iniziali e, per ora, non usa database.

## Concetti da coprire

### Backend

- server HTTP base
- endpoint `/health`
- endpoint `GET /products`
- struct Go per prodotti
- JSON encode
- CORS per frontend Vite

### Frontend

- stato locale vs stato globale
- Context API
- provider
- custom hook per leggere/modificare il carrello
- derived state: totale e quantità
- componenti separati per lista prodotti e carrello

## Stato: da iniziare

- [x] Backend Go inizializzato
- [x] Backend: `GET /health`
- [x] Backend: `GET /products`
- [x] Backend: CORS per frontend Vite
- [ ] Frontend React inizializzato
- [ ] Frontend: lista prodotti
- [ ] Frontend: Context per carrello
- [ ] Frontend: aggiunta/rimozione/modifica quantità
- [ ] Frontend: totale carrello

## Metodo

Seguiamo il metodo del percorso: teoria breve, step piccolo, codice scritto da te, revisione insieme.

Non partiamo dalla soluzione completa.
