package services

import (
	"database/sql"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"

	"golang.org/x/crypto/bcrypt"
)

// UserStore инкапсулирует операции с пользователями.
// Принимает интерфейс db.DBTX, что позволяет мокать зависимость
// в unit-тестах через sqlmock или stub-реализации.
type UserStore struct {
	DB db.DBTX
}

// NewUserStore создает новый экземпляр UserStore.
func NewUserStore(database db.DBTX) *UserStore {
	return &UserStore{DB: database}
}

// GetUserByID возвращает пользователя по ID.
func (s *UserStore) GetUserByID(id int) (models.User, error) {
	var u models.User
	row := s.DB.QueryRow(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE id=$1`, id)
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}
	return u, nil
}

// GetUserByEmail возвращает пользователя по email.
func (s *UserStore) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	row := s.DB.QueryRow(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}
	return u, nil
}

// RegisterUser создает нового пользователя с хешированным паролем.
func (s *UserStore) RegisterUser(email, fullName, role, password string) (models.User, error) {
	role = normalizeRole(role)
	if !roleIsValid(role) {
		return models.User{}, ErrInvalidRole
	}

	// Проверяем уникальность email
	if _, err := s.GetUserByEmail(email); err == nil {
		return models.User{}, ErrDuplicateUser
	} else if err != ErrNotFound {
		return models.User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	createdAt := time.Now()
	query := `INSERT INTO users (email, password_hash, full_name, role, created_at)
	          VALUES ($1,$2,$3,$4,$5) RETURNING id`
	var newId int
	if err := s.DB.QueryRow(query, email, string(hashed), fullName, role, createdAt).Scan(&newId); err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:        newId,
		Email:     email,
		FullName:  fullName,
		Role:      role,
		CreatedAt: createdAt,
	}, nil
}

// AuthenticateUser проверяет email/пароль и возвращает userID.
func (s *UserStore) AuthenticateUser(email, password string) (int, error) {
	u, err := s.GetUserByEmail(email)
	if err != nil {
		if err == ErrNotFound {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return 0, ErrInvalidCredentials
	}

	return u.Id, nil
}

// normalizeRole приводит роль к нижнему регистру и обрезает пробелы.
func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}
