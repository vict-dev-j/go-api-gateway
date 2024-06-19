FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY customers/go.mod customers/go.sum ./
RUN go mod download

COPY customers/ ./
RUN go build -o customers ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/customers/customers /app/customers/customers

ENV POSTGRES_HOST=localhost \
    POSTGRES_PORT=5432 \
    POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=postgres \
    POSTGRES_DB=customers

EXPOSE 8080

CMD [ "/app/customers/customers" ]
