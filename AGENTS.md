# AGENTS.md — Scifi (Task Management System)

## Общая информация о проекте

Scifi — это веб-приложение для управления задачами и научными проектами, вдохновлённое Jira, с поддержкой Kanban-доски, управления гипотезами, датасетов, спринтов, системы ролей (RBAC), грантов и генерации отчётов по ГОСТ 7.32. Проект разрабатывается в рамках дипломной работы.

- **Язык комментариев и документации:** русский
- **Язык интерфейса:** русский
- **Go module:** `project-MVP`
- **Go version:** 1.25.1

## Технологический стек

| Компонент | Технология | Назначение |
|-----------|-----------|------------|
| Backend | Go 1.25.1 | Серверная логика, REST API |
| HTTP-роутер | Gorilla Mux v1.8.1 | Маршрутизация HTTP-запросов |
| База данных | PostgreSQL 12+ | Реляционное хранилище данных |
| Драйвер БД | `lib/pq` v1.12.0 | Драйвер PostgreSQL для Go |
| Хеширование паролей | `golang.org/x/crypto/bcrypt` | Безопасное хранение паролей |
| Переменные окружения | `joho/godotenv` v1.5.1 | Загрузка `.env` |
| Генерация Excel | `github.com/xuri/excelize/v2` | Экспорт проектов в Excel |
| Генерация PDF | `github.com/jung-kurt/gofpdf` | Экспорт проектов в PDF |
| Swagger | `github.com/swaggo/http-swagger` | Документация API |
| Frontend | Vue.js 3 (global build из CDN) | Реактивный UI |
| CSS-фреймворк | Bootstrap 5.3.0 | Стилизация компонентов |
| Иконки | Font Awesome 6.0.0 | Иконографика |
| Форматирование | Prettier 3.8.3 | Форматирование кода (devDependency) |

---

## Структура проекта

