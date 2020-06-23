package apiv1

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgInvalidToken         = &i18n.Message{ID: "model.access_token_verify.json.app_error", Other: "token is invalid or has already expired"}
	msgUserFromJSON         = &i18n.Message{ID: "api.user.create_user.json.app_error", Other: "could not decode user json data"}
	msgRefreshTokenFromJSON = &i18n.Message{ID: "api.user.create_user.json.app_error", Other: "could not decode token json data"}
	msgInvalidEmail         = &i18n.Message{ID: "api.sendUserVerificationEmail.email.app_error", Other: "invalid email provided"}
	msgInvalidPassword      = &i18n.Message{ID: "api.resetUserPassword.password.app_error", Other: "invalid password provided"}
)

// InitUser inits the user routes
func InitUser(a *API) {
	a.BaseRoutes.Users.Post("/", a.createUser)
	a.BaseRoutes.Users.Post("/login", a.login)
	a.BaseRoutes.Users.Post("/logout", a.AuthRequired(a.logout))
	a.BaseRoutes.Users.Post("/token/refresh", a.refresh)
	a.BaseRoutes.Users.Post("/email/verify", a.verifyUserEmail)
	a.BaseRoutes.Users.Post("/email/verify/send", a.sendVerificationEmail)
	a.BaseRoutes.Users.Post("/password/reset", a.resetUserPassword)
	a.BaseRoutes.Users.Post("/password/reset/send", a.sendPasswordResetEmail)

	a.BaseRoutes.Users.Get("/protected", a.AuthRequired(a.protected))
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

	tokenMeta, err := a.app.IssueTokens(user)
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

	tokenMeta, err := a.app.IssueTokens(user)
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
	a.app.DeleteSessionCookies(w)
	ad, err := a.app.ExtractTokenMetadata(r)
	if err != nil {
		respondError(w, err)
		return
	}
	deleted, err := a.app.DeleteAuth(ad.AccessUUID)
	if err != nil || deleted == 0 {
		respondError(w, err)
		return
	}
	respondOK(w)
}

func (a *API) refresh(w http.ResponseWriter, r *http.Request) {
	rt, e := model.RefreshTokenFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("refresh", model.ErrInternal, locale.GetUserLocalizer("en"), msgRefreshTokenFromJSON, http.StatusInternalServerError, nil))
		return
	}

	meta, err := a.app.RefreshToken(rt)
	if err != nil {
		respondError(w, err)
		return
	}

	a.app.AttachSessionCookies(w, meta)
	respondOK(w)
}

func (a *API) protected(w http.ResponseWriter, r *http.Request) {
	ad, err := a.app.ExtractTokenMetadata(r)
	if err != nil {
		respondError(w, err)
		return
	}
	userID, err := a.app.GetAuth(ad)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"userID": userID})
}

func (a *API) sendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	props := model.MapStrStrFromJSON(r.Body)
	email := props["email"]
	email = model.NormalizeEmail(email)

	if len(email) == 0 || !model.IsValidEmail(email) {
		respondError(w, model.NewAppErr("api.sendVerificationEmail", model.ErrInvalid, locale.GetUserLocalizer("en"), msgInvalidEmail, http.StatusBadRequest, nil))
		return
	}

	user, err := a.app.GetUserByEmail(email)
	if err != nil {
		// don't leak whether email is valid and exists - maybe for demonstration return some err
		respondOK(w)
		return
	}

	if err := a.app.SendVerificationEmail(user, email); err != nil {
		// don't leak whether email is valid and exists - maybe for demonstration return some err
		respondOK(w)
		return
	}

	respondOK(w)
}

func (a *API) verifyUserEmail(w http.ResponseWriter, r *http.Request) {
	props := model.MapStrStrFromJSON(r.Body)
	token := props["token"]

	if len(token) == 0 {
		respondError(w, model.NewAppErr("api.sendVerificationEmail", model.ErrInvalid, locale.GetUserLocalizer("en"), msgInvalidToken, http.StatusBadRequest, nil))
		return
	}

	if err := a.app.VerifyUserEmail(token); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}

func (a *API) sendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	props := model.MapStrStrFromJSON(r.Body)
	email := props["email"]
	email = model.NormalizeEmail(email)

	if len(email) == 0 || !model.IsValidEmail(email) {
		respondError(w, model.NewAppErr("api.sendPasswordResetEmail", model.ErrInvalid, locale.GetUserLocalizer("en"), msgInvalidEmail, http.StatusBadRequest, nil))
		return
	}

	if err := a.app.SendPasswordResetEmail(email); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) resetUserPassword(w http.ResponseWriter, r *http.Request) {
	props := model.MapStrStrFromJSON(r.Body)
	token := props["token"]
	oldPassword := props["old_password"]
	newPassword := props["new_password"]

	if len(oldPassword) == 0 || len(newPassword) == 0 {
		respondError(w, model.NewAppErr("api.resetUserPassword", model.ErrInvalid, locale.GetUserLocalizer("en"), msgInvalidPassword, http.StatusBadRequest, nil))
		return
	}

	if err := a.app.ResetUserPassword(token, oldPassword, newPassword); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}
