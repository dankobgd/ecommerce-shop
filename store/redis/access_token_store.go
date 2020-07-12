package redis

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgSaveAuth   = &i18n.Message{ID: "store.redis.access_token.save_auth.app_error", Other: "could not save auth data"}
	msgGetAuth    = &i18n.Message{ID: "store.redis.access_token.get_auth.app_error", Other: "auth token is invalid or has already expired"}
	msgDeleteAuth = &i18n.Message{ID: "store.redis.access_token.delete_auth.app_error", Other: "could not delete auth data"}
)

// RdAccessTokenStore is the redis implementation
type RdAccessTokenStore struct {
	RdStore
}

// NewRedisAccessTokenStore creates the new access token store
func NewRedisAccessTokenStore(rdst *RdStore) store.AccessTokenStore {
	return &RdAccessTokenStore{*rdst}
}

// SaveAuth saves the user auth meta data
func (s RdAccessTokenStore) SaveAuth(userID int64, meta *model.TokenMetadata) *model.AppErr {
	c := context.TODO()
	now := time.Now()

	err := s.client.Set(c, meta.AccessUUID, userID, meta.AccessExpires.Sub(now)).Err()
	if err != nil {
		return model.NewAppErr("RdAccessTokenStore.SaveAuth", model.ErrConflict, locale.GetUserLocalizer("en"), msgSaveAuth, http.StatusInternalServerError, nil)
	}
	err2 := s.client.Set(c, meta.RefreshUUID, userID, meta.RefreshExpires.Sub(now)).Err()
	if err2 != nil {
		return model.NewAppErr("RdAccessTokenStore.SaveAuth", model.ErrConflict, locale.GetUserLocalizer("en"), msgSaveAuth, http.StatusInternalServerError, nil)
	}

	return nil
}

// GetAuth gets the user auth meta data. err = token expired
func (s RdAccessTokenStore) GetAuth(ad *model.AccessData) (int64, *model.AppErr) {
	uid, err := s.client.Get(context.TODO(), ad.AccessUUID).Result()
	if err != nil {
		return 0, model.NewAppErr("RdAccessTokenStore.GetAuth", model.ErrConflict, locale.GetUserLocalizer("en"), msgGetAuth, http.StatusInternalServerError, nil)
	}
	userID, _ := strconv.ParseInt(uid, 10, 64)
	return userID, nil
}

// DeleteAuth deletes the token
func (s RdAccessTokenStore) DeleteAuth(uuid string) (int64, *model.AppErr) {
	deleted, err := s.client.Del(context.TODO(), uuid).Result()
	if err != nil {
		return 0, model.NewAppErr("RdAccessTokenStore.DeleteAuth", model.ErrConflict, locale.GetUserLocalizer("en"), msgDeleteAuth, http.StatusInternalServerError, nil)
	}
	return deleted, nil
}
