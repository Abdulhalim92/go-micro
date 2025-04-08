# base go image
FROM golang:1.24-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

RUN chmod +x /app/mailerApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/mailerApp /app
COPY --from=builder /app/templates /app/templates

RUN ls -l /app/templates && sleep 10

CMD ["/app/mailerApp"]