FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/notification-service
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notification-service/main .

EXPOSE 50054
CMD ["./main"]