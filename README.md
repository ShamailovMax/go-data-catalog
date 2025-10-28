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
JWT_SECRET=changeme_super_secret
TOKEN_TTL=60
```

### 5. Запустите сервер

```bash
go run cmd/server/main.go
```

Сервер запустится на порту, указанном в `SERVER_PORT` (по умолчанию 8080)

## API Endpoints

### Health Check
- `GET /health` - проверка состояния сервера

### Аутентификация
- `POST /api/v1/auth/register` — регистрация (email, password, name) → JWT токен
- `POST /api/v1/auth/login` — логин → JWT токен

Все ниже — под Bearer JWT.

### Команды (рабочие пространства)
- `GET /api/v1/teams?search=<q>` — поиск команд
- `POST /api/v1/teams` — создать команду (создатель становится owner)
- `POST /api/v1/teams/:teamId/join` — запрос на вступление
- `GET /api/v1/teams/:teamId/requests` — запросы на вступление (owner/admin)
- `POST /api/v1/teams/:teamId/requests/:id/(approve|reject)` — решение по запросу (owner/admin)
- `GET /api/v1/me/teams` — мои команды

### Артефакты (в контексте команды)
- `GET /api/v1/teams/:teamId/artifacts`
- `GET /api/v1/teams/:teamId/artifacts/:id`
- `POST /api/v1/teams/:teamId/artifacts`
- `PUT /api/v1/teams/:teamId/artifacts/:id`
- `DELETE /api/v1/teams/:teamId/artifacts/:id`

Поля артефактов:
- `GET /api/v1/teams/:teamId/artifacts/:id/fields`
- `POST /api/v1/teams/:teamId/artifacts/:id/fields`
- `GET /api/v1/teams/:teamId/fields/:id`
- `PUT /api/v1/teams/:teamId/fields/:id`
- `DELETE /api/v1/teams/:teamId/fields/:id`

### Контакты (в контексте команды)
- `GET /api/v1/teams/:teamId/contacts`
- `GET /api/v1/teams/:teamId/contacts/:id`
- `POST /api/v1/teams/:teamId/contacts`
- `PUT /api/v1/teams/:teamId/contacts/:id`
- `DELETE /api/v1/teams/:teamId/contacts/:id`

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
- [x] Добавить аутентификацию
- [x] Разграничение по командам (multi-tenant)
