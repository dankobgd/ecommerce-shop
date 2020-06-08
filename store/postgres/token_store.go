package postgres

import (
	"github.com/dankobgd/ecommerce-shop/store"
)

// PgTokenStore is the postgres implementation
type PgTokenStore struct {
	PgStore
}

// NewPgTokenStore creates the new token store
func NewPgTokenStore(pgst *PgStore) store.TokenStore {
	return &PgTokenStore{*pgst}
}

// Cleanup ...
func (s PgTokenStore) Cleanup() {}

// Get ...
func (s PgTokenStore) Get() {}

// RemoveAllTokensByType ...
func (s PgTokenStore) RemoveAllTokensByType(tokenType string) {}
