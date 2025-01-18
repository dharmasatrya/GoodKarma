# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the entire monorepo for proto access
COPY . .

# Build the service
WORKDIR /app/donation-service
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/donation-service/main .

EXPOSE 50052
CMD ["./main"]