FROM golang:1.23-alpine AS builder
RUN apk update && apk add --no-cache \
    make \
    protobuf-dev \
    ca-certificates \
    && update-ca-certificates \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY . .
RUN cd donation-service && go mod download
RUN cd donation-service && GOARCH=amd64 GOOS=linux go build -o main .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/donation-service/main .
COPY --from=builder /app/donation-service/config ./config
COPY --from=builder /app/donation-service/.env ./.env 
EXPOSE 50052

# Ensure binary is executable
RUN chmod +x /app/main
CMD ["./main"]