```
project-MVP/
├── main.go                          # Точка входа
├── go.mod / go.sum                  # Go-модули
├── package.json                     # Node devDependencies (Prettier)
├── .env                             # Переменные окружения
├── .gitignore                       # Исключения Git
│
├── db/
│   ├── db.go                        # Подключение к PostgreSQL, глобальная переменная DB *sql.DB
│   ├── interfaces.go                # Интерфейсы DBTX и DBPool для тестируемости
│   └── schema.sql                   # Полная схема БД (~560 строк, 21+ таблиц)
│
├── models/                          # DTO / структуры данных (18 файлов)
│   ├── user.go
│   ├── project.go                   # + ResearchGoal, MainHypothesis, Novelty, ExpectedResult
│   ├── task.go                      # + Type, HypothesisID, ResearchMethod, ResearchContribution, DOI
│   ├── comment.go                   # + ParentID, SoftDelete, Nested Replies
│   ├── attachment.go
│   ├── sprint.go
│   ├── tag.go
│   ├── tasktag.go
│   ├── task_history.go
│   ├── team.go                      # + TeamMemberInfo, UpdateTeamRequest, AddMemberRequest
│   ├── project_member.go            # + RoleId (FK к roles)
│   ├── hypothesis.go
│   ├── dataset.go                   # + Parameters (JSONB)
│   ├── experimentdataset.go
│   ├── grant.go                     # + GrantType, Status constants
│   ├── project_grant_funding.go
│   └── rbac.go                      # Role, Permission, RolePermission, ProjectMemberWithRole
│
├── handlers/                        # HTTP-обработчики (20+ файлов)
│   ├── auth_handler.go              # POST /auth/register, /auth/login, GET /logout
│   ├── user_handler.go              # POST/GET /users, GET /me
│   ├── project_handler.go           # CRUD проектов, прогресс, assignable users
│   ├── task_handler.go              # CRUD задач, статус, спринт, гипотезы
│   ├── comment_handler.go           # CRUD комментариев с soft-delete и проверкой прав
│   ├── attachment_handler.go        # Загрузка файлов, список вложений, отчёт ГОСТ 7.32
│   ├── project_member_handler.go    # Добавление/удаление участников, смена ролей
│   ├── sprint_handler.go            # CRUD спринтов, start/complete
│   ├── team_handler.go              # CRUD команд, управление участниками
│   ├── tag_handler.go
│   ├── task_tag_handler.go
│   ├── dataset_handler.go
│   ├── experiment_dataset_handler.go
│   ├── grant_handler.go             # Полный CRUD грантов + финансирование проектов
│   ├── rbac_handler.go              # Роли, права, аудит-лог
│   ├── export_handler.go            # Excel/PDF экспорт
│   ├── activity_handler.go          # Лента активности
│   ├── ui_handler.go                # Раздача HTML и редиректы
│   ├── helpers.go                   # GetUserID, RequireAuth, RequirePermission, RequireAdmin, AtoiParam
│   ├── auth_handler_test.go
│   └── helpers_test.go
│
├── services/                        # Бизнес-логика (30+ файлов)
│   ├── user_service.go / user_store.go
│   ├── project_service.go / project_store.go
│   ├── task_service.go / task_store.go
│   ├── comment_service.go / comment_store.go
│   ├── attachment_service.go
│   ├── project_member_service.go
│   ├── sprint_service.go / sprint_store.go
│   ├── team_service.go / team_store.go
│   ├── tag_service.go / tag_store.go
│   ├── task_tag_service.go
│   ├── dataset_service.go / dataset_store.go
│   ├── experiment_service.go / experiment_store.go
│   ├── grant_service.go
│   ├── project_grant_funding_service.go
│   ├── export_service.go            # Генерация Excel и PDF
│   ├── activity_service.go
│   ├── rbac_service.go / rbac_store.go / rbac_seed.go
│   ├── init.go                      # InitDefaultStores — инициализация store-экземпляров
│   └── *_test.go (17 тестовых файлов)
│
├── middleware/
│   ├── auth.go                      # AuthMiddleware, GetUserID
│   └── auth_test.go                 # 6 unit-тестов middleware
│
├── routes/
│   ├── routes.go                    # Все маршруты + NoCacheMiddleware
│   └── routes_test.go               # Тесты роутов
│
├── web/
│   ├── index.html                   # Основное SPA на Vue.js (~5900 строк)
│   └── login.html                   # Страница входа/регистрации
│
├── static/
│   ├── uploads/                     # Загруженные пользователями файлы
│   └── *.png                        # Логотипы, бейджи
│
├── docs/
│   ├── docs.go, swagger.json, swagger.yaml  # Swagger-документация
│
└── migrations/                      # 9 SQL-миграций
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

---

## Системная архитектура и проектирование

### Трёхслойная архитектура

Приложение построено по классической трёхслойной архитектуре с дополнительным слоем хранилища (Store pattern):

```
┌─────────────────────────────────────────────────────────────┐
│  Клиент (Браузер)                                           │
│  Vue.js 3 SPA → fetch() API                                 │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP / JSON
┌──────────────────────▼──────────────────────────────────────┐
│  Слой маршрутизации и middleware                            │
│  Gorilla Mux + AuthMiddleware + NoCacheMiddleware           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Handlers (HTTP-обработчики)                                │
│  Валидация входных данных, извлечение параметров            │
│  Вызов сервисов, формирование HTTP-ответов                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Services (Бизнес-логика)                                   │
│  Валидация бизнес-правил, транзакции, координация store     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Stores (Слой доступа к данным)                             │
│  SQL-запросы к PostgreSQL через db.DB / tx                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  PostgreSQL                                                 │
└─────────────────────────────────────────────────────────────┘
```

### Поток данных через слои (пример: обновление статуса задачи на Kanban-доске)

1. **Клиент**: пользователь перетаскивает задачу в другую колонку. Vue.js вызывает `fetch('/tasks/123/status', {method: 'PATCH', body: JSON.stringify({status: 'В РАБОТЕ'})})`.
2. **Middleware**: `AuthMiddleware` извлекает `session` cookie, парсит `user_id=42`, помещает в `request.Context()`. `NoCacheMiddleware` добавляет заголовки `Cache-Control`.
3. **Router**: Gorilla Mux направляет запрос в `handlers.UpdateTaskStatus`.
4. **Handler**: извлекает `user_id` из контекста через `middleware.GetUserID(r)`, парсит `id` из URL через `AtoiParam`, декодирует JSON. Проверяет права через `RequirePermission(w, r, projectID, "task.change_status")`.
5. **Service**: `services.UpdateTaskStatus` проверяет, что статус допустим, запускает транзакцию: обновляет `tasks.status`, вставляет запись в `task_history` (field_name, old_value, new_value), вызывает `LogActivity` для ленты событий.
6. **Store**: `task_store.go` выполняет параметризованные SQL-запросы (`$1`, `$2`) через `db.DB`.
7. **Ответ**: Handler возвращает JSON с обновлённой задачей и HTTP 200.
8. **Клиент**: Vue.js получает ответ, обновляет локальное состояние задачи в массиве `tasks`, Kanban-доска перерисовывается.

### Store Pattern для тестируемости

Большинство сервисов разделены на два файла: `*_service.go` (бизнес-логика) и `*_store.go` (SQL-операции). Store-структуры принимают интерфейс `db.DBTX` (обобщает `*sql.DB` и `*sql.Tx`), что позволяет:
- Использовать одни и те же методы внутри транзакций и вне их.
- Мокать БД в unit-тестах через `sqlmock` (передаётся `sqlmock.NewMockDB()` как `DBTX`).

Пример: `DefaultTaskStore` имеет метод `UpdateTaskStatus(ctx context.Context, dbtx db.DBTX, taskID int, status string) error`. В production передаётся `db.DB`, в тестах — `mockDB`.

---

## Проектирование базы данных

### Схема БД (21 таблица + связующие)

Полная схема описана в `db/schema.sql` (~560 строк). Ниже — концептуальная модель.

#### Основные сущности и связи

```
users (1) ───────< (N) projects (created_by)
users (1) ───────< (N) tasks (created_by, assignee_id)
users (1) ───────< (N) comments
users (1) ───────< (N) attachments
users (1) ───────< (N) activities
users (1) ───────< (N) grants (created_by, principal_investigator_id)
users (1) ───────< (N) hypotheses (created_by)

projects (1) ────< (N) tasks
projects (1) ────< (N) sprints
projects (1) ────< (N) hypotheses
projects (1) ────< (N) datasets
projects (1) ────< (N) project_members
projects (1) ────< (N) attachments (через tasks)
projects (1) ────< (N) project_grant_funding
projects (1) ────< (1) teams (team_id, optional)

tasks (1) ───────< (N) task_history
 tasks (1) ───────< (N) comments
 tasks (1) ───────< (N) attachments
 tasks (N) ───────> (N) tags (через task_tags)
 tasks (N) ───────> (N) datasets (через experiment_datasets)
 tasks (1) ───────> (1) sprints (sprint_id, optional)
 tasks (1) ───────> (1) hypotheses (hypothesis_id, optional)

teams (1) ───────< (N) team_members
 teams (1) ───────< (N) projects

roles (1) ───────< (N) project_members (role_id)
roles (N) ───────> (N) permissions (через role_permissions)
```

#### Ключевые поля и типы

- **JSONB**: `datasets.parameters`, `tasks.parameters`, `tasks.metrics` — гибкие структурированные данные.
- **Soft delete**: `comments` имеет `deleted_at` (в модели), обновляется через `UPDATE`, не физическое удаление.
- **Генерируемые ID**: большинство таблиц используют `GENERATED ALWAYS AS IDENTITY` (современный PostgreSQL-стандарт), некоторые `serial`.
- **Внешние ключи**: строгая ссылочная целостность с `ON DELETE CASCADE` для зависимых сущностей (task_tags, attachments, comments) и `SET NULL` для создателей/исполнителей (чтобы не терять историю при удалении пользователя).
- **Уникальные ограничения**: `users(email)`, `roles(name)`, `permissions(code)`, `tags(name)`, `project_members(project_id, user_id)`.

#### Миграционная стратегия

Миграции применяются вручную через `psql`. История миграций:
1. `add_task_fields.sql` — добавление `start_date`, `team` в `tasks`.
2. `fix_task_history.sql` — расширение аудита: `changed_by`, `field_name`, `old_value`, `new_value`.
3. `2026-05-10_add_new_tables_and_fields.sql` — научные сущности: `datasets`, `experiment_datasets`, `hypotheses`, `tags`, `task_tags`; поля `type`, `hypothesis_id`, `conclusion`, `parameters`, `metrics` в `tasks`.
4. `2026-05-16_add_grants_and_funding.sql` — `grants`, `project_grant_funding`.
5. `2026-05-16_add_rbac_system.sql` — `roles`, `permissions`, `role_permissions`; `role_id` в `project_members`; начальное заполнение ролей и прав.
6. `2026-05-16_fix_users_role_check.sql` — ограничение `users.role` на `admin`/`user`/`guest`.
7. `2026-05-16_link_projects_to_teams.sql` — `team_id`, `execution_type` в `projects`; миграция старых ролей.
8. `2026-05-16_unify_team_roles.sql` — унификация ролей: `admin`→`project_lead`, `member`→`researcher` и т.д.
9. `add_task_doi_column.sql` — `doi` в `tasks`.

---

## Проектирование API

### Принципы REST

- Ресурсы именуются существительными: `/projects`, `/tasks`, `/comments`.
- HTTP-методы определяют действие: `GET` — чтение, `POST` — создание, `PATCH` — частичное обновление, `DELETE` — удаление.
- Идентификаторы в URL: `/tasks/{id}/status`.
- Формат данных: JSON. Файлы: `multipart/form-data`.

### Глобальное middleware

Все запросы проходят через два middleware (регистрируются в `routes.go` через `r.Use(...)`):

1. **`middleware.AuthMiddleware`**:
   - Читает cookie `session`.
   - Если cookie валидна (целое число), помещает `user_id` в `request.Context()` под ключом `userIDKey`.
   - **НЕ прерывает цепочку** при отсутствии/невалидности cookie. Это позволяет публичным endpoint'ам (`/auth/register`, `/auth/login`, `/static/`, `/swagger/`, `/login`, `/`) работать без авторизации.
   - Защищённые endpoint'ы сами решают, требовать ли авторизацию (через `RequireAuth` или `RequirePermission`).

2. **`routes.NoCacheMiddleware`**:
   - Устанавливает заголовки `Cache-Control: no-cache, no-store, must-revalidate`, `Pragma: no-cache`, `Expires: 0`.
   - Предотвращает кэширование приватных страниц в браузере.

### Группы endpoint'ов

| Группа | Endpoint'ы | Особенности |
|--------|-----------|-------------|
| Auth | `/auth/register`, `/auth/login`, `/logout` | Куки session, bcrypt |
| Users | `/users`, `/me` | CRUD пользователей |
| Projects | `/projects`, `/user/projects`, `/projects/{id}/...` | Авто-назначение создателя как lead |
| Tasks | `/tasks`, `/tasks/{id}/status`, `/tasks/{id}/sprint` | История статусов, спринты |
| Sprints | `/projects/{id}/sprints`, `/sprints/{id}/start` | Жизненный цикл спринта |
| Comments | `/comments`, `/comments/{type}/{id}` | Soft-delete, nested replies |
| Attachments | `/tasks/{id}/attachments` | Multipart upload, `./static/uploads/` |
| Teams | `/teams`, `/teams/{id}/members` | Управление командами и ролями |
| Project Members | `/project-members`, `/projects/{id}/members` | Роли проектных участников |
| RBAC | `/roles`, `/permissions`, `/projects/{id}/activities` | Аудит, права |
| Grants | `/grants`, `/projects/{id}/grants` | Финансирование, бюджет |
| Reports | `/projects/{id}/report`, `/export/excel`, `/export/pdf` | ГОСТ 7.32, Excel, PDF |
| Activities | `/activities` | Лента событий |
| Static | `/static/*` | `http.FileServer` |
| Swagger | `/swagger/*` | Документация API |

---

## Проектирование фронтенда

### Архитектура SPA

Фронтенд — это **Single Page Application** на Vue.js 3 (global build из CDN) без сборки (no build step). Весь UI находится в двух HTML-файлах:
- `web/login.html` — автономная страница аутентификации.
- `web/index.html` — основное приложение (~5900 строк).

### Управление состоянием и маршрутизация

**Нет Vue Router и Vuex/Pinia.** Маршрутизация реализована императивно через реактивные свойства корневого компонента:
- `currentPage` — текущая страница.
- `activeSidebarTab` — выбранный пункт бокового меню (`recent`, `flagged`, `sections`, `filters`, `dashboards`, `teams`, `grants`).
- `currentProject` — если установлен, UI переключается в "режим проекта".
- `currentProjectView` — под-виды проекта: `list`, `board`, `backlog`, `attachments`, `funding`, `audit`, `discussion`.
- `currentGrant` — если установлен, показывается детальная карточка гранта.

Персистентность части состояния через `localStorage`:
- `scifi_saved_filters_<userId>` — сохранённые фильтры.
- `scifi_recent_<userId>` — недавние элементы.
- `scifi_flagged_<userId>` — отмеченные элементы.
- `appTheme` — тёмная/светлая тема.

### Коммуникация с backend

- **Transport**: нативный `fetch()` API (без Axios).
- **Auth**: session cookie (`HttpOnly: true`, `SameSite: Lax`). `credentials: 'include'` на запросах.
- **Content-Type**: `application/json` для API, `multipart/form-data` для файлов.

### Основные UI-секции

**Глобальная раскладка:**
- **Top Navbar**: логотип, глобальный поиск с выпадающим списком (проекты + задачи + быстрые фильтры), кнопка "Создать", иконки уведомлений/помощи/настроек, выпадающее меню профиля (инициалы, имя, email, переключатель темы, выход).
- **Left Sidebar**: разделы "Для Вас" (Недавние, Отмеченные, Проекты) и "Рекомендуется" (Фильтры, Дашборды, Команды, Гранты).
- **Main Content**: динамически меняется на основе состояния.

**Виды проекта:**
- **List** (`list`): разворачиваемая Bootstrap-таблица задач. Клик по строке открывает панель деталей с описанием, тегами, научным контекстом (гипотеза, вклад, метод), параметрами (эксперимент, исследование, анализ), комментариями.
- **Kanban Board** (`board`): 4 колонки (todo, inProgress, review, done). HTML5 drag-and-drop (`draggable`, `@dragstart`, `@drop`). При drop вызывается `PATCH /tasks/{id}/status`.
- **Backlog** (`backlog`): бэклог проекта.
- **Attachments** (`attachments`): список файлов проекта.
- **Funding** (`funding`): связанные гранты и их бюджеты.
- **Audit** (`audit`): аудит-лог проекта.
- **Discussion** (`discussion`): комментарии к проекту.

**Модальные окна:**
- Создание проекта (пошаговый мастер: шаблон → название/key → научный паспорт → участники).
- Создание/редактирование задачи (динамическая форма в зависимости от типа: research, experiment, data_collection, analysis, dev, doc/publication).
- Управление спринтом, командой, грантом, участниками проекта.

### Кастомные директивы

- `v-click-outside` — закрывает выпадающий поиск при клике вне его области.

---

## Разработка backend

### Последовательность запуска (main.go)

```go
func main() {
    godotenv.Load()          // 1. Загрузка .env
    db.Connect()             // 2. Подключение к PostgreSQL
    services.InitDefaultStores()  // 3. Инициализация глобальных store-экземпляров
    services.InitRoles()     // 4. Сидинг системных и проектных ролей
    services.InitPermissions()    // 5. Сидинг прав
    services.InitRolePermissions() // 6. Сидинг связей роль→права
    r := routes.InitRoutes() // 7. Маршрутизация
    http.ListenAndServe(":8080", r) // 8. Запуск сервера
}
```

### Слой Handlers

Handlers не содержат бизнес-логики. Их задачи:
1. Декодировать JSON/multipart из запроса.
2. Извлечь параметры (URL vars, query strings).
3. Проверить авторизацию/права через хелперы из `helpers.go`.
4. Вызвать соответствующий сервис.
5. Сформировать HTTP-ответ (JSON, файл, редирект, ошибка с кодом).

**Хелперы авторизации (`helpers.go`):**
- `GetUserID(r)` — получает `user_id` из контекста (установлен `AuthMiddleware`).
- `RequireAuth(w, r)` — возвращает 401, если пользователь не авторизован.
- `RequirePermission(w, r, projectID, permissionCode)` — проверяет право через RBAC-store; возвращает 403 при отказе.
- `RequireAdmin(w, r)` — проверяет системную роль `admin`.
- `AtoiParam(s)` — безопасный парсинг URL-параметров в `int` (отклоняет 0 и отрицательные).

**Пример взаимодействия в handler**: `CreateTaskInProject`
- Парсит `project_id` из URL.
- Декодирует JSON в `models.Task`.
- Устанавливает `ProjectId` из URL (переопределяет тело запроса).
- Проверяет `RequirePermission(..., "task.create")`.
- Вызывает `services.CreateTask(task)`.
- При успехе возвращает 201 + JSON с созданной задачей.
- При ошибках валидации возвращает 400 с текстом ошибки.

### Слой Services

Services содержат всю бизнес-логику, правила валидации и координацию транзакций.

**Паттерн Service + Store:**
- `*_service.go` — высокоуровневая логика (валидация, координация нескольких store, аудит).
- `*_store.go` — низкоуровневые SQL-операции (CRUD, JOIN, агрегации).

**Пример**: `task_service.go` / `task_store.go`
- `task_service.go`:
  - `CreateTask` — проверяет существование проекта и создателя через `project_service` и `user_service`, валидирует обязательные поля (`project_id`, `title`, `created_by`), проверяет что `assignee_id` является участником проекта (`project_member_service.IsUserInProject`), нормализует статус, создаёт задачу через `DefaultTaskStore`, при наличии тегов вызывает `task_tag_service.AddTagToTask`, логирует активность.
  - `UpdateTaskStatus` — проверяет валидность статуса, вызывает store для обновления `tasks.status`, вставляет запись в `task_history`, вызывает `activity_service.LogActivity`.
- `task_store.go`:
  - `CreateTask` — `INSERT INTO tasks ... RETURNING id, task_num`.
  - `GetTaskByID` — `SELECT` с `JOIN` на `users` (assignee), `STRING_AGG` для тегов.
  - `UpdateTask` — `UPDATE tasks SET ...`.
  - `DeleteTask` — в транзакции удаляет `task_tags`, `task_history`, `comments`, затем `tasks`.
  - `GetProjectReportData` — сложный `SELECT` с агрегацией данных для отчёта ГОСТ 7.32.

**Транзакции:**
- Используются для операций, требующих атомарности: создание проекта с автодобавлением создателя как lead, удаление проекта с каскадной очисткой, завершение спринта с возвратом незавершённых задач в бэклог.
- Транзакция начинается в service или store через `db.DB.Begin()`, передаётся как `*sql.Tx` (реализует `db.DBTX`) в store-методы.

**Публичные ошибки сервисов:**
```go
var (
    ErrNotFound           = errors.New("not found")
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrInvalidRole        = errors.New("invalid role")
    ErrDuplicateUser      = errors.New("user with this email already exists")
)
```

Handlers проверяют конкретные ошибки и возвращают соответствующие HTTP-коды (400, 401, 403, 404, 409).

### Аутентификация и авторизация

**Аутентификация:**
- Сессионная модель через HTTP-only cookie `session`.
- При логине (`POST /auth/login`) сервис `AuthenticateUser` сверяет bcrypt-хеш пароля. При успехе handler устанавливает cookie со строковым значением `user_id`, `HttpOnly: true`, `SameSite: Lax`, `MaxAge: 86400` (24 часа).
- При логауте (`GET /logout`) cookie сбрасывается (`MaxAge: -1`) и выполняется редирект на `/login`.
- Системные роли пользователей: `admin`, `user`, `guest`.

**RBAC (Role-Based Access Control):**
- Реализован через три слоя: `rbac_seed.go` (начальное заполнение), `rbac_service.go` (бизнес-логика), `rbac_store.go` (SQL).
- **Системные роли** (`is_system = true`): `admin`, `user`, `guest` — определяют глобальные права.
- **Проектные роли** (`is_system = false`): `project_lead`, `scientific_supervisor`, `researcher`, `analyst`, `developer`, `reviewer`, `viewer`.
- **Права** (permissions): 17 кодов прав (`project.view`, `task.create`, `experiment.approve`, `audit.view` и др.).
- **Матрица доступа**: `rolePermissionMappings` в `rbac_seed.go` определяет какие права есть у каждой роли.
- **Проверка прав** (`RequirePermission`):
  1. Получает `user_id` из контекста.
  2. Запрашивает `rbac_store` — функция `CheckPermission` выполняет SQL-запрос: находит роль пользователя в проекте (через `project_members` → `roles`), затем проверяет наличие связи `role_permissions` с нужным `permission.code`.
  3. Если право есть — продолжает; если нет — 403 Forbidden.
- **Аудит**: лента активности проекта (`GET /projects/{id}/activities`) использует таблицу `activities`. Записи создаются при CRUD-операциях через `services.LogActivity`.

### Файловые вложения

- Загрузка через `multipart/form-data` (`POST /tasks/{id}/attachments`).
- Файл сохраняется на диск в `./static/uploads/` под уникальным именем.
- Метаданные (filename, URL, task_id, user_id) записываются в таблицу `attachments`.
- Статические файлы раздаются через `http.FileServer` по префиксу `/static/`.

---

## Разработка frontend

### Vue.js 3 приложение в index.html

Приложение монтируется через `createApp({...}).mount('#app')`.

**Данные (data):**
- `currentUser` — объект пользователя, загружаемый при старте через `GET /me`.
- `projects`, `tasks`, `teams`, `grants` — массивы сущностей.
- `currentProject`, `currentTask`, `currentGrant` — выбранные сущности.
- `isDarkTheme` — булево для тёмной темы.
- `searchQuery`, `searchResults` — глобальный поиск.
- Множество булевых флагов для открытия/закрытия модальных окон.

**Вычисляемые свойства (computed):**
- `filteredTasks` — фильтрация задач по статусу, приоритету, исполнителю, спринту.
- `kanbanColumns` — группировка задач по статусам для доски.
- `projectProgress` — процент выполненных задач.

**Методы:**
- `loadProjects()`, `loadTasks()`, `loadTeams()`, `loadGrants()` — загрузка данных с API.
- `createTask()`, `updateTaskStatus()`, `deleteTask()` — CRUD задач.
- `onDragStart(task)`, `onDrop(column)` — обработка DnD Kanban.
- `saveFilter()`, `applyFilter()` — работа с фильтрами.
- `toggleTheme()` — переключение темы (добавляет/убирает класс `dark-theme` на `<body>`).

### Динамические формы задач

Форма создания/редактирования задачи изменяется в зависимости от выбранного `type`:
- **research** — поля для гипотезы, метода, вклада.
- **experiment** — поля для параметров, метрик, ожидаемых результатов.
- **data_collection** — поля для датасетов.
- **analysis** — поля для методов анализа.
- **dev** — стандартные поля разработки.
- **doc/publication** — поля для DOI, публикаций.

Vue использует `v-if`/`v-show` для условного отображения блоков формы.

### Kanban-доска

- Данные: массив `tasks` фильтруется по `currentProject.id`.
- Колонки: `К выполнению`, `В работе`, `На проверке`, `Готово`.
- Каждая задача (`<div>`) имеет `draggable="true"` и обработчики `@dragstart` (сохраняет `task.id` в `dataTransfer`) и `@dragend`.
- Колонки имеют `@dragover.prevent` и `@drop` (читает `task.id`, обновляет `status` соответствующей колонке, вызывает API `PATCH /tasks/{id}/status`).
- После успешного API-вызова задача перемещается в другой массив колонки в реактивном состоянии Vue.

---

## Реализация ключевых функциональных модулей

### 1. Управление проектами

**Создание проекта** (`POST /projects`):
1. Handler валидирует обязательные поля (`name`, `key`).
2. `project_service.ValidateAndNormalizeProject` нормализует ключ, проверяет тип исполнения (`execution_type`: `team` или `manual`).
3. `project_store.CreateProject` выполняет **транзакцию**:
   - `INSERT INTO projects ... RETURNING id`.
   - Получает `role_id` для роли `project_lead` из таблицы `roles`.
   - `INSERT INTO project_members (project_id, user_id, role, role_id)` — создатель становится руководителем.
   - Если указан `team_id` и `execution_type == "team"`, синхронизирует участников команды в проект через `SyncTeamMembersToProject`.
4. Возвращается проект с `id`.

**Прогресс проекта** (`GET /projects/{id}/progress`):
- SQL-запрос: `COUNT(*) FILTER (WHERE status = 'ГОТОВО') / COUNT(*) * 100`.
- Возвращает процент выполнения.

**Удаление проекта** (`DELETE /projects/{id}`):
- Требует право `project.manage_members`.
- Транзакция: удаляет `project_members`, `tasks` (каскадно), затем сам проект.

### 2. Задачи и жизненный цикл

**Создание задачи**:
- Обязательные поля: `project_id`, `title`, `created_by`.
- `task_num` генерируется автоматически (последовательность в рамках проекта).
- Ключ задачи формируется как `{ProjectKey}-{TaskNum}` (например, `SCI-12`).
- Если указан `assignee_id`, проверяется, что пользователь состоит в проекте.
- Если в `tags` передана строка с тегами через запятую, они создаются/привязываются через `task_tag_service`.

**Обновление статуса** (`PATCH /tasks/{id}/status`):
- Handler проверяет право `task.change_status`.
- `task_service.UpdateTaskStatus`:
  1. Валидирует новый статус (допустимые: `К выполнению`, `В работе`, `На проверке`, `Готово`).
  2. В транзакции обновляет `tasks.status`.
  3. Вставляет запись в `task_history` (`task_id`, `changed_by`, `field_name`, `old_value`, `new_value`, `changed_at`).
  4. Вызывает `activity_service.LogActivity` для ленты событий.

### 3. Спринты

**Создание спринта** (`POST /projects/{id}/sprints`):
- Требует право `sprint.manage`.
- Статус по умолчанию: `planned`.

**Запуск спринта** (`PATCH /sprints/{id}/start`):
- Проверяет, что спринт в статусе `planned`.
- Обновляет статус на `active`.

**Завершение спринта** (`PATCH /sprints/{id}/complete`):
- Транзакция:
  1. Обновляет `sprints.status` на `completed`.
  2. Находит все задачи спринта со статусом ≠ `Готово`.
  3. Сбрасывает у них `sprint_id` в `NULL` (возврат в бэклог).

### 4. Команды и участники

**Создание команды** (`POST /teams`):
- Создатель автоматически становится `project_lead` в `team_members`.
- Можно сразу указать `MemberEmails` — участники приглашаются по email.
- Команда может быть привязана к проекту через `projects.team_id`.

**Синхронизация команды с проектом** (`SyncTeamMembersToProject`):
- При создании проекта с `execution_type == "team"` все участники команды копируются в `project_members`.
- Роли команды маппятся на проектные роли (`team lead` → `project_lead`, `researcher` → `researcher`).

### 5. Гранты и финансирование

**Грант** (`grants`):
- Поля: название, код, организация-финансировщик, страна, описание, научное направление, тип (`state`, `university`, `international`, `corporate`, `internal`), статус (`draft` → `submitted` → ... → `completed`), бюджет, валюта, даты.

**Финансирование проекта** (`project_grant_funding`):
- Связь many-to-many между грантом и проектом с указанием выделенной суммы (`allocated_amount`).
- Сервис валидирует, что сумма выделения не превышает общий бюджет гранта.
- `GET /grants/{id}/budget` возвращает `total`, `allocated`, `remaining`.

### 6. Отчёты и экспорт

**ГОСТ 7.32** (`GET /projects/{id}/report`):
- Handler `GenerateGostReport` вызывает `services.GetProjectReportData`.
- SQL-запрос агрегирует данные проекта: название, цель, гипотеза, новизна, ожидаемый результат, список задач с их статусами и описаниями.
- Формируется текстовый отчёт в формате ГОСТ 7.32.
- Возвращается как `text/plain` с `Content-Disposition: attachment; filename=report.txt`.

**Excel / PDF** (`GET /projects/{id}/export/...`):
- `export_service.go` использует библиотеки `excelize` и `gofpdf`.
- Получает проект и его задачи, генерирует файл в памяти (`bytes.Buffer`), возвращает как бинарный поток с правильным `Content-Type`.

### 7. Комментарии

- Поддерживаются комментарии к задачам (`entity_type = "task"`) и проектам (`entity_type = "project"`).
- **Вложенность**: поле `parent_id` позволяет создавать ответы на комментарии.
- **Soft delete**: `DELETE /comments/{id}` не удаляет строку, а устанавливает `deleted_at = NOW()`.
- **Права**: удалять/редактировать может только автор или руководитель проекта (`isAuthor || isLead`).
- **Дерево комментариев**: `comment_service.BuildCommentTree` преобразует плоский список из БД во вложенную структуру (`Replies []Comment`).

### 8. Аудит и активность

- **`task_history`**: фиксирует изменения статусов задач (кто, когда, с какого статуса на какой).
- **`activities`**: лента событий для проектов и пользователей. Заполняется через `services.LogActivity`. `GET /activities` возвращает последние 50 событий по проектам, в которых участвует пользователь. `GET /projects/{id}/activities` возвращает события конкретного проекта.
- **`activities`**: лента событий для пользователя. Заполняется через `services.LogActivity`. `GET /activities` возвращает последние 50 событий по проектам, в которых участвует пользователь.

---

## Тестирование

### Общая информация

В проекте **21 тестовый файл** с unit-тестами. Основной акцент на слое `services/` (17 тестовых файлов). Тесты используют `testify/assert` и `sqlmock` для мокирования БД.

### Структура тестов

| Слой | Тестовые файлы | Количество | Покрытие |
|------|---------------|------------|----------|
| `middleware/` | `auth_test.go` | 6 тестов | AuthMiddleware, GetUserID |
| `handlers/` | `auth_handler_test.go`, `helpers_test.go` | 9 тестов | Logout, AtoiParam, RequireAuth |
| `routes/` | `routes_test.go` | — | Роуты |
| `services/` | 17 файлов `*_test.go` | ~100+ тестов | Store, Service, RBAC seeding |

### Тестируемый слой Stores

Store-файлы протестированы через `sqlmock`:
- Подключается `sqlmock.New()`.
- Создаётся экземпляр Store с `mockDB` (реализует `db.DBTX`).
- Устанавливаются ожидания (`ExpectQuery`, `ExpectExec`) с регулярными выражениями для SQL.
- Вызываются методы store.
- Проверяется, что все ожидания выполнены (`mock.ExpectationsWereMet()`).

**Пример тестируемых сценариев:**
- `user_store_test.go` — создание пользователя, поиск по email, обработка `ErrNoRows`.
- `task_store_test.go` — создание задачи, получение по ID, обновление статуса, удаление с каскадом.
- `rbac_store_test.go` — проверка прав, назначение ролей, аудит-лог.
- `project_store_test.go` — создание проекта с транзакцией, синхронизация команды.
- `team_store_test.go` — создание команды, добавление участников, изменение ролей.
- `sprint_store_test.go` — запуск и завершение спринта.
- `grant_service_test.go` — валидация бюджета, дат, типов грантов.
- `export_service_test.go` — генерация Excel/PDF.

### Интерфейс DBTX для тестирования

```go
// db/interfaces.go
type DBTX interface {
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
```

Этот интерфейс позволяет передавать в store-методы как `*sql.DB`, так и `*sql.Tx`, а в тестах — `sqlmock.Sqlmock`.

### Что покрыто тестами

- **CRUD операции** всех основных сущностей (users, projects, tasks, teams, grants, datasets, comments).
- **Транзакции** (создание проекта с авто-добавлением lead, завершение спринта).
- **RBAC** (проверка прав, сидинг ролей и permission mappings).
- **Валидация бизнес-правил** (дубликаты email, невалидные роли, бюджетные ограничения грантов, форматы дат).
- **Генерация отчётов** (Excel/PDF не падают при валидных данных).

### Что НЕ покрыто тестами (ручное тестирование)

- End-to-end тесты frontend-backend.
- Интеграционные тесты с реальной PostgreSQL.
- Тесты файловой загрузки (multipart).
- Тесты drag-and-drop Kanban-доски.
- Нагрузочное тестирование.

---

## Безопасность

| Угроза | Мера защиты |
|--------|-------------|
| Перехват паролей | `bcrypt` с `bcrypt.DefaultCost` |
| XSS / Session hijacking | HTTP-only cookie `session`, `SameSite: Lax` |
| Кэширование приватных данных | `NoCacheMiddleware` на всех запросах |
| SQL-инъекции | Параметризованные запросы (`$1`, `$2`) |
| Несанкционированный доступ | RBAC с проверкой на каждом защищённом endpoint'е |
| Утечка .env | Файл `.env` в `.gitignore` |
| IDOR (доступ к чужим данным) | Проверка `project_members` перед доступом к проектным ресурсам |
| File upload risks | Файлы сохраняются вне web-root с раздачей через handler |

---

## Сборка и запуск

### Требования
- Go 1.25.1+
- PostgreSQL 12+
- Node.js (опционально, только для Prettier)

### Настройка
```bash
# Создать .env
cat > .env <<EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=пароль
DB_NAME=SRA
EOF

# Инициализация БД
psql -h localhost -U postgres -d SRA -f db/schema.sql
for f in migrations/*.sql; do psql -h localhost -U postgres -d SRA -f "$f"; done

# Зависимости
go mod download
```

### Запуск
```bash
go run main.go
# Сервер на :8080
```

### Сборка бинарника
```bash
go build -o app.exe main.go
```

### Деплой
Приложение — одиночный Go-бинарник, который:
1. Слушает порт `:8080`.
2. Раздаёт статические файлы из `./static/`.
3. Раздаёт HTML из `./web/`.
4. Обрабатывает REST API.

Достаточно скопировать бинарник, папки `static/`, `web/`, `migrations/` и настроить переменные окружения.

---

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

---

## Полезные файлы для ознакомления

- `README.md` — пользовательская документация с примерами API-запросов
- `db/schema.sql` — полная схема базы данных
- `migrations/*.sql` — история изменений схемы
- `docs/swagger.yaml` — Swagger-документация API
