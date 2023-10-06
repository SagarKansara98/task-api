package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect connect with database
func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	// setting ideal connection
	// max connection

	return db, err
}
