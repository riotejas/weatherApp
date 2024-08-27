# Simple Makefile

# Build the application
all: build

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Build Docker container
docker-build:
	@docker build --rm -t weatherapp:v1.0 .

# Shutdown DB container
docker-run:
	@docker run -p 8080:8080 weatherapp:v1.0

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main


.PHONY: all build run test clean
