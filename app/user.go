package app

import (
	"net/http"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgVerifyTokenExpired = &i18n.Message{ID: "model.token.expired.app_error", Other: "token has expired"}
)

// CreateUser creates the new user in the system
func (a *App) CreateUser(u *model.User) (*model.User, *model.AppErr) {
	rawpw := u.Password
	u.PreSave()
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := a.IsValidPassword(rawpw); err != nil {
		return nil, err
	}

	user, err := a.Srv().Store.User().Save(u)

	if err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}

	user.Sanitize(map[string]bool{})
	return user, nil
}

// Login handles the user login
func (a *App) Login(u *model.UserLogin) (*model.User, *model.AppErr) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	user, err := a.Srv().Store.User().GetByEmail(u.Email)
	if err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}
	if err := a.CheckUserPassword(user, u.Password); err != nil {
		return nil, err
	}

	user.Sanitize(map[string]bool{})
	return user, nil
}

// SaveAuth saves the user auth information
func (a *App) SaveAuth(userID int64, meta *model.TokenMetadata) *model.AppErr {
	if err := a.Srv().Store.AccessToken().SaveAuth(userID, meta); err != nil {
		return model.NewAppErr("createUser", model.ErrInternal, locale.GetUserLocalizer("en"), model.MsgInvalidUser, http.StatusInternalServerError, nil)
	}
	return nil
}

// GetAuth gets the auth details information
func (a *App) GetAuth(ad *model.AccessData) (int64, *model.AppErr) {
	return a.Srv().Store.AccessToken().GetAuth(ad)
}

// DeleteAuth deletes the user auth details
func (a *App) DeleteAuth(uuid string) (int64, *model.AppErr) {
	return a.Srv().Store.AccessToken().DeleteAuth(uuid)
}

// GetUserByID gets the user by his id
func (a *App) GetUserByID(id int64) (*model.User, *model.AppErr) {
	return a.Srv().Store.User().Get(id)
}

// GetUserByEmail gets the user by his email
func (a *App) GetUserByEmail(email string) (*model.User, *model.AppErr) {
	return a.Srv().Store.User().GetByEmail(email)
}

// SendVerificationEmail sends the email for veryfing the user account
func (a *App) SendVerificationEmail(user *model.User, email string) *model.AppErr {
	token, err := a.createVerifyEmailToken(user.ID, email)
	if err != nil {
		return err
	}
	if err := a.SendEmailVerificationEmail(email, token.Token, a.SiteURL(), user.Locale); err != nil {
		return err
	}
	return nil
}

func (a *App) createVerifyEmailToken(id int64, email string) (*model.Token, *model.AppErr) {
	token := model.NewToken(model.TokenTypeEmailVerification, id)
	if err := a.Srv().Store.Token().Save(token); err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyUserEmail verifies the user's email account
func (a *App) VerifyUserEmail(tkn string) *model.AppErr {
	token, err := a.Srv().Store.Token().GetByToken(tkn)
	if err != nil {
		return err
	}

	if time.Now().After(token.ExpiresAt) {
		return model.NewAppErr("VerifyUserEmail", model.ErrInternal, locale.GetUserLocalizer("en"), msgVerifyTokenExpired, http.StatusInternalServerError, nil)
	}

	if err := a.Srv().Store.User().VerifyEmail(token.UserID); err != nil {
		return err
	}

	if err := a.Srv().Store.Token().Delete(token.Token); err != nil {
		zlog.Error("could not delete token", zlog.Int64("user_id", token.UserID), zlog.String("token_type", token.Type), zlog.Err(err))
	}

	return nil
}

// SendPasswordResetEmail sends the pwd recovery email
func (a *App) SendPasswordResetEmail() {}

// ResetUserPassword resets the user pwd
func (a *App) ResetUserPassword() {}
