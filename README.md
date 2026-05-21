# Scifi — Task Management System

Современная веб-система управления задачами и научными проектами, вдохновлённая Jira. Разработана в рамках дипломной работы.

## 🎯 Назначение

Scifi предназначена для организации работы над научными проектами: планирование задач, управление гипотезами, отслеживание прогресса, работа с датасетами и генерация отчётов по ГОСТ 7.32. Система поддерживает Kanban-доску, спринты, команды и гибкую систему ролей.

## 🚀 Особенности

- ✅ **Задачи и Kanban-доска** — полный жизненный цикл задач с drag-and-drop
- ✅ **Управление проектами** — создание проектов с целями, гипотезами, новизной
- ✅ **Гипотезы** — научные гипотезы в рамках проекта
- ✅ **Датасеты и эксперименты** — привязка данных к задачам
- ✅ **Спринты** — планирование итераций внутри проектов
- ✅ **Команды** — управление командами и их участниками
- ✅ **Система ролей (RBAC)** — роли admin, user, guest; роли в проектах и командах
- ✅ **Комментарии и вложения** — обсуждение задач с файлами
- ✅ **История изменений** — аудит изменений задач (task_history)
- ✅ **Теги** — гибкая система тегов для задач
- ✅ **Приоритизация** — 5 уровней приоритета
- ✅ **Отчёты по ГОСТ 7.32** — автоматическая генерация научных отчётов
- ✅ **Гранты и финансирование** — учёт грантов для проектов
- ✅ **Тёмная тема** — переключение между светлой и тёмной темой

## 🛠 Технологический стек

| Компонент | Технология |
|-----------|------------|
| Backend | Go 1.25.1 |
| HTTP-роутер | Gorilla Mux v1.8.1 |
| База данных | PostgreSQL |
| Драйвер БД | `lib/pq` v1.12.0 |
| Хеширование паролей | `golang.org/x/crypto/bcrypt` |
| Переменные окружения | `joho/godotenv` v1.5.1 |
| Frontend | Vue.js 3 (inline в HTML), Bootstrap 5.3.0, Font Awesome 6.0.0 |
| Форматирование | Prettier 3.8.3 (devDependency) |

## 📋 Структура проекта

```
project-MVP/
├── main.go                    # Точка входа: загружает .env, подключается к БД, запускает сервер на :8080
├── go.mod / go.sum            # Go-модули
├── package.json               # Node devDependencies (только Prettier)
├── .env                       # Переменные окружения (НЕ коммитится!)
├── .gitignore                 # Стандартный gitignore для Go + Node
│
├── db/
│   ├── db.go                  # Инициализация подключения к PostgreSQL
│   └── schema.sql             # Полная схема БД (таблицы, внешние ключи)
│
├── models/                    # DTO / структуры данных
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
│   ├── experimentdataset.go
│   └── grant.go
│
├── handlers/                  # HTTP-обработчики
│   ├── auth_handler.go        # Регистрация, вход, выход (сессионные куки)
│   ├── user_handler.go        # CRUD пользователей, /me
│   ├── project_handler.go     # CRUD проектов, прогресс, отчёты
│   ├── task_handler.go        # CRUD задач, обновление статуса, спринты, гипотезы
│   ├── comment_handler.go     # Комментарии к задачам
│   ├── attachment_handler.go  # Загрузка файлов
│   ├── project_member_handler.go
│   ├── sprint_handler.go
│   ├── team_handler.go
│   ├── tag_handler.go
│   ├── task_tag_handler.go
│   ├── dataset_handler.go
│   ├── experiment_dataset_handler.go
│   ├── grant_handler.go
│   └── ui_handler.go          # Раздача HTML-страниц
│
├── services/                  # Бизнес-логика
│   ├── user_service.go
│   ├── project_service.go
│   ├── task_service.go
│   └── ...
│
├── routes/
│   └── routes.go              # Маршруты API и middleware
│
├── middleware/
│   └── auth.go                # Аутентификация и авторизация
│
├── web/                       # Frontend
│   ├── index.html             # Основное SPA-приложение на Vue.js
│   └── login.html             # Страница входа/регистрации
│
├── static/                    # Статические файлы
│   ├── uploads/               # Загруженные пользователями файлы
│   └── *.png                  # Бейджи и логотипы
│
└── migrations/                # SQL-миграции (применяются вручную через psql)
    ├── 2026-05-10_add_new_tables_and_fields.sql
    ├── 2026-05-16_add_grants_and_funding.sql
    ├── 2026-05-16_add_rbac_system.sql
    ├── 2026-05-16_fix_users_role_check.sql
    ├── 2026-05-16_link_projects_to_teams.sql
    ├── 2026-05-16_unify_team_roles.sql
    ├── add_task_doi_column.sql
    ├── add_task_fields.sql
    └── fix_task_history.sql
```

## 📦 Требования

- Go 1.25.1+
- PostgreSQL 12+
- Node.js (опционально, только для Prettier)
- Веб-браузер с поддержкой ES6

## 🔧 Установка и настройка

### 1. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=ваш_пароль
DB_NAME=SRA
```

**⚠️ ВАЖНО:** Файл `.env` добавлен в `.gitignore` и никогда не должен попадать в репозиторий!

### 2. Инициализация базы данных

```bash
# Создание таблиц
psql -h localhost -U postgres -d SRA -f db/schema.sql

