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

// Product returns the Product tag store implementation
func (s *Supplier) Product() store.ProductStore {
	return postgres.NewPgProductStore(s.Pgst)
}

// ProductTag returns the Product store implementation
func (s *Supplier) ProductTag() store.ProductTagStore {
	return postgres.NewPgProductTagStore(s.Pgst)
}

// ProductImage returns the Product image store implementation
func (s *Supplier) ProductImage() store.ProductImageStore {
	return postgres.NewPgProductImageStore(s.Pgst)
}

// Order returns the Order store implementation
func (s *Supplier) Order() store.OrderStore {
	return postgres.NewPgOrderStore(s.Pgst)
}

// OrderDetail returns the OrderDetail store implementation
func (s *Supplier) OrderDetail() store.OrderDetailStore {
	return postgres.NewPgOrderDetailStore(s.Pgst)
}

// Address returns the Address store implementation
func (s *Supplier) Address() store.AddressStore {
	return postgres.NewPgAddressStore(s.Pgst)
}

// Category returns the Category store implementation
func (s *Supplier) Category() store.CategoryStore {
	return postgres.NewPgCategoryStore(s.Pgst)
}

// Brand returns the Brand store implementation
func (s *Supplier) Brand() store.BrandStore {
	return postgres.NewPgBrandStore(s.Pgst)
}

// Tag returns the Tag store implementation
func (s *Supplier) Tag() store.TagStore {
	return postgres.NewPgTagStore(s.Pgst)
}
