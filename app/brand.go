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
	msgBrandSizeExceeded = &i18n.Message{ID: "app.brand.create_brand.image_size.app_error", Other: "upload image size exceeded"}
	msgBrandFileErr      = &i18n.Message{ID: "app.brand.create_brand.formfile.app_error", Other: "error parsing files"}
)

// CreateBrand creates the new brand in the system
func (a *App) CreateBrand(b *model.Brand, fh *multipart.FileHeader) (*model.Brand, *model.AppErr) {
	if fh.Size > model.FileUploadSizeLimit {
		return nil, model.NewAppErr("CreateBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandSizeExceeded, http.StatusInternalServerError, nil)
	}
	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("CreateBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandSizeExceeded, http.StatusInternalServerError, nil)
	}
	bb, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("CreateBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandFileErr, http.StatusInternalServerError, nil)
	}

	b.PreSave()
	if err := b.Validate(); err != nil {
		return nil, err
	}

	details, uErr := a.UploadImage(bytes.NewBuffer(bb), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}
	b.SetLogoURL(details.SecureURL)

	brand, bErr := a.Srv().Store.Brand().Save(b)
	if bErr != nil {
		a.Log().Error(bErr.Error(), zlog.Err(bErr))
		return nil, bErr
	}

	return brand, nil
}

// PatchBrand patches the brand
func (a *App) PatchBrand(bid int64, patch *model.BrandPatch) (*model.Brand, *model.AppErr) {
	old, err := a.Srv().Store.Brand().Get(bid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	ubrand, err := a.Srv().Store.Brand().Update(bid, old)
	if err != nil {
		return nil, err
	}

	return ubrand, nil
}

// GetBrand gets the brand by the id
func (a *App) GetBrand(bid int64) (*model.Brand, *model.AppErr) {
	return a.Srv().Store.Brand().Get(bid)
}

// GetBrands gets all brands from the db
func (a *App) GetBrands() ([]*model.Brand, *model.AppErr) {
	return a.Srv().Store.Brand().GetAll()
}

// DeleteBrand hard deletes the brand from the db
func (a *App) DeleteBrand(bid int64) *model.AppErr {
	return a.Srv().Store.Brand().Delete(bid)
}
