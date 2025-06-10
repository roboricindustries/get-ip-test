FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY src/ ./src/

WORKDIR /app/src

RUN go build -o api

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache curl

COPY --from=builder /app/src/api .

CMD ["./api"]
