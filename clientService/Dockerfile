FROM golang:1.20-alpine3.18 AS builder

RUN mkdir app
COPY . /app

WORKDIR /app

RUN go build -o main cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app .

CMD ["/app/main"]