# Применение миграций
psql -h localhost -U postgres -d SRA -f migrations/2026-05-10_add_new_tables_and_fields.sql
psql -h localhost -U postgres -d SRA -f migrations/2026-05-16_add_grants_and_funding.sql
psql -h localhost -U postgres -d SRA -f migrations/2026-05-16_add_rbac_system.sql
psql -h localhost -U postgres -d SRA -f migrations/2026-05-16_fix_users_role_check.sql
psql -h localhost -U postgres -d SRA -f migrations/2026-05-16_link_projects_to_teams.sql
psql -h localhost -U postgres -d SRA -f migrations/2026-05-16_unify_team_roles.sql
psql -h localhost -U postgres -d SRA -f migrations/add_task_doi_column.sql
psql -h localhost -U postgres -d SRA -f migrations/add_task_fields.sql
psql -h localhost -U postgres -d SRA -f migrations/fix_task_history.sql
```

### 3. Установка зависимостей Go

```bash
go mod download
```

### 4. Запуск приложения

```bash
go run main.go
```

Приложение будет доступно по адресу: `http://localhost:8080`

### 5. Сборка бинарника (опционально)

```bash
go build -o app.exe main.go
```

## 🎯 Использование

### Регистрация и вход

1. Откройте `http://localhost:8080`
2. Зарегистрируйтесь или войдите с существующим аккаунтом
3. Сессия сохраняется в HTTP-only cookie на 24 часа

### Создание проекта

1. Откройте раздел **"Разделы"** в левой панели
2. Нажмите **"Создать раздел"**
3. Заполните данные проекта (название, цель, гипотеза, новизна, ожидаемый результат)
4. Создатель автоматически становится менеджером проекта

### Работа с задачами

1. Нажмите кнопку **"Создать"** в верхней панели
2. Заполните форму:
   - **Раздел** — выберите проект
   - **Тип задачи** — Задача, Баг, Функция или Эпик
   - **Статус** — К выполнению, В РАБОТЕ, ГОТОВО
   - **Резюме** — краткое описание
   - **Описание** — подробное описание
   - **Исполнитель** — назначьте себя или другого пользователя
   - **Приоритет** — от Lowest до Highest
   - **Срок выполнения** — дедлайн
   - **Спринт** — привязка к спринту проекта
3. Нажмите **"Создать"**

### Kanban-доска

1. Перейдите в раздел **"Разделы"**
2. Выберите нужный раздел
3. Нажмите вкладку **"Доска"**
4. Перетаскивайте задачи между колонками для изменения статуса
5. Статус обновляется автоматически с записью в историю

### Спринты

1. Внутри проекта перейдите на вкладку **"Спринты"**
2. Создайте новый спринт (название, цель, даты)
3. Запустите спринт — задачи будут привязаны к активному спринту
4. Завершите спринт по окончании итерации

### Гипотезы и датасеты

- Гипотезы создаются в рамках проекта и привязываются к задачам
- Датасеты загружаются и связываются с задачами через эксперименты

### Отчёт по ГОСТ 7.32

1. Откройте проект
2. Нажмите **"Отчёт по ГОСТ 7.32"**
3. Система сформирует структурированный отчёт на основе данных проекта

### Управление командой

1. Раздел **"Команды"** в левой панели
2. Создайте команду и добавьте участников
3. Команды привязываются к задачам и проектам

## 🔌 API Endpoints

### Аутентификация
- `POST /auth/register` — регистрация
- `POST /auth/login` — вход
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
- `GET /projects/{id}/report` — отчёт по ГОСТ 7.32

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

### Комментарии и вложения
- `POST /tasks/{id}/comments`, `GET /tasks/{id}/comments`
- `POST /tasks/{id}/attachments`, `GET /tasks/{id}/attachments`

## 📝 Примеры API запросов

### Создание задачи
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": 1,
    "title": "Исправить баг в логине",
    "description": "Пользователи не могут войти через Google",
    "status": "К выполнению",
    "priority": "High",
    "assignee_id": 2,
    "created_by": 1,
    "due_date": "2026-05-20T00:00:00Z",
    "team": "Backend Team"
  }'
```

### Получение списка проектов
```bash
curl http://localhost:8080/projects
```

### Обновление статуса задачи
```bash
curl -X PATCH http://localhost:8080/tasks/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "В РАБОТЕ"}'
```

## 🎨 Кастомизация

### Темы
- Светлая тема (по умолчанию)
- Тёмная тема (переключение в меню профиля)

### Цвета
Основной цвет приложения: `#0d6efd` (синий)

Для изменения цветов отредактируйте CSS-переменные в `web/index.html`.

## 🐛 Решение проблем

### Ошибка подключения к БД
```
Cannot connect: could not connect to server
```
**Решение:** Проверьте, что PostgreSQL запущён и соединение правильно настроено в `.env`.

### Порт 8080 уже занят
```bash
# Измените порт в main.go
# или освободите порт:
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### Задача не сохраняется
- Проверьте, что `project_id` существует в БД
- Проверьте, что таблица `tasks` имеет все необходимые колонки
- Выполните все миграции из папки `migrations/`

## 🔐 Безопасность

- Пароли хешируются через `bcrypt` с `bcrypt.DefaultCost`
- Сессионные куки: `HttpOnly: true`, `SameSite: Lax`, срок действия 24 часа
- Параметризованные SQL-запросы защищают от SQL-инъекций
- Файл `.env` исключён из Git

## 📄 Лицензия

Этот проект является частью дипломной работы.

---

**Версия:** 1.0.0  
**Последнее обновление:** май 2026  
**Разработчик:** MVP Team
