package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestUserStore_RegisterUser_Success проверяет успешную регистрацию пользователя.
// Использует sqlmock для имитации запросов к БД без реального подключения.
func TestUserStore_RegisterUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	// Ожидаем проверку email (GetUserByEmail)
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=\$1`).
		WithArgs("test@example.com").
		WillReturnError(sql.ErrNoRows)

	// Ожидаем INSERT с RETURNING id
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("test@example.com", sqlmock.AnyArg(), "Test User", "user", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user, err := store.RegisterUser("test@example.com", "Test User", "user", "password123")

	require.NoError(t, err)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, "user", user.Role)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_RegisterUser_DuplicateEmail проверяет, что регистрация
// с существующим email возвращает ErrDuplicateUser.
func TestUserStore_RegisterUser_DuplicateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	// Пользователь уже существует
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at"}).
		AddRow(1, "test@example.com", "hash", "Existing User", "user", time.Now())
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=\$1`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	_, err = store.RegisterUser("test@example.com", "Test User", "user", "password123")

	assert.ErrorIs(t, err, ErrDuplicateUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_RegisterUser_InvalidRole проверяет валидацию роли
// до обращения к базе данных.
func TestUserStore_RegisterUser_InvalidRole(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	_, err = store.RegisterUser("test@example.com", "Test User", "superuser", "password123")

	assert.ErrorIs(t, err, ErrInvalidRole)
}

// TestUserStore_AuthenticateUser_Success проверяет успешную аутентификацию.
func TestUserStore_AuthenticateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	// Хеш пароля "secret" с bcrypt.DefaultCost
	hash, _ := bcryptHash("secret")
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at"}).
		AddRow(5, "user@example.com", hash, "User", "user", time.Now())

	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=\$1`).
		WithArgs("user@example.com").
		WillReturnRows(rows)

	userID, err := store.AuthenticateUser("user@example.com", "secret")

	require.NoError(t, err)
	assert.Equal(t, 5, userID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_AuthenticateUser_WrongPassword проверяет отказ
// при неверном пароле.
func TestUserStore_AuthenticateUser_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	hash, _ := bcryptHash("correct_password")
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at"}).
		AddRow(5, "user@example.com", hash, "User", "user", time.Now())

	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=\$1`).
		WithArgs("user@example.com").
		WillReturnRows(rows)

	_, err = store.AuthenticateUser("user@example.com", "wrong_password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_AuthenticateUser_UserNotFound проверяет отказ
// при отсутствии пользователя.
func TestUserStore_AuthenticateUser_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email=\$1`).
		WithArgs("unknown@example.com").
		WillReturnError(sql.ErrNoRows)

	_, err = store.AuthenticateUser("unknown@example.com", "password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_GetUserByID_Success проверяет получение пользователя по ID.
func TestUserStore_GetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at"}).
		AddRow(7, "admin@example.com", "hash", "Admin User", "admin", time.Now())

	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE id=\$1`).
		WithArgs(7).
		WillReturnRows(rows)

	user, err := store.GetUserByID(7)

	require.NoError(t, err)
	assert.Equal(t, 7, user.Id)
	assert.Equal(t, "admin", user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserStore_GetUserByID_NotFound проверяет обработку отсутствующего пользователя.
func TestUserStore_GetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := NewUserStore(db)

	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE id=\$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetUserByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// bcryptHash — вспомогательная функция для генерации хеша пароля в тестах.
func bcryptHash(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(fromPassword), err
}
