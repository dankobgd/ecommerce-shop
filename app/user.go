package app

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
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

// GetAuth ...
func (a *App) GetAuth(ad *model.AccessData) (int64, *model.AppErr) {
	return a.Srv().Store.AccessToken().GetAuth(ad)
}

// DeleteAuth ...
func (a *App) DeleteAuth(uuid string) (int64, *model.AppErr) {
	return a.Srv().Store.AccessToken().DeleteAuth(uuid)
}
