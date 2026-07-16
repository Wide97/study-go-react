# Primo stage: qui usiamo Node solo per trasformare il codice TypeScript/React
# nella cartella statica dist. Il server finale non ha bisogno di Node.
FROM node:24-alpine AS builder

WORKDIR /app

# Copiare prima i manifest permette a Docker di riusare la cache delle
# dipendenze quando cambiano solo i file sorgenti. npm ci usa package-lock.json
# e installa esattamente le versioni gia verificate in locale.
COPY package.json package-lock.json ./
RUN npm ci

COPY . ./

# Le variabili VITE_ vengono lette e incorporate nella build da Vite.
# In Compose il valore /api funziona grazie al reverse proxy nginx.
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# Secondo stage: nginx serve soltanto file statici, quindi l'immagine finale
# e molto piu piccola di una che contiene anche Node e le dipendenze di build.
FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
