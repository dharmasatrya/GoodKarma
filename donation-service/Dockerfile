FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN cd donation-service && go mod download
RUN cd donation-service && CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/donation-service/main .
COPY --from=builder /app/donation-service/config ./config
COPY --from=builder /app/donation-service/.env ./.env 
EXPOSE 50052
CMD ["./main"]