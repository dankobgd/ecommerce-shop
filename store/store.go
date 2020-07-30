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
	ProductTag() ProductTagStore
	ProductImage() ProductImageStore
	Order() OrderStore
	OrderDetail() OrderDetailStore
	Address() AddressStore
}

// UserStore ris the user store
type UserStore interface {
	BulkInsert([]*model.User) *model.AppErr
	Save(*model.User) (*model.User, *model.AppErr)
	Get(id int64) (*model.User, *model.AppErr)
	GetAll() ([]*model.User, *model.AppErr)
	GetByEmail(email string) (*model.User, *model.AppErr)
	Update(id int64, u *model.User) (*model.User, *model.AppErr)
	Delete(id int64) *model.AppErr
	UpdateAvatar(id int64, url *string, publicID *string) (*string, *string, *model.AppErr)
	DeleteAvatar(id int64) *model.AppErr
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
	Save(p *model.Product) (*model.Product, *model.AppErr)
	Get(id int64) (*model.Product, *model.AppErr)
	ListByIDS(ids []int64) ([]*model.Product, *model.AppErr)
	GetAll() ([]*model.Product, *model.AppErr)
	Update(id int64, u *model.Product) (*model.Product, *model.AppErr)
	Delete(id int64) *model.AppErr
}

// ProductTagStore is the product tag store
type ProductTagStore interface {
	BulkInsert(tags []*model.ProductTag) ([]int64, *model.AppErr)
	Get(id int64) (*model.ProductTag, *model.AppErr)
	GetAll(pid int64) ([]*model.ProductTag, *model.AppErr)
	Update(id int64, tag *model.ProductTag) (*model.ProductTag, *model.AppErr)
	Delete(id int64) *model.AppErr
}

// ProductImageStore is the product image store
type ProductImageStore interface {
	BulkInsert(imgs []*model.ProductImage) ([]int64, *model.AppErr)
	Get(id int64) (*model.ProductImage, *model.AppErr)
	GetAll(pid int64) ([]*model.ProductImage, *model.AppErr)
	Update(id int64, img *model.ProductImage) (*model.ProductImage, *model.AppErr)
	Delete(id int64) *model.AppErr
}

// OrderStore is the order store
type OrderStore interface {
	Save(order *model.Order, shipAddr *model.Address, billAddr *model.Address) (*model.Order, *model.AppErr)
	Get(id int64) (*model.Order, *model.AppErr)
	Update(id int64, o *model.Order) (*model.Order, *model.AppErr)
}

// OrderDetailStore is the order detail store
type OrderDetailStore interface {
	BulkInsert(items []*model.OrderDetail) *model.AppErr
	Save(od *model.OrderDetail) (*model.OrderDetail, *model.AppErr)
	Get(id int64) (*model.OrderDetail, *model.AppErr)
}

// AddressStore is the contact address store
type AddressStore interface {
	Save(addr *model.Address, userID int64, addrType model.AddrType) (*model.Address, *model.AppErr)
	Get(id int64) (*model.Address, *model.AppErr)
	Update(id int64, addr *model.Address) (*model.Address, *model.AppErr)
	Delete(id int64) *model.AppErr
}
