FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go build -o server .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]
