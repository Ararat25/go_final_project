FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o server.bin ./cmd

FROM alpine:latest

WORKDIR /app/cmd

COPY --from=builder /app/server.bin ./

WORKDIR /app

COPY --from=builder /app/.env ./
COPY --from=builder /app/web ./web

WORKDIR /app/cmd

CMD ["./server.bin"]