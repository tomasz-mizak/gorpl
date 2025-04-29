.PHONY: build test clean run docker-build docker-run

# Variables
BINARY_NAME=gorpl
DOCKER_IMAGE=gorpl:latest

# Build the application
build:
	go build -o $(BINARY_NAME) ./cmd/...

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run the application
run:
	go run ./cmd/...

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Default target
all: deps fmt vet test build 