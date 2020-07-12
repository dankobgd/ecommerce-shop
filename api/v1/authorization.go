package apiv1

import (
	"context"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/app"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgAdminRequired = &i18n.Message{ID: "api.admin_session_required.app_error", Other: "insufficient permissions"}
)

// SessionRequired requires session to access the resource
func (a *API) SessionRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := a.app.TokenValid(r)
		if err != nil {
			respondError(w, err)
			return
		}

		ad, err := a.app.ExtractTokenMetadata(r)
		if err != nil {
			respondError(w, err)
			return
		}
		if _, err := a.app.GetAuth(ad); err != nil {
			respondError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), app.AccessDataCtxKey, ad)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// AdminSessionRequired requires admin role to access the resource
func (a *API) AdminSessionRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := a.app.TokenValid(r)
		if err != nil {
			respondError(w, err)
			return
		}

		ad, err := a.app.ExtractTokenMetadata(r)
		if err != nil {
			respondError(w, err)
			return
		}
		if _, err := a.app.GetAuth(ad); err != nil {
			respondError(w, err)
			return
		}

		if ad.Role != model.AdminRole {
			respondError(w, model.NewAppErr("AdminSessionRequired", model.ErrUnauthorized, locale.GetUserLocalizer("en"), msgAdminRequired, http.StatusForbidden, nil))
			return
		}

		ctx := context.WithValue(r.Context(), app.AccessDataCtxKey, ad)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
