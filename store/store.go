package store

import (
	"github.com/dankobgd/ecommerce-shop/model"
)

// Store represents all stores
type Store interface {
	User() UserStore
	AccessToken() AccessTokenStore
}

// UserStore ris the user store
type UserStore interface {
	Save(*model.User) (*model.User, *model.AppErr)
	Get(id int) (*model.User, *model.AppErr)
	GetAll() ([]*model.User, *model.AppErr)
	GetByEmail(email string) (*model.User, *model.AppErr)
	Update(id int, u *model.User) (*model.User, *model.AppErr)
	Delete(id int) (*model.User, *model.AppErr)
}

// AccessTokenStore is the access token store
type AccessTokenStore interface {
	SaveAuth(userID int64, meta *model.TokenMetadata) *model.AppErr
	GetAuth(ad *model.AccessData) (int64, *model.AppErr)
	DeleteAuth(uuid string) (int64, *model.AppErr)
}
