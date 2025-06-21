# Build Stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o main ./cmd/server

# Run Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
