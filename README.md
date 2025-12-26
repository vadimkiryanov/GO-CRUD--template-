# 1. Инициализация нового модуля (пропустить этот шаг, если вы уже имеете свой проект)
go mod init your-module-name

# 2. Скачать все зависимости
go mod download

# 3. Установить зависимости и обновить go.mod и go.sum
go mod tidy

## 4. Настройка базы данных

Для запуска и работы с базой данных PostgreSQL выполните следующие шаги:

1.  **Запуск PostgreSQL с Docker (если используется):**
    ```bash
    docker run --name=any-name-for-container -e POSTGRES_PASSWORD=123456 -p 5436:5432 -d --rm postgres
    ```

2.  **Создание новой базы данных (миграции):**
    *Создание миграций:*
    ```bash
    migrate create -ext sql -dir ./migrations -seq init
    ```

    *Запуск миграций, применение их к базе данных:*
    ```bash
    migrate -database "postgres://postgres:123456@localhost:5436/postgres?sslmode=disable" -path ./migrations up
    ```
3.  **Подключение к созданной базе данных при помощи docker (опционально):**
    ```bash
    docker exec -it any-name-for-container psql -U postgres
    ```

