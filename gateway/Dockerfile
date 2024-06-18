FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

COPY gateway/ ./
RUN go build -o gateway ./gateway.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway/gateway /app/gateway/gateway

EXPOSE 8081

CMD [ "/app/gateway/gateway" ]
