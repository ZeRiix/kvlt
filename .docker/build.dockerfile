FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN mkdir ./data

ENV CLEANER_TIME="@every 3s"

EXPOSE 8080

RUN chmod +x /root/main

CMD ["./main"]