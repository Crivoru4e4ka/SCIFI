package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project-MVP/db"
	"project-MVP/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
	ErrDuplicateUser      = errors.New("user with this email already exists")
)

var allowedRoles = map[string]bool{
	"admin": true, // Системный администратор
	"user":  true, // Обычный сотрудник/студент
	"guest": true, // Внешний эксперт/рецензент
}

func roleIsValid(role string) bool {
	return allowedRoles[strings.ToLower(strings.TrimSpace(role))]
}

func GetUsers() ([]models.User, error) {
	rows, err := db.DB.Query(`SELECT id, email, full_name, role, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Email, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetUserByID(id int) (models.User, error) {
	var u models.User
	row := db.DB.QueryRow(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE id=$1`, id)
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}
	return u, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	row := db.DB.QueryRow(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		}
		return u, err
	}
	return u, nil
}

func RegisterUser(email, fullName, role, password string) (models.User, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if !roleIsValid(role) {
		return models.User{}, ErrInvalidRole
	}

	if _, err := GetUserByEmail(email); err == nil {
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
	if err := db.DB.QueryRow(query, email, string(hashed), fullName, role, createdAt).Scan(&newId); err != nil {
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

func AuthenticateUser(email, password string) (int, error) {
	u, err := GetUserByEmail(email)
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
