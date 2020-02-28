package postgres

import (
	"log"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib" // pg driver
	"github.com/jmoiron/sqlx"
)

// postgres error codes
const (
	uniqueConstraintViolation = "23505"
)

// IsUniqueConstraintError checks for postgres unique constraint error code
func IsUniqueConstraintError(err error) bool {
	if pqErr, ok := err.(pgx.PgError); ok && pqErr.Code == uniqueConstraintViolation {
		return true
	}
	return false
}

// PgStore has the pg db driver
type PgStore struct {
	db *sqlx.DB
}

// Connect establishes connection to postgres db
func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", "host=localhost user=test password=test dbname=ecommerce sslmode=disable")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return db, nil
}

// NewStore initializes postgres based store
func NewStore(db *sqlx.DB) *PgStore {
	return &PgStore{db: db}
}
