FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /api

EXPOSE 8080

CMD ["/api"]
