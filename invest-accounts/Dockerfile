FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY invest-accounts/go.mod invest-accounts/go.sum ./
RUN go mod download

COPY invest-accounts/ ./
RUN go build -o invest-accounts ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/invest-accounts/invest-accounts /app/invest-accounts/invest-accounts

EXPOSE 8082

CMD [ "/app/invest-accounts/invest-accounts" ]
