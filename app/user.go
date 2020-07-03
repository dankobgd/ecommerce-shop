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
	msgTokenExpired = &i18n.Message{ID: "model.token.expired.app_error", Other: "token has expired"}
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
	return a.Srv().Store.AccessToken().SaveAuth(userID, meta)
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
	token, err := a.createTokenAndPersist(user.ID, model.TokenTypeEmailVerification, a.Cfg().AuthSettings.EmailVerificationExpiryHours)
	if err != nil {
		return err
	}
	if err := a.SendEmailVerificationEmail(email, token.Token, a.SiteURL(), user.Locale); err != nil {
		return err
	}
	return nil
}

// VerifyUserEmail verifies the user's email account
func (a *App) VerifyUserEmail(tokenString string) *model.AppErr {
	token, err := a.Srv().Store.Token().GetByToken(tokenString)
	if err != nil {
		return err
	}

	if time.Now().After(token.ExpiresAt) {
		return model.NewAppErr("app.VerifyUserEmail", model.ErrInternal, locale.GetUserLocalizer("en"), msgTokenExpired, http.StatusInternalServerError, nil)
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
func (a *App) SendPasswordResetEmail(email string) *model.AppErr {
	user, err := a.Srv().Store.User().GetByEmail(email)
	if err != nil {
		return err
	}
	token, err := a.createTokenAndPersist(user.ID, model.TokenTypePasswordRecovery, a.Cfg().AuthSettings.PasswordResetExpiryHours)
	if err != nil {
		return err
	}
	if err := a.SendPasswordRecoveryEmail(email, user.Username, token, a.SiteURL(), user.Locale); err != nil {
		return err
	}
	return nil
}

// ResetUserPassword resets the user pwd
func (a *App) ResetUserPassword(tokenString, oldPassword, newPassword string) *model.AppErr {
	token, err := a.Srv().Store.Token().GetByToken(tokenString)
	if err != nil {
		return err
	}

	if time.Now().After(token.ExpiresAt) {
		return model.NewAppErr("app.ResetUserPassword", model.ErrInternal, locale.GetUserLocalizer("en"), msgTokenExpired, http.StatusInternalServerError, nil)
	}

	user, err := a.Srv().Store.User().Get(token.UserID)
	if err != nil {
		return err
	}

	if err := a.CheckUserPassword(user, oldPassword); err != nil {
		return err
	}

	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	if err := a.Srv().Store.Token().Delete(token.Token); err != nil {
		zlog.Error("could not delete token", zlog.Int64("user_id", token.UserID), zlog.String("token_type", token.Type), zlog.Err(err))
	}

	go func() {
		if err := a.SendPasswordUpdatedEmail(user.Email, user.Username, token.Token, a.SiteURL(), user.Locale); err != nil {
			zlog.Error("could not send password reset completed email", zlog.Int64("user_id", token.UserID), zlog.Err(err))
		}
	}()

	return nil
}

// UpdatePassword sets the new user password
func (a *App) UpdatePassword(user *model.User, newPassword string) *model.AppErr {
	if err := a.IsValidPassword(newPassword); err != nil {
		return err
	}
	hashed := model.HashPassword(newPassword)
	if err := a.Srv().Store.User().UpdatePassword(user.ID, hashed); err != nil {
		return err
	}
	return nil
}

func (a *App) createTokenAndPersist(userID int64, tokenType model.TokenType, expiryHours ...int) (*model.Token, *model.AppErr) {
	token := model.NewToken(tokenType, userID, expiryHours...)
	if err := a.Srv().Store.Token().Save(token); err != nil {
		return nil, err
	}
	return token, nil
}
