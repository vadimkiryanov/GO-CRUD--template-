
.PHONY: run-all-clean clean-db run-db migrate-up run-server

DB_CONTAINER_NAME = any-name-for-container
DB_PASSWORD = 123456
DB_PORT = 5436
DB_NAME = postgres

run-all-clean: clean-db run-db-with-delete migrate-up run-server
run-all: run-db-check migrate-up run-server

# Check if DB container exists and start it if needed
run-db-check:
	@if docker ps -q -f name=$(DB_CONTAINER_NAME) | grep -q . ; then \
		echo "Container $(DB_CONTAINER_NAME) is already running..."; \
	else \
		if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
			echo "Container $(DB_CONTAINER_NAME) is stopped, starting..."; \
			docker start $(DB_CONTAINER_NAME); \
		else \
			echo "Container $(DB_CONTAINER_NAME) does not exist, creating new..."; \
			make run-db; \
		fi \
	fi

clean-db:
	@echo "Stopping and removing PostgreSQL container..."
	-docker stop $(DB_CONTAINER_NAME) || true
	-docker rm $(DB_CONTAINER_NAME) || true
run-db:
	@echo "Checking existing container..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
		echo "Container $(DB_CONTAINER_NAME) already exists, stopping..."; \
		docker stop $(DB_CONTAINER_NAME); \
		docker rm $(DB_CONTAINER_NAME); \
	fi
	@echo "Starting PostgreSQL container with Docker..."
	docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d postgres

run-db-with-delete:
	@echo "Checking existing container..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
		echo "Container $(DB_CONTAINER_NAME) already exists, stopping..."; \
		docker stop $(DB_CONTAINER_NAME); \
		docker rm $(DB_CONTAINER_NAME); \
	fi
	@echo "Starting PostgreSQL container with Docker..."
	docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d --rm postgres

migrate-up:
	@echo "Running migrations..."
	@echo "Waiting for database readiness..."
	@timeout 30 bash -c 'until docker exec $(DB_CONTAINER_NAME) pg_isready > /dev/null 2>&1; do sleep 1; done' || echo "Database may be unavailable"
	@sleep 2  # Additional wait time for full readiness
	migrate -database "postgres://postgres:$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path ./migrations up

run-server: 
	@echo "Starting server..."
	go run ./cmd/main.go
