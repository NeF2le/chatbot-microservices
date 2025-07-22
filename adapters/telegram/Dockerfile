FROM golang:1.23.6-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o bin/telegram ./cmd

FROM alpine:3.18
RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /app/bin/telegram ./

COPY config.yaml ./

ENTRYPOINT ["./telegram"]