# Variables
BINARY_NAME=openlarq
BUILD_DIR=build
MAIN_PATH=cmd/api/api.go

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build

# Default target
.DEFAULT_GOAL := run

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run:
	@echo "Starting $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Run the application in development mode with hot reload (requires air)
dev:
	@echo "Starting $(BINARY_NAME) in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .
	@echo "Docker image built"

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Docker Compose commands
docker-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping all services..."
	docker-compose down

docker-logs:
	@echo "Showing logs for all services..."
	docker-compose logs -f

docker-rebuild:
	@echo "Rebuilding and starting all services..."
	docker-compose up --build -d

docker-clean:
	@echo "Stopping services and removing containers..."
	docker-compose down
	@echo "Removing unused Docker resources..."
	docker system prune -f
