FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN cd user-service && go mod download
RUN cd user-service && CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/user-service/main .
COPY --from=builder /app/user-service/config ./config
COPY --from=builder /app/payment-service/.env ./.env 
EXPOSE 50051
CMD ["./main"]