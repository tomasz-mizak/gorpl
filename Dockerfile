# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy templates
COPY --from=builder /app/templates ./templates

# Create directories for static files and XML data
RUN mkdir -p /app/static /app/data

# Expose the application port
EXPOSE 1532

# Run the application
CMD ["./main"] 