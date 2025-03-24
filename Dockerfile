FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o ./users_app ./cmd/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/users_app .
COPY .env ./

EXPOSE 8080

ENTRYPOINT [ "/app/users_app" ]

