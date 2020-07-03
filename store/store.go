package store

import (
	"github.com/dankobgd/ecommerce-shop/model"
)

// Store represents all stores
type Store interface {
	AccessToken() AccessTokenStore
	User() UserStore
	Token() TokenStore
	Product() ProductStore
}

// UserStore ris the user store
type UserStore interface {
	BulkInsert([]*model.User) *model.AppErr
	Save(*model.User) (*model.User, *model.AppErr)
	Get(id int64) (*model.User, *model.AppErr)
	GetAll() ([]*model.User, *model.AppErr)
	GetByEmail(email string) (*model.User, *model.AppErr)
	Update(id int64, u *model.User) (*model.User, *model.AppErr)
	Delete(id int64) (*model.User, *model.AppErr)
	VerifyEmail(usrerID int64) *model.AppErr
	UpdatePassword(userID int64, hashedPassword string) *model.AppErr
}

// AccessTokenStore is the access token store
type AccessTokenStore interface {
	SaveAuth(userID int64, meta *model.TokenMetadata) *model.AppErr
	GetAuth(ad *model.AccessData) (int64, *model.AppErr)
	DeleteAuth(uuid string) (int64, *model.AppErr)
}

// TokenStore is the access token store
type TokenStore interface {
	Save(token *model.Token) *model.AppErr
	GetByToken(token string) (*model.Token, *model.AppErr)
	Delete(token string) *model.AppErr
	Cleanup() *model.AppErr
	RemoveByType(tokenType model.TokenType) *model.AppErr
}

// ProductStore is the product store
type ProductStore interface {
	BulkInsert([]*model.Product) *model.AppErr
	Save(*model.Product) (*model.Product, *model.AppErr)
	Get(id int64) (*model.Product, *model.AppErr)
	GetAll() ([]*model.Product, *model.AppErr)
	Update(id int64, u *model.Product) (*model.Product, *model.AppErr)
	Delete(id int64) (*model.Product, *model.AppErr)
	InsertTag(tag *model.ProductTag) *model.AppErr
	InsertCategory(category *model.ProductCategory) *model.AppErr
	InsertBrand(brand *model.ProductBrand) *model.AppErr
}
