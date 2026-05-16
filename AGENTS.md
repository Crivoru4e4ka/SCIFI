# AGENTS.md — Scifi (Task Management System)

## Общая информация о проекте

Scifi — это веб-приложение для управления задачами и научными проектами, вдохновленное Jira, с поддержкой Kanban-доски, управления гипотезами, датасетами и генерацией отчетов по ГОСТ 7.32. Проект разрабатывается в рамках дипломной работы.

- **Язык комментариев и документации:** русский
- **Язык интерфейса:** русский
- **Go module:** `project-MVP`
- **Go version:** 1.25.1

## Технологический стек

| Компонент | Технология |
|-----------|-----------|
| Backend | Go 1.25.1 |
| HTTP-роутер | Gorilla Mux v1.8.1 |
| База данных | PostgreSQL |
| Драйвер БД | `lib/pq` v1.12.0 |
| Хеширование паролей | `golang.org/x/crypto/bcrypt` |
| Переменные окружения | `joho/godotenv` v1.5.1 |
| Frontend | Vue.js 3 (inline в HTML), Bootstrap 5.3.0, Font Awesome 6.0.0 |
| Форматирование | Prettier 3.8.3 (devDependency) |

## Структура проекта

```
project-MVP/
├── main.go                 # Точка входа: загружает .env, подключается к БД, запускает сервер на :8080
├── go.mod / go.sum         # Go-модули
├── package.json            # Node devDependencies (только Prettier)
├── .env                    # Переменные окружения (НЕ коммитится!)
├── .gitignore              # Стандартный gitignore для Go + Node
│
├── db/
│   ├── db.go               # Инициализация подключения к PostgreSQL через переменные окружения
│   └── schema.sql          # Полная схема БД (таблицы, внешние ключи)
│
├── models/                 # DTO / структуры данных
│   ├── user.go
│   ├── project.go
│   ├── task.go
│   ├── comment.go
│   ├── attachment.go
│   ├── sprint.go
│   ├── tag.go
│   ├── tasktag.go
│   ├── task_history.go
│   ├── team.go
│   ├── project_member.go
│   ├── hypothesis.go
│   ├── dataset.go
│   └── experimentdataset.go
│
├── handlers/               # HTTP-обработчики (handlers layer)
│   ├── auth_handler.go     # Регистрация, вход, выход (сессионные куки)
│   ├── user_handler.go     # CRUD пользователей, /me
│   ├── project_handler.go  # CRUD проектов, прогресс
│   ├── task_handler.go     # CRUD задач, обновление статуса, спринты, гипотезы
│   ├── comment_handler.go  # Комментарии к задачам
│   ├── attachment_handler.go # Загрузка файлов
│   ├── project_member_handler.go # Участники проекта
│   ├── sprint_handler.go   # Управление спринтами
│   ├── team_handler.go     # Управление командами и их участниками
│   ├── tag_handler.go      # Теги
│   ├── task_tag_handler.go # Связь задач и тегов
│   ├── dataset_handler.go  # Датасеты
│   ├── experiment_dataset_handler.go
│   └── ui_handler.go       # Раздача HTML-страниц и редиректы (index, login)
│
├── services/               # Бизнес-логика (service layer)
│   ├── user_service.go     # Аутентификация, bcrypt, роли
│   ├── project_service.go  # Логика проектов, транзакции
│   ├── task_service.go     # Логика задач, история изменений, теги
│   ├── comment_service.go
│   ├── attachment_service.go
│   ├── project_member_service.go
│   ├── sprint_service.go
│   ├── team_service.go
│   ├── tag_service.go
│   ├── task_tag_service.go
│   ├── dataset_service.go
│   └── experiment_service.go
│
├── routes/
│   └── routes.go           # Определение всех маршрутов и middleware (NoCacheMiddleware)
│
├── web/                    # Frontend-файлы
│   ├── index.html          # Основное SPA-приложение на Vue.js (~3900 строк)
│   └── login.html          # Страница входа/регистрации
│
├── static/                 # Статические файлы
│   ├── uploads/            # Загруженные пользователями файлы
│   └── *.png               # Бейджи и логотипы
│
└── migrations/             # SQL-миграции (применяются вручную через psql)
    ├── 2026-05-10_add_new_tables_and_fields.sql
    ├── add_task_fields.sql
    └── fix_task_history.sql
```

## Сборка и запуск

### Предварительные требования
- Go 1.25.1+
- PostgreSQL 12+
- Node.js (только для Prettier, опционально)

### Настройка окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=ваш_пароль
DB_NAME=SRA
```

**Важно:** файл `.env` добавлен в `.gitignore` и никогда не должен попадать в репозиторий.

### Инициализация базы данных

```bash
# Создание таблиц
psql -h localhost -U postgres -d SRA -f db/schema.sql

# Применение миграций
psql -h localhost -U postgres -d SRA -f migrations/2026-05-10_add_new_tables_and_fields.sql
psql -h localhost -U postgres -d SRA -f migrations/add_task_fields.sql
psql -h localhost -U postgres -d SRA -f migrations/fix_task_history.sql
```

### Установка зависимостей

```bash
go mod download
```

### Запуск приложения

```bash
go run main.go
```

Приложение доступно по адресу: `http://localhost:8080`

### Сборка бинарника

```bash
go build -o app.exe main.go
```

## Архитектура приложения

### Трехслойная архитектура

