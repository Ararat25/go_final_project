FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server.bin ./cmd

FROM ubuntu:latest

WORKDIR /app/cmd

EXPOSE 8080/tcp

COPY --from=builder /app/server.bin ./

WORKDIR /app

COPY --from=builder /app/.env ./
COPY --from=builder /app/web ./web

WORKDIR /app/cmd

CMD ["./server.bin"]