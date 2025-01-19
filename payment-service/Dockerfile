FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/payment-service
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-service/main .

EXPOSE 50053
CMD ["./main"]