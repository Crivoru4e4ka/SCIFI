package services

import (
	"log"
	"strings"
	"sync"

	"project-MVP/db"
)

// systemRoles — обязательные системные роли
var systemRoles = []struct {
	Name        string
	Description string
}{
	{"admin", "Системный администратор"},
	{"user", "Обычный пользователь"},
	{"guest", "Гость"},
}

// projectRoles — обязательные проектные роли (единый источник)
var projectRoles = []struct {
	Name        string
	Description string
}{
	{"project_lead", "Руководитель проекта"},
	{"scientific_supervisor", "Научный руководитель"},
	{"researcher", "Исследователь"},
	{"analyst", "Аналитик"},
	{"developer", "Разработчик"},
	{"reviewer", "Рецензент"},
	{"viewer", "Наблюдатель"},
}

// permissionDefinitions — все права системы
var permissionDefinitions = []struct {
	Code     string
	Name     string
	Category string
}{
	// Проект
	{"project.view", "Просмотр проекта", "project"},
	{"project.edit", "Редактирование проекта", "project"},
	{"project.delete", "Удаление проекта", "project"},
	{"project.manage_members", "Управление участниками проекта", "project"},
	// Задачи
	{"task.create", "Создание задач", "task"},
	{"task.edit", "Редактирование задач", "task"},
	{"task.delete", "Удаление задач", "task"},
	{"task.assign", "Назначение исполнителей", "task"},
	{"task.change_status", "Изменение статуса задач", "task"},
	// Комментарии
	{"comment.create", "Создание комментариев", "comment"},
	// Вложения
	{"attachment.upload", "Загрузка вложений", "attachment"},
	// Спринты
	{"sprint.manage", "Управление спринтами", "sprint"},
	// Исследования
	{"hypothesis.edit", "Редактирование гипотез", "research"},
	{"experiment.create", "Создание экспериментов", "research"},
	{"experiment.approve", "Утверждение экспериментов", "research"},
	{"analysis.perform", "Проведение анализа", "research"},
	{"results.validate", "Валидация результатов", "research"},
	// Данные
	{"dataset.upload", "Загрузка датасетов", "data"},
	{"dataset.export", "Экспорт датасетов", "data"},
	// Публикации
	{"publication.create", "Создание публикаций", "publication"},
	{"publication.review", "Рецензирование публикаций", "publication"},
	{"publication.approve", "Утверждение публикаций", "publication"},
	// Аудит
	{"audit.view", "Просмотр аудит-лога", "audit"},
}

// rolePermissionMappings — матрица доступа роль → права
var rolePermissionMappings = map[string][]string{
	"admin": {
		"project.view", "project.edit", "project.delete", "project.manage_members",
		"task.create", "task.edit", "task.delete", "task.assign", "task.change_status",
		"comment.create", "attachment.upload", "sprint.manage",
		"hypothesis.edit", "experiment.create", "experiment.approve", "analysis.perform", "results.validate",
		"dataset.upload", "dataset.export",
		"publication.create", "publication.review", "publication.approve",
		"audit.view",
	},
	"user": {
		"project.view", "project.edit",
		"task.create", "task.edit", "task.change_status",
		"comment.create", "attachment.upload",
		"hypothesis.edit", "experiment.create", "analysis.perform",
		"dataset.upload", "publication.create",
	},
	"guest": {
		"project.view", "publication.review",
	},
	"project_lead": {
		"project.view", "project.edit", "project.delete", "project.manage_members",
		"task.create", "task.edit", "task.delete", "task.assign", "task.change_status",
		"comment.create", "attachment.upload", "sprint.manage",
		"hypothesis.edit", "experiment.create", "experiment.approve", "analysis.perform", "results.validate",
		"dataset.upload", "dataset.export",
		"publication.create", "publication.review", "publication.approve",
		"audit.view",
	},
	"scientific_supervisor": {
		"project.view", "project.edit", "project.manage_members",
		"task.create", "task.edit", "task.assign", "task.change_status",
		"comment.create", "attachment.upload", "sprint.manage",
		"hypothesis.edit", "experiment.create", "experiment.approve", "analysis.perform", "results.validate",
		"dataset.upload", "dataset.export",
		"publication.create", "publication.review", "publication.approve",
		"audit.view",
	},
	"researcher": {
		"project.view",
		"task.create", "task.edit", "task.change_status",
		"comment.create", "attachment.upload",
		"hypothesis.edit", "experiment.create", "analysis.perform",
		"dataset.upload", "publication.create",
	},
	"analyst": {
		"project.view",
		"task.create", "task.edit", "task.change_status",
		"comment.create", "attachment.upload",
		"analysis.perform", "results.validate", "dataset.upload", "dataset.export",
	},
	"developer": {
		"project.view",
		"task.create", "task.edit", "task.change_status",
		"comment.create", "attachment.upload",
		"dataset.upload",
	},
	"reviewer": {
		"project.view", "publication.review", "task.edit",
		"comment.create",
	},
	"viewer": {
		"project.view",
	},
}

