package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgGetToken            = &i18n.Message{ID: "store.postgres.token.get_by_token.app_error", Other: "could not get token from db"}
	msgSaveToken           = &i18n.Message{ID: "store.postgres.token.save.app_error", Other: "could not save token to db"}
	msgCleanup             = &i18n.Message{ID: "store.postgres.token.cleanup.app_error", Other: "could not cleanup all tokens"}
	msgRemoveAllTokensType = &i18n.Message{ID: "store.postgres.token.RemoveAllTokensByType.app_error", Other: "could not remove all tokens by type"}
)

// PgTokenStore is the postgres implementation
type PgTokenStore struct {
	PgStore
}

// NewPgTokenStore creates the new token store
func NewPgTokenStore(pgst *PgStore) store.TokenStore {
	return &PgTokenStore{*pgst}
}

// Save saves the token
func (s PgTokenStore) Save(token *model.Token) *model.AppErr {
	q := "INSERT INTO public.token (user_id, token, type, created_at, expires_at) values(:user_id, :token, :type, :created_at, :expires_at)"
	if _, err := s.db.NamedExec(q, token); err != nil {
		return model.NewAppErr("PgTokenStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveToken, http.StatusInternalServerError, nil)
	}
	return nil
}

// GetByToken searches the token by the token string
func (s PgTokenStore) GetByToken(token string) (*model.Token, *model.AppErr) {
	var tkn model.Token
	if err := s.db.Get(&tkn, "SELECT * FROM public.token where token = $1", token); err != nil {
		return nil, model.NewAppErr("PgUserStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetToken, http.StatusInternalServerError, nil)
	}
	return &tkn, nil
}

// Delete deletes the single token
func (s PgTokenStore) Delete(token string) *model.AppErr {
	m := map[string]interface{}{"token": token}
	if _, err := s.db.NamedExec("DELETE FROM public.token WHERE token = :token", m); err != nil {
		return model.NewAppErr("PgTokenStore.DeleteToken", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteToken, http.StatusInternalServerError, nil)
	}
	return nil
}

// Cleanup deletes all the tokens
func (s PgTokenStore) Cleanup() *model.AppErr {
	if _, err := s.db.Exec("DELETE FROM public.token"); err != nil {
		return model.NewAppErr("PgTokenStore.Cleanup", model.ErrInternal, locale.GetUserLocalizer("en"), msgCleanup, http.StatusInternalServerError, nil)
	}
	return nil
}

// RemoveByType deletes all tokens of same type
func (s PgTokenStore) RemoveByType(tokenType model.TokenType) *model.AppErr {
	m := map[string]interface{}{"type": tokenType.String()}
	if _, err := s.db.Exec("DELETE FROM public.token WHERE type = :type", m); err != nil {
		return model.NewAppErr("PgTokenStore.RemoveByType", model.ErrInternal, locale.GetUserLocalizer("en"), msgRemoveAllTokensType, http.StatusInternalServerError, nil)
	}
	return nil
}
