# EM-test: Сервис управления подписками

REST API для управления онлайн-подписками пользователей с подсчётом суммарной стоимости.

## Требования

- Go 1.24+
- Docker и Docker Compose (для контейнеризации)
- PostgreSQL (можно запускать локально или через Docker)

## Установка и запуск

### 1. Клонирование репозитория
git clone https://github.com/yourusername/em-test.git
cd em-test

### 2. Конфигурация
Скопируйте пример файла .env.example в .env и отредактируйте параметры подключения к базе данных и порта:
DATABASE_URL=postgres://user:password@localhost:5432/DBname?sslmode=disable
SUBSCRIPTION_PORT=:8080

### 3. Запуск миграций
Если используете локальный PostgreSQL, примените миграции вручную из:
em-test/cmd/internal/migrations

### 4. Запуск сервиса локально
go run cmd/main.go

### 5. Запуск с Docker Compose
В директории с docker-compose.yml:
docker-compose up --build

Это запустит сервис и PostgreSQL в контейнерах.

## API
Документация по API доступна через Swagger:
http://localhost:8080/swagger/index.html

## Тестирование
Запуск тестов:
go test ./cmd/internal/tests/...

## Логирование
Сервис ведёт логи в stdout с таймстампом: ошибки и важные события.

### Контакты и поддержка
Для вопросов и помощи обращайтесь в Issues репозитория.

## Кратко о работе
- CRUDL-эндпоинты для подписок (/subscriptions)
- Отчёт по суммарной стоимости за месяц с фильтрами (/subscriptions/report)
- Поддержка фильтров по пользователю и провайдеру