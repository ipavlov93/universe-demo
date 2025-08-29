package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ipavlov93/universe-demo/product-sv/internal/config"
)

type PostgresAdapter struct {
	sqlxDB *sqlx.DB
}

// NewWithConfig tries to connect to Postgres database. Otherwise, it returns an error.
func NewWithConfig(cfg config.PostgresConfig) (PostgresAdapter, error) {
	return NewPostgresAdapter(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName, cfg.Options)
}

// NewPostgresAdapter tries to connect to Postgres database. Otherwise, it returns an error.
func NewPostgresAdapter(host string, port int, user, password, dbname, option string) (PostgresAdapter, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		host, port, user, password, dbname, option)

	connection, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return PostgresAdapter{}, err
	}

	return PostgresAdapter{sqlxDB: connection}, nil
}

func (db *PostgresAdapter) GetConnection() *sqlx.DB {
	return db.sqlxDB
}

func (db *PostgresAdapter) CloseConnection() error {
	if db.sqlxDB == nil {
		return nil
	}
	return db.sqlxDB.Close()
}
