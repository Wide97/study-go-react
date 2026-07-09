# 03 — Dashboard realtime

Scopo: costruire una piccola dashboard full-stack in cui un backend Go invia dati in tempo reale e un frontend React li visualizza live.

## Idea del progetto

Il backend simula un dato che cambia nel tempo, per esempio:

- un contatore
- una temperatura finta
- un valore sensore casuale

Il frontend si collega al backend e aggiorna la UI quando arrivano nuovi dati, senza dover premere refresh.

## Concetti da coprire

### Backend

- endpoint HTTP base (`/health`)
- differenza tra HTTP classico e WebSocket
- upgrade della connessione HTTP a WebSocket
- invio periodico di messaggi dal server al client
- JSON per messaggi realtime
- gestione errori e chiusura connessione

### Frontend

- `useEffect` per aprire/chiudere una connessione WebSocket
- `useState` per visualizzare l'ultimo valore ricevuto
- gestione stato realtime che cambia nel tempo
- cleanup dell'effect quando il componente viene smontato
- visualizzazione semplice dei dati ricevuti

## Stato: da iniziare

- [x] Backend Go inizializzato
- [x] Backend: `GET /health`
- [x] Backend: endpoint WebSocket
- [x] Backend: invio dati finti periodici
- [ ] Frontend React inizializzato
- [ ] Frontend: connessione WebSocket
- [ ] Frontend: visualizzazione valore live
- [ ] Frontend: gestione stato connessione

## Metodo

Seguiamo il metodo del percorso: teoria breve, step piccolo, codice scritto da te, revisione insieme.

Non partiamo dalla soluzione completa.
