package app

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgTokenExpired     = &i18n.Message{ID: "model.token.expired.app_error", Other: "token has expired"}
	msgUploadUserAvatar = &i18n.Message{ID: "app.upload_user_avatar.app_error", Other: "could not upload user avatar"}
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
		a.Log().Error(err.Error(), zlog.Err(err))
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
		a.Log().Error(err.Error(), zlog.Err(err))
		return nil, err
	}
	if err := a.CheckUserPassword(user, u.Password); err != nil {
		return nil, err
	}

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

// GetUserByIDWithPassword gets the user by his id and includes pwd (for some checks for the old pwd)
func (a *App) GetUserByIDWithPassword(id int64) (*model.User, *model.AppErr) {
	user, err := a.Srv().Store.User().Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID gets the user by his id
func (a *App) GetUserByID(id int64) (*model.User, *model.AppErr) {
	user, err := a.Srv().Store.User().Get(id)
	if err != nil {
		return nil, err
	}
	user.Sanitize(map[string]bool{})
	return user, nil
}

// GetUserByEmail gets the user by his email
func (a *App) GetUserByEmail(email string) (*model.User, *model.AppErr) {
	user, err := a.Srv().Store.User().GetByEmail(email)
	if err != nil {
		return nil, err
	}
	user.Sanitize(map[string]bool{})
	return user, nil
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
func (a *App) ResetUserPassword(tokenString, newPassword string) *model.AppErr {
	token, err := a.Srv().Store.Token().GetByToken(tokenString)
	if err != nil {
		return err
	}

	if time.Now().After(token.ExpiresAt) {
		return model.NewAppErr("app.ResetUserPassword", model.ErrInternal, locale.GetUserLocalizer("en"), msgTokenExpired, http.StatusInternalServerError, nil)
	}

	user, err := a.GetUserByID(token.UserID)
	if err != nil {
		return err
	}

	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	if err := a.Srv().Store.Token().Delete(token.Token); err != nil {
		a.Log().Error("could not delete token", zlog.Int64("user_id", token.UserID), zlog.String("token_type", token.Type), zlog.Err(err))
	}

	go func() {
		if err := a.SendPasswordUpdatedEmail(user.Email, user.Username, a.SiteURL(), user.Locale); err != nil {
			zlog.Error("could not send password reset completed email", zlog.Int64("user_id", token.UserID), zlog.Err(err))
		}
	}()

	return nil
}

// ChangeUserPassword updates the user pwd from the app
func (a *App) ChangeUserPassword(uid int64, oldPassword, newPassword string) *model.AppErr {
	user, err := a.GetUserByIDWithPassword(uid)
	if err != nil {
		return err
	}

	if err := a.CheckUserPassword(user, oldPassword); err != nil {
		return err
	}

	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	go func() {
		if err := a.SendPasswordUpdatedEmail(user.Email, user.Username, a.SiteURL(), user.Locale); err != nil {
			a.Log().Error("could not send password reset completed email", zlog.Int64("user_id", uid), zlog.Err(err))
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

// PatchUserProfile patches the user profile
func (a *App) PatchUserProfile(id int64, patch *model.UserPatch) (*model.User, *model.AppErr) {
	old, err := a.Srv().Store.User().Get(id)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	if err := patch.Validate(); err != nil {
		return nil, err
	}
	user, err := a.Srv().Store.User().Update(id, old)
	if err != nil {
		return nil, err
	}

	user.Sanitize(map[string]bool{})
	return user, nil
}

// DeleteUser soft deletes the user account
func (a *App) DeleteUser(id int64) *model.AppErr {
	return a.Srv().Store.User().Delete(id)
}

// UploadUserAvatar uploads the user profile image and returns the avatar url
// it returns the url, public id and error
func (a *App) UploadUserAvatar(userID int64, f multipart.File, fh *multipart.FileHeader) (*string, *string, *model.AppErr) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return model.NewString(""), model.NewString(""), model.NewAppErr("uploadUserAvatar", model.ErrInternal, locale.GetUserLocalizer("en"), msgUploadUserAvatar, http.StatusInternalServerError, nil)
	}

	details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
	if uErr != nil {
		return model.NewString(""), model.NewString(""), uErr
	}

	return a.Srv().Store.User().UpdateAvatar(userID, model.NewString(details.SecureURL), model.NewString(details.PublicID))
}

// DeleteUserAvatar deletes the user profile image
func (a *App) DeleteUserAvatar(userID int64, publicID string) *model.AppErr {
	go func() {
		if err := a.DeleteImage(publicID); err != nil {
			a.Log().Error("could not delete user avatar from cloudinary", zlog.Int64("user_id", userID), zlog.String("public_id", publicID), zlog.Err(err))
		}
	}()

	return a.Srv().Store.User().DeleteAvatar(userID)
}

// CreateUserAddress creates the user addresss
func (a *App) CreateUserAddress(addr *model.Address, userID int64) (*model.Address, *model.AppErr) {
	if err := addr.Validate(); err != nil {
		return nil, err
	}

	var latitude, longitude float64

	geocode, err := a.GetAddressGeocodeResult(addr)
	if err != nil {
		a.Log().Warn(err.Message, zlog.Err(err), zlog.Any("geocode_result", geocode))
		addr.Latitude = model.NewFloat64(latitude)
		addr.Longitude = model.NewFloat64(longitude)
	} else {
		lat, _ := strconv.ParseFloat(geocode.Lat, 64)
		lon, _ := strconv.ParseFloat(geocode.Lon, 64)
		addr.Latitude = &lat
		addr.Longitude = &lon
	}

	addr.PreSave()
	return a.Srv().Store.Address().Save(addr, userID)
}

// GetUserAddress gets the user addresss
func (a *App) GetUserAddress(userID, addressID int64) (*model.Address, *model.AppErr) {
	return a.Srv().Store.Address().Get(userID, addressID)
}

// GetUserAddresses gets the user addresses
func (a *App) GetUserAddresses(userID int64) ([]*model.Address, *model.AppErr) {
	return a.Srv().Store.Address().GetAll(userID)
}

// PatchUserAddress patches the user addresss
func (a *App) PatchUserAddress(userID, addressID int64, patch *model.AddressPatch) (*model.Address, *model.AppErr) {
	old, err := a.Srv().Store.Address().Get(userID, addressID)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	uaddress, err := a.Srv().Store.Address().Update(addressID, old)
	if err != nil {
		return nil, err
	}

	return uaddress, nil
}

// DeleteUserAddress hard deletes the user address
func (a *App) DeleteUserAddress(id int64) *model.AppErr {
	return a.Srv().Store.Address().Delete(id)
}

// GetOrdersForUser gets all user orders
func (a *App) GetOrdersForUser(uid int64, limit, offset int) ([]*model.Order, *model.AppErr) {
	return a.Srv().Store.User().GetAllOrders(uid, limit, offset)
}

// CreateWishlistForUser adds new product t the wishlist
func (a *App) CreateWishlistForUser(uid, pid int64) *model.AppErr {
	return a.Srv().Store.User().CreateWishlist(uid, pid)
}

// GetWishlistForUser gets all wishlist products for user
func (a *App) GetWishlistForUser(uid int64) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.User().GetWishlist(uid)
}

// DeleteWishlistForUser gets all user orders
func (a *App) DeleteWishlistForUser(uid, pid int64) *model.AppErr {
	return a.Srv().Store.User().DeleteWishlist(uid, pid)
}

// ClearWishlistForUser gets all user orders
func (a *App) ClearWishlistForUser(uid int64) *model.AppErr {
	return a.Srv().Store.User().ClearWishlist(uid)
}
