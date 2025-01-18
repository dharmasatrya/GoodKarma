FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/event-service
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/event-service/main .

EXPOSE 50055
CMD ["./main"]