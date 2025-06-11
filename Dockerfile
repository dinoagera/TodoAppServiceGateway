FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/main.go

FROM alpine:3.18

COPY --from=builder /app/main /main

COPY .env .env

EXPOSE 8083

CMD ["/main"]