// projectRoleNamesCache — кэш имён проектных ролей для быстрой валидации
var projectRoleNamesCache = make(map[string]bool)
var projectRoleNamesMutex sync.RWMutex

// InitRoles гарантирует существование всех системных и проектных ролей в БД.
// Вызывается один раз при старте приложения.
func InitRoles() {
	// Системные роли
	for _, r := range systemRoles {
		_, err := db.DB.Exec(`
			INSERT INTO roles (name, description, is_system) VALUES ($1, $2, true)
			ON CONFLICT (name) DO NOTHING
		`, r.Name, r.Description)
		if err != nil {
			log.Printf("InitRoles: ошибка создания системной роли %s: %v", r.Name, err)
		}
	}

	// Проектные роли
	for _, r := range projectRoles {
		_, err := db.DB.Exec(`
			INSERT INTO roles (name, description, is_system) VALUES ($1, $2, false)
			ON CONFLICT (name) DO NOTHING
		`, r.Name, r.Description)
		if err != nil {
			log.Printf("InitRoles: ошибка создания проектной роли %s: %v", r.Name, err)
		}
	}

	refreshProjectRoleCache()
	log.Println("InitRoles: роли инициализированы")
}

// InitPermissions гарантирует существование всех permissions в БД.
func InitPermissions() {
	for _, p := range permissionDefinitions {
		_, err := db.DB.Exec(`
			INSERT INTO permissions (code, name, category) VALUES ($1, $2, $3)
			ON CONFLICT (code) DO NOTHING
		`, p.Code, p.Name, p.Category)
		if err != nil {
			log.Printf("InitPermissions: ошибка создания права %s: %v", p.Code, err)
		}
	}
	log.Println("InitPermissions: права инициализированы")
}

// InitRolePermissions гарантирует корректность связей роль → права.
func InitRolePermissions() {
	for roleName, permCodes := range rolePermissionMappings {
		var roleID int
		err := db.DB.QueryRow(`SELECT id FROM roles WHERE name = $1`, roleName).Scan(&roleID)
		if err != nil {
			log.Printf("InitRolePermissions: роль %s не найдена: %v", roleName, err)
			continue
		}

		for _, code := range permCodes {
			var permID int
			err := db.DB.QueryRow(`SELECT id FROM permissions WHERE code = $1`, code).Scan(&permID)
			if err != nil {
				log.Printf("InitRolePermissions: право %s не найдено: %v", code, err)
				continue
			}

			_, err = db.DB.Exec(`
				INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, roleID, permID)
			if err != nil {
				log.Printf("InitRolePermissions: ошибка связи %s → %s: %v", roleName, code, err)
			}
		}
	}
	log.Println("InitRolePermissions: связи ролей и прав инициализированы")
}

// refreshProjectRoleCache перечитывает проектные роли из БД в кэш
func refreshProjectRoleCache() {
	if db.DB == nil {
		return
	}
	rows, err := db.DB.Query(`SELECT name FROM roles WHERE is_system = false`)
	if err != nil {
		log.Printf("refreshProjectRoleCache: ошибка чтения ролей: %v", err)
		return
	}
	defer rows.Close()

	newCache := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			newCache[name] = true
		}
	}

	projectRoleNamesMutex.Lock()
	projectRoleNamesCache = newCache
	projectRoleNamesMutex.Unlock()
}

// IsValidProjectRole проверяет, является ли имя роли допустимой проектной ролью
func IsValidProjectRole(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}

	projectRoleNamesMutex.RLock()
	valid := projectRoleNamesCache[name]
	projectRoleNamesMutex.RUnlock()

	// Если не нашли в кэше, попробуем обновить кэш (на случай ручного добавления в БД)
	if !valid {
		refreshProjectRoleCache()
		projectRoleNamesMutex.RLock()
		valid = projectRoleNamesCache[name]
		projectRoleNamesMutex.RUnlock()
	}
	return valid
}

// GetProjectRoleNames возвращает список имён всех проектных ролей
func GetProjectRoleNames() []string {
	projectRoleNamesMutex.RLock()
	names := make([]string, 0, len(projectRoleNamesCache))
	for name := range projectRoleNamesCache {
		names = append(names, name)
	}
	projectRoleNamesMutex.RUnlock()
	return names
}
