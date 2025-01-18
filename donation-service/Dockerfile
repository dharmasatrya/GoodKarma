FROM golang:1.23.2

WORKDIR /app

# Copy proto repo first
COPY proto-repo /app/proto-repo

# Copy users service files
COPY goodkarma-donation-service /app/donation-service

WORKDIR /app/donation-service

# Download dependencies
RUN go mod download

# Build
RUN go build -o main .

EXPOSE 50054

CMD ["./main"]