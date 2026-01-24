
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
		echo "Контейнер $(DB_CONTAINER_NAME) уже запущен..."; \
	else \
		if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
			echo "Контейнер $(DB_CONTAINER_NAME) остановлен, запускаем его..."; \
			docker start $(DB_CONTAINER_NAME); \
		else \
			echo "Контейнер $(DB_CONTAINER_NAME) не существует, создаем новый..."; \
			make run-db; \
		fi \
	fi

clean-db:
	@echo "Остановка и удаление контейнера PostgreSQL..."
	-docker stop $(DB_CONTAINER_NAME) || true
	-docker rm $(DB_CONTAINER_NAME) || true
run-db:
	@echo "Проверка существующего контейнера..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
		echo "Контейнер $(DB_CONTAINER_NAME) уже существует, останавливаем его..."; \
		docker stop $(DB_CONTAINER_NAME); \
		docker rm $(DB_CONTAINER_NAME); \
	fi
	@echo "Запуск контейнера PostgreSQL с Docker..."
	docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d postgres

run-db-with-delete:
	@echo "Проверка существующего контейнера..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
		echo "Контейнер $(DB_CONTAINER_NAME) уже существует, останавливаем его..."; \
		docker stop $(DB_CONTAINER_NAME); \
		docker rm $(DB_CONTAINER_NAME); \
	fi
	@echo "Запуск контейнера PostgreSQL с Docker..."
	docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d --rm postgres

migrate-up:
	@echo "Запуск миграций..."
	@echo "Ожидание готовности базы данных..."
	@timeout 30 bash -c 'until docker exec $(DB_CONTAINER_NAME) pg_isready > /dev/null 2>&1; do sleep 1; done' || echo "База данных может быть недоступна"
	@sleep 2  # Additional wait time for full readiness
	migrate -database "postgres://postgres:$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path ./migrations up

run-server: 
	@echo "Запуск сервера..."
	go run ./cmd/main.go
