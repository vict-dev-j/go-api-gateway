FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN cd customers && go mod download && go build -o customers ./main.go

RUN cd gateway && go mod download && go build -o gateway ./gateway.go

RUN cd invest-accounts && go mod download && go build -o invest-accounts ./main.go

FROM alpine:latest

COPY --from=builder /app/customers/customers /app/customers
COPY --from=builder /app/gateway/gateway /app/gateway
COPY --from=builder /app/invest-accounts/invest-accounts /app/invest-accounts

EXPOSE 8080 8081 8082

CMD ["/app/customers/customers"]
CMD ["/app/gateway/gateway"]
CMD ["/app/invest-accounts/invest-accounts"]