1. **Handlers** (`handlers/`) — принимают HTTP-запросы, валидируют входные данные, возвращают HTTP-ответы. Не содержат бизнес-логики.
2. **Services** (`services/`) — содержат всю бизнес-логику, транзакции БД, валидацию бизнес-правил.
3. **Models** (`models/`) — чистые структуры данных с JSON-тегами.
4. **DB** (`db/`) — подключение к PostgreSQL и глобальная переменная `DB *sql.DB`.

### Аутентификация

- Сессионная аутентификация через HTTP-only куки с именем `session`
- Значение куки — строковое представление `user_id`
- Пароли хешируются через `bcrypt`
- Middleware `NoCacheMiddleware` запрещает кэширование страниц для защиты личного кабинета

### Роли пользователей

- `admin` — системный администратор
- `user` — обычный сотрудник/студент
- `guest` — внешний эксперт/рецензент

При создании проекта создатель автоматически становится `manager` через таблицу `project_members`.

## Основные сущности БД

- **users** — пользователи системы
- **projects** — проекты/научные разделы (с полями цели, гипотезы, новизны, ожидаемого результата)
- **tasks** — задачи (с поддержкой спринтов, гипотез, тегов, research_method, research_contribution)
- **comments** — комментарии к задачам
- **attachments** — прикрепленные файлы
- **sprints** — спринты внутри проектов
- **tags / task_tags** — система тегов
- **teams / team_members** — команды и их участники
- **project_members** — участники проектов с ролями
- **hypotheses** — научные гипотезы в рамках проекта
- **datasets / experiment_datasets** — датасеты и их связь с задачами
- **task_history** — аудит изменений задач (особенно статусов)

## API Endpoints (основные)

### Аутентификация
- `POST /auth/register` — регистрация
- `POST /auth/login` — вход (устанавливает куку session)
- `GET /logout` — выход

### Пользователи
- `GET /users`, `POST /users`
- `GET /me` — текущий пользователь

### Проекты
- `GET /projects`, `POST /projects`
- `GET /user/projects` — проекты текущего пользователя
- `GET /projects/{id}/tasks`
- `GET /projects/{id}/progress`
- `GET /projects/{id}/hypotheses`, `POST /projects/{id}/hypotheses`
- `GET /projects/{id}/report` — отчет по ГОСТ 7.32

### Задачи
- `POST /tasks`
- `PATCH /tasks/{id}/status` — обновление статуса с записью в task_history
- `POST /tasks/{id}/comments`, `GET /tasks/{id}/comments`
- `PATCH /tasks/{id}/sprint`
- `GET /user/tasks` — все задачи пользователя

### Спринты
- `GET /projects/{id}/sprints`, `POST /projects/{id}/sprints`
- `PATCH /sprints/{id}/start`, `PATCH /sprints/{id}/complete`

### Команды
- `GET /user/{id}/teams`
- `POST /teams`, `PATCH /teams/{id}`, `DELETE /teams/{id}`
- `GET /teams/{id}/members`, `POST /teams/{id}/members`
- `DELETE /teams/{id}/members/{userID}`, `PATCH /teams/{id}/members/{userID}/role`

## Стиль кода и соглашения

### Язык комментариев
Весь код комментируется на **русском языке**. При внесении изменений сохраняйте этот стиль.

### Именование
- Экспортируемые функции/типы: `PascalCase`
- Локальные переменные: `camelCase`
- Структуры моделей: `PascalCase` с JSON-тегами в `snake_case`
- Файлы: `snake_case.go`

### Обработка ошибок
- В сервисах определены публичные переменные ошибок: `ErrNotFound`, `ErrInvalidCredentials`, `ErrInvalidRole`, `ErrDuplicateUser`
- Хендлеры проверяют конкретные ошибки и возвращают соответствующие HTTP-коды
- Всегда используйте `defer rows.Close()` после `db.DB.Query()`

### SQL-запросы
- Параметризованные запросы через `$1`, `$2` (стиль PostgreSQL)
- Для nullable полей используйте `sql.NullInt64`, `sql.NullString` и конвертируйте в указатели

### Фронтенд
- Весь фронтенд реализован внутри двух HTML-файлов (`web/index.html`, `web/login.html`) как Vue.js 3 SPA без сборки
- Bootstrap 5.3.0 и Font Awesome 6.0.0 подключаются через CDN
- Кастомные стили пишутся в `<style>` блоке внутри HTML

## Безопасность

- Пароли хешируются через `bcrypt` с `bcrypt.DefaultCost`
- Сессионные куки: `HttpOnly: true`, `SameSite: Lax`, срок действия 24 часа
- `NoCacheMiddleware` предотвращает кэширование приватных страниц в браузере
- Параметризованные SQL-запросы защищают от SQL-инъекций
- Файл `.env` исключен из Git

## Тестирование

На текущий момент в проекте **отсутствуют автоматические тесты**. Тестирование выполняется вручную.

## Деплой

Приложение представляет собой одиночный Go-бинарник, который:
1. Слушает порт `:8080`
2. Раздает статические файлы из `./static/`
3. Раздает HTML-страницы из `./web/`
4. Обрабатывает REST API

Для production-деплоя достаточно скомпилировать бинарник и установить переменные окружения на сервере (либо через `.env`, либо через системные env vars).

## Полезные файлы для ознакомления

- `README.md` — пользовательская документация с примерами API-запросов
- `KANBAN_GUIDE.md` — документация по Kanban-доске
- `db/schema.sql` — полная схема базы данных
- `migrations/*.sql` — история изменений схемы
