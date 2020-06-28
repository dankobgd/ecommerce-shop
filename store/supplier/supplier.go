package supplier

import (
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/store/postgres"
	"github.com/dankobgd/ecommerce-shop/store/redis"
)

// Supplier contains the stores
type Supplier struct {
	Pgst *postgres.PgStore
	Rdst *redis.RdStore
}

// AccessToken returns the AccessToken store implementation
func (s *Supplier) AccessToken() store.AccessTokenStore {
	return redis.NewRedisAccessTokenStore(s.Rdst)
}

// User returns the User store implementation
func (s *Supplier) User() store.UserStore {
	return postgres.NewPgUserStore(s.Pgst)
}

// Token returns the Token store implementation
func (s *Supplier) Token() store.TokenStore {
	return postgres.NewPgTokenStore(s.Pgst)
}

// Product returns the Product store implementation
func (s *Supplier) Product() store.ProductStore {
	return postgres.NewPgProductStore(s.Pgst)
}
