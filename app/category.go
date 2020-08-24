package app

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgCategorySizeExceeded = &i18n.Message{ID: "app.category.create_category.image_size.app_error", Other: "upload image size exceeded"}
	msgCategoryFileErr      = &i18n.Message{ID: "app.category.create_category.formfile.app_error", Other: "error parsing files"}
)

// CreateCategory creates the new category in the system
func (a *App) CreateCategory(c *model.Category, fh *multipart.FileHeader) (*model.Category, *model.AppErr) {
	if fh.Size > model.FileUploadSizeLimit {
		return nil, model.NewAppErr("CreateCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategorySizeExceeded, http.StatusInternalServerError, nil)
	}
	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("CreateCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategorySizeExceeded, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("CreateCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryFileErr, http.StatusInternalServerError, nil)
	}

	c.PreSave()
	if err := c.Validate(); err != nil {
		return nil, err
	}

	details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}
	c.SetLogoURL(details.SecureURL)

	category, cErr := a.Srv().Store.Category().Save(c)
	if cErr != nil {
		a.Log().Error(cErr.Error(), zlog.Err(cErr))
		return nil, cErr
	}

	return category, nil
}

// PatchCategory patches the category
func (a *App) PatchCategory(cid int64, patch *model.CategoryPatch) (*model.Category, *model.AppErr) {
	old, err := a.Srv().Store.Category().Get(cid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	ucat, err := a.Srv().Store.Category().Update(cid, old)
	if err != nil {
		return nil, err
	}

	return ucat, nil
}

// GetCategory gets the category by the id
func (a *App) GetCategory(cid int64) (*model.Category, *model.AppErr) {
	return a.Srv().Store.Category().Get(cid)
}

// GetCategories gets all categories from the db
func (a *App) GetCategories() ([]*model.Category, *model.AppErr) {
	return a.Srv().Store.Category().GetAll()
}

// DeleteCategory hard deletes the category from the db
func (a *App) DeleteCategory(cid int64) *model.AppErr {
	return a.Srv().Store.Category().Delete(cid)
}
