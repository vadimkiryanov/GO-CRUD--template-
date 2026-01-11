
.PHONY: run-all-clean clean-db run-db migrate-up run-server

DB_CONTAINER_NAME = any-name-for-container
DB_PASSWORD = 123456
DB_PORT = 5436
DB_NAME = postgres

run-all-clean: clean-db run-db migrate-up run-server

clean-db:
	@echo "Остановка и удаление контейнера PostgreSQL..."
	-docker stop $(DB_CONTAINER_NAME) || true
	-docker rm $(DB_CONTAINER_NAME) || true

run-db:
	@echo "Запуск контейнера PostgreSQL с Docker..."
	docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d --rm postgres

migrate-up:
	@echo "Запуск миграций..."
	# Wait a bit for the database to be ready
	sleep 5
	migrate -database "postgres://postgres:$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path ./migrations up

run-server: 
	@echo "Запуск сервера..."
	go run ./cmd/main.go
