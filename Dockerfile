FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o go-file-parsing .

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/go-file-parsing .

# Copy the config file and sample CSV
COPY --from=builder /app/config.json .
COPY --from=builder /app/sample.csv .

# Create a directory for data files
RUN mkdir -p /app/data

# Set the environment variable for Valkey
ENV VALKEY_URLS="valkey:6379"

# Command to run the application
ENTRYPOINT ["/app/go-file-parsing"]
