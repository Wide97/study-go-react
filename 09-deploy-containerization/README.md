# 09 — Deploy Containerization

Scopo: prendere una piccola app full-stack Go + React e farla girare con Docker.

## Idea del progetto

Costruiamo un backend Go minimo e un frontend React minimo, poi li containerizziamo.

Il punto non è creare nuove feature applicative: il focus è capire come si prepara un progetto per girare fuori dal dev server locale.

## Target finale

```text
backend container  -> espone API HTTP
frontend container -> serve la build React
docker compose     -> avvia tutto insieme
```

## Concetti da coprire

### Backend

- build di un binario Go
- Dockerfile per Go
- differenza tra `go run` e binario compilato
- porta esposta dal container
- variabili d'ambiente minime

### Frontend

- build React con Vite
- differenza tra dev server e build statica
- servire file statici con nginx
- configurazione API base URL

### Docker

- immagine
- container
- Dockerfile
- `.dockerignore`
- multi-stage build
- `docker compose`
- network tra servizi

## Struttura prevista

```text
09-deploy-containerization/
  README.md
  docker-compose.yml
  backend/
    Dockerfile
    .dockerignore
    main.go
    go.mod
  frontend/
    Dockerfile
    .dockerignore
    nginx.conf
    package.json
    src/
```

## Stato

- [x] Scheletro progetto creato
- [x] Backend Go inizializzato
- [x] Backend: `GET /health`
- [ ] Backend: Dockerfile
- [ ] Backend: build immagine
- [ ] Frontend React inizializzato
- [ ] Frontend: chiamata API backend
- [ ] Frontend: Dockerfile
- [ ] Frontend: nginx config
- [ ] Frontend: build immagine
- [ ] Docker Compose
- [ ] Avvio full-stack con `docker compose up`
- [ ] Aggiornare `STUDY.md`

## Metodo

Procediamo a piccoli passi.

Prima facciamo funzionare backend e frontend in locale, poi li mettiamo nei container. Docker arriva dopo che sappiamo cosa deve eseguire.
