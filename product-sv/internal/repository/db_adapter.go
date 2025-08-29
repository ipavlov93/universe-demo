package repository

import (
	"github.com/jmoiron/sqlx"
)

type DBAdapter interface {
	GetConnection() *sqlx.DB
	CloseConnection() error
}
