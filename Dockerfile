FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM debian:buster-slim

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.conf ./configs/config.conf
COPY --from=builder /app/web ./web
COPY --from=builder /app/storage ./storage

RUN mkdir -p storage/images storage/thumbnail \
    && chmod -R 777 storage

EXPOSE 8000
CMD ["./main"]