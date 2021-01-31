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
	ProductReview() ProductReviewStore
	Order() OrderStore
	OrderDetail() OrderDetailStore
	Address() AddressStore
	Category() CategoryStore
	Brand() BrandStore
	Tag() TagStore
	Promotion() PromotionStore
}

// UserStore ris the user store
type UserStore interface {
	BulkInsert([]*model.User) *model.AppErr
	Save(*model.User) (*model.User, *model.AppErr)
	Get(id int64) (*model.User, *model.AppErr)
	GetAll(limit, offset int) ([]*model.User, *model.AppErr)
	GetByEmail(email string) (*model.User, *model.AppErr)
	Update(id int64, u *model.User) (*model.User, *model.AppErr)
	Delete(id int64) *model.AppErr
	BulkDelete(ids []int) *model.AppErr
	UpdateAvatar(id int64, url *string, publicID *string) (*string, *string, *model.AppErr)
	DeleteAvatar(id int64) *model.AppErr
	VerifyEmail(userID int64) *model.AppErr
	UpdatePassword(userID int64, hashedPassword string) *model.AppErr
	GetAllOrders(userID int64, limit, offset int) ([]*model.Order, *model.AppErr)
	CreateWishlist(userID, productID int64) *model.AppErr
	GetWishlist(userID int64) ([]*model.Product, *model.AppErr)
	DeleteWishlist(userID, productID int64) *model.AppErr
	ClearWishlist(userID int64) *model.AppErr
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
	GetAll(filters map[string][]string, limit, offset int) ([]*model.Product, *model.AppErr)
	GetFeatured(limit, offset int) ([]*model.Product, *model.AppErr)
	Update(id int64, u *model.Product) (*model.Product, *model.AppErr)
	Delete(id int64) *model.AppErr
	BulkDelete(ids []int) *model.AppErr
	GetReviews(id int64) ([]*model.ProductReview, *model.AppErr)
	Search(query string) ([]*model.Product, *model.AppErr)
	GetLatestPricing(pid int64) (*model.ProductPricing, *model.AppErr)
	InsertPricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr)
	UpdatePricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr)
}

// ProductTagStore is the product tag store
type ProductTagStore interface {
	BulkInsert(tags []*model.ProductTag) *model.AppErr
	Save(pid int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr)
	Get(pid, tid int64) (*model.ProductTag, *model.AppErr)
	GetAll(pid int64) ([]*model.ProductTag, *model.AppErr)
	Update(pid, tid int64, tag *model.ProductTag) (*model.ProductTag, *model.AppErr)
	Replace(pid int64, tagIDs []int) ([]*model.ProductTag, *model.AppErr)
	Delete(pid, id int64) *model.AppErr
	BulkDelete(pid int64, ids []int) *model.AppErr
}

// ProductImageStore is the product image store
type ProductImageStore interface {
	BulkInsert(imgs []*model.ProductImage) *model.AppErr
	Save(pid int64, img *model.ProductImage) (*model.ProductImage, *model.AppErr)
	Get(pid, id int64) (*model.ProductImage, *model.AppErr)
	GetAll(pid int64) ([]*model.ProductImage, *model.AppErr)
	Update(pid, id int64, img *model.ProductImage) (*model.ProductImage, *model.AppErr)
	Delete(pid, id int64) *model.AppErr
	BulkDelete(pid int64, ids []int) *model.AppErr
}

// ProductReviewStore is the review store
type ProductReviewStore interface {
	BulkInsert(reviews []*model.ProductReview) *model.AppErr
	Save(pid int64, review *model.ProductReview) (*model.ProductReview, *model.AppErr)
	Get(pid, rid int64) (*model.ProductReview, *model.AppErr)
	GetAll(pid int64) ([]*model.ProductReview, *model.AppErr)
	Update(pid, rid int64, rev *model.ProductReview) (*model.ProductReview, *model.AppErr)
	Delete(pid, rid int64) *model.AppErr
	BulkDelete(pid int64, ids []int) *model.AppErr
}

// OrderStore is the order store
type OrderStore interface {
	Save(order *model.Order) (*model.Order, *model.AppErr)
	Get(id int64) (*model.Order, *model.AppErr)
	GetAll(limit, offset int) ([]*model.Order, *model.AppErr)
	Update(id int64, order *model.Order) (*model.Order, *model.AppErr)
	Delete(id int64) *model.AppErr
}

// OrderDetailStore is the order detail store
type OrderDetailStore interface {
	BulkInsert(items []*model.OrderDetail) *model.AppErr
	Save(od *model.OrderDetail) (*model.OrderDetail, *model.AppErr)
	Get(id int64) (*model.OrderDetail, *model.AppErr)
	GetAll(limit, offset int) ([]*model.OrderDetail, *model.AppErr)
}

// AddressStore is the contact address store
type AddressStore interface {
	Save(addr *model.Address, userID int64) (*model.Address, *model.AppErr)
	Get(userID, addressID int64) (*model.Address, *model.AppErr)
	GetAll(userID int64) ([]*model.Address, *model.AppErr)
	Update(addressID int64, addr *model.Address) (*model.Address, *model.AppErr)
	Delete(addressID int64) *model.AppErr
}

// CategoryStore is the category store
type CategoryStore interface {
	BulkInsert(categories []*model.Category) *model.AppErr
	Save(c *model.Category) (*model.Category, *model.AppErr)
	Get(id int64) (*model.Category, *model.AppErr)
	GetAll(limit, offset int) ([]*model.Category, *model.AppErr)
	GetFeatured(limit, offset int) ([]*model.Category, *model.AppErr)
	Update(id int64, addr *model.Category) (*model.Category, *model.AppErr)
	Delete(id int64) *model.AppErr
	BulkDelete(ids []int) *model.AppErr
}

// BrandStore is the brand store
type BrandStore interface {
	BulkInsert(brands []*model.Brand) *model.AppErr
	Save(b *model.Brand) (*model.Brand, *model.AppErr)
	Get(id int64) (*model.Brand, *model.AppErr)
	GetAll(limit, offset int) ([]*model.Brand, *model.AppErr)
	Update(id int64, addr *model.Brand) (*model.Brand, *model.AppErr)
	Delete(id int64) *model.AppErr
	BulkDelete(ids []int) *model.AppErr
}

// TagStore is the tag store
type TagStore interface {
	BulkInsert(tags []*model.Tag) *model.AppErr
	Save(t *model.Tag) (*model.Tag, *model.AppErr)
	Get(id int64) (*model.Tag, *model.AppErr)
	GetAll(limit, offset int) ([]*model.Tag, *model.AppErr)
	Update(id int64, addr *model.Tag) (*model.Tag, *model.AppErr)
	Delete(id int64) *model.AppErr
	BulkDelete(ids []int) *model.AppErr
}

// PromotionStore is the promotion store
type PromotionStore interface {
	BulkInsert(promotions []*model.Promotion) *model.AppErr
	Save(p *model.Promotion) (*model.Promotion, *model.AppErr)
	Get(code string) (*model.Promotion, *model.AppErr)
	GetAll(limit, offset int) ([]*model.Promotion, *model.AppErr)
	Update(code string, p *model.Promotion) (*model.Promotion, *model.AppErr)
	Delete(code string) *model.AppErr
	BulkDelete(codes []string) *model.AppErr
	InsertDetail(pd *model.PromotionDetail) (*model.PromotionDetail, *model.AppErr)
	IsValid(code string) *model.AppErr
	IsUsed(code string, userID int64) *model.AppErr
}
