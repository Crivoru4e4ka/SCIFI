package db

import (
	"database/sql"
)

// DBTX абстрагирует операции с базой данных, позволяя использовать
// как *sql.DB, так и *sql.Tx внутри транзакций.
// Это ключевой интерфейс для unit-тестирования сервисов с sqlmock.
type DBTX interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DBPool расширяет DBTX операциями Begin для транзакций.
// Реализуется *sql.DB.
type DBPool interface {
	DBTX
	Begin() (*sql.Tx, error)
}
