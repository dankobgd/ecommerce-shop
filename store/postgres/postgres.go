package postgres

import (
	"log"

	_ "github.com/jackc/pgx/stdlib" // pg driver
	"github.com/jmoiron/sqlx"
)

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
