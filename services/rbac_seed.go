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

// refreshProjectRoleCache перечитывает проектные роли из БД в кэш
func refreshProjectRoleCache() {
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
