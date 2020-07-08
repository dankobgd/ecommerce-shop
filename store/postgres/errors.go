package postgres

import "github.com/jackc/pgx"

// postgres error codes
const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

// IsUniqueConstraintViolationError checks for postgres unique constraint error code
func IsUniqueConstraintViolationError(err error) bool {
	if pqErr, ok := err.(pgx.PgError); ok && pqErr.Code == uniqueViolation {
		return true
	}
	return false
}

// IsForeignKeyConstraintViolationError checks for postgres unique constraint error code
func IsForeignKeyConstraintViolationError(err error) bool {
	if pqErr, ok := err.(pgx.PgError); ok && pqErr.Code == foreignKeyViolation {
		return true
	}
	return false
}
