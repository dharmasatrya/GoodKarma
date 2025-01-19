FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/user-service
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service/main .

EXPOSE 50051
CMD ["./main"]