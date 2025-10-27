# Go Data Catalog

Простой каталог данных на Go для управления артефактами данных (таблицы, представления, API и т.д.)

## Возможности

- ✅ Полный CRUD для артефактов и контактов
- ✅ Валидация входных данных
- ✅ Логирование запросов
- ✅ Обработка ошибок
- ✅ CORS для работы с фронтендом
- ✅ RESTful API

## Требования

- Go 1.20+
- PostgreSQL 12+

## Установка и запуск

### 1. Клонируйте репозиторий

```bash
git clone <url>
cd go-data-catalog
```

### 2. Установите зависимости

```bash
go mod download
```

### 3. Настройте базу данных

Создайте базу данных PostgreSQL и выполните миграцию:

```bash
psql -U your_user -d your_database < migrations/001_init_schema.sql
```

### 4. Настройте переменные окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
SERVER_PORT=8080
```

### 5. Запустите сервер

```bash
go run cmd/server/main.go
```

Сервер запустится на порту, указанном в `SERVER_PORT` (по умолчанию 8080)

## API Endpoints

### Health Check
- `GET /health` - проверка состояния сервера

### Артефакты
- `GET /api/v1/artifacts` - получить все артефакты
- `GET /api/v1/artifacts/:id` - получить артефакт по ID
- `POST /api/v1/artifacts` - создать новый артефакт
- `PUT /api/v1/artifacts/:id` - обновить артефакт
- `DELETE /api/v1/artifacts/:id` - удалить артефакт

### Контакты
- `GET /api/v1/contacts` - получить все контакты
- `GET /api/v1/contacts/:id` - получить контакт по ID
- `POST /api/v1/contacts` - создать новый контакт
- `PUT /api/v1/contacts/:id` - обновить контакт
- `DELETE /api/v1/contacts/:id` - удалить контакт

## Примеры запросов

### Создание контакта
```bash
curl -X POST http://localhost:8080/api/v1/contacts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Петр Петров",
    "telegram_contact": "@ivanov"
  }'
```

### Создание артефакта
```bash
curl -X POST http://localhost:8080/api/v1/artifacts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "users_table",
    "type": "table",
    "description": "Таблица с пользователями",
    "project_name": "DWH",
    "developer_id": 1
  }'
```

## Типы артефактов

Поддерживаемые типы:
- `table` - таблица БД
- `view` - представление
- `procedure` - процедура
- `function` - функция
- `index` - индекс
- `dataset` - датасет
- `api` - API endpoint
- `file` - файл

## Структура проекта

```
go-data-catalog/
├── cmd/
│   └── server/
│       └── main.go         # Точка входа
├── internal/
│   ├── config/             # Конфигурация
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Middleware (логирование, CORS и т.д.)
│   ├── models/             # Модели данных
│   └── repository/         # Слой работы с БД
│       └── postgres/
├── migrations/             # SQL миграции
├── go.mod
├── go.sum
└── README.md
```

## Разработка

### Запуск в режиме разработки

```bash
# Установите air для hot reload
go install github.com/cosmtrek/air@latest

# Запустите
air
```

### Форматирование кода

```bash
go fmt ./...
```

### Проверка кода

```bash
go vet ./...
```

## TODO

- [ ] Добавить пагинацию
- [ ] Добавить фильтрацию и поиск
- [ ] Добавить Swagger документацию
- [ ] Написать тесты
- [ ] Добавить Docker compose
- [ ] Добавить аутентификацию