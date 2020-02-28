package apiv1

import (
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgUserFromJSON = &i18n.Message{ID: "api.user.create_user.json.app_error", Other: "could not decode user json data"}
)

// InitUser inits the user routes
func InitUser(a *API) {
	a.BaseRoutes.Users.Post("/", a.createUser)
	a.BaseRoutes.Users.Post("/login", a.login)
	a.BaseRoutes.Users.Post("/logout", a.AuthRequired(a.logout))
	a.BaseRoutes.Users.Get("/test", a.AuthRequired(a.test))
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	u, e := model.UserFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createUser", model.ErrInternal, locale.GetUserLocalizer("en"), msgUserFromJSON, http.StatusInternalServerError, nil))
		return
	}

	user, err := a.app.CreateUser(u)
	if err != nil {
		respondError(w, err)
		return
	}

	tokenMeta, err := a.app.IssueTokens(user.ID)
	if err != nil {
		respondError(w, err)
	}

	if err := a.app.SaveAuth(user.ID, tokenMeta); err != nil {
		respondError(w, err)
	}

	a.app.AttachSessionCookies(w, tokenMeta)

	respondJSON(w, http.StatusCreated, user)
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	u, e := model.UserLoginFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("login", model.ErrInternal, locale.GetUserLocalizer("en"), msgUserFromJSON, http.StatusInternalServerError, nil))
		return
	}

	user, err := a.app.Login(u)
	if err != nil {
		respondError(w, err)
		return
	}

	tokenMeta, err := a.app.IssueTokens(user.ID)
	if err != nil {
		respondError(w, err)
	}

	if err := a.app.SaveAuth(user.ID, tokenMeta); err != nil {
		respondError(w, err)
	}

	a.app.AttachSessionCookies(w, tokenMeta)

	respondJSON(w, http.StatusOK, user)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	ad, err := a.app.ExtractTokenMetadata(r)
	if err != nil {
		w.Write([]byte("unauthorized"))
		return
	}
	deleted, err := a.app.DeleteAuth(ad.AccessUUID)
	if err != nil || deleted == 0 {
		w.Write([]byte("unauthorized"))
		return
	}
	w.Write([]byte("success logout"))
}

func (a *API) test(w http.ResponseWriter, r *http.Request) {
	ad, err := a.app.ExtractTokenMetadata(r)
	if err != nil {
		w.Write([]byte("unaothorized"))
		return
	}
	userID, err := a.app.GetAuth(ad)
	if err != nil {
		w.Write([]byte("unaothorized"))
		return
	}

	w.Write([]byte("userID: " + strconv.FormatInt(userID, 10)))
}
