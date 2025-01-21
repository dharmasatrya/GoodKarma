# Build stage
FROM golang:1.21-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy the entire monorepo for proto access
COPY . .

# Download dependencies
RUN cd donation-service && go mod download

# Build the service
WORKDIR /app/donation-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:3.19

# Install certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/donation-service/main .
# Copy any config files if needed
COPY --from=builder /app/donation-service/config ./config
# Copy proto files (since you're mounting them as volume, this might be redundant)
COPY --from=builder /app/proto-repo ./proto-repo

EXPOSE 50052
CMD ["./main"]