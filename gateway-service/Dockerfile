FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/gateway-service
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway-service/main .

EXPOSE 8080
CMD ["./main"]