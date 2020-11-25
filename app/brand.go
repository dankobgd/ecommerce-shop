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
	b.PreSave()
	if err := b.Validate(fh); err != nil {
		return nil, err
	}

	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("CreateBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandSizeExceeded, http.StatusInternalServerError, nil)
	}
	bb, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("CreateBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandFileErr, http.StatusInternalServerError, nil)
	}

	details, uErr := a.UploadImage(bytes.NewBuffer(bb), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}
	b.SetLogoDetails(details)

	brand, bErr := a.Srv().Store.Brand().Save(b)
	if bErr != nil {
		a.Log().Error(bErr.Error(), zlog.Err(bErr))
		return nil, bErr
	}

	return brand, nil
}

// PatchBrand patches the brand
func (a *App) PatchBrand(bid int64, patch *model.BrandPatch, fh *multipart.FileHeader) (*model.Brand, *model.AppErr) {
	if err := patch.Validate(fh); err != nil {
		return nil, err
	}

	old, err := a.Srv().Store.Brand().Get(bid)
	if err != nil {
		return nil, err
	}

	oldPublicID := old.LogoPublicID

	if fh != nil {
		thumbnail, err := fh.Open()
		if err != nil {
			return nil, model.NewAppErr("patchBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandSizeExceeded, http.StatusInternalServerError, nil)
		}
		b, err := ioutil.ReadAll(thumbnail)
		if err != nil {
			return nil, model.NewAppErr("patchBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandFileErr, http.StatusInternalServerError, nil)
		}

		details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return nil, uErr
		}

		old.SetLogoDetails(details)
	}

	old.Patch(patch)
	old.PreUpdate()
	ubrand, err := a.Srv().Store.Brand().Update(bid, old)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := a.DeleteImage(oldPublicID); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
		}
	}()

	return ubrand, nil
}

// GetBrand gets the brand by the id
func (a *App) GetBrand(bid int64) (*model.Brand, *model.AppErr) {
	return a.Srv().Store.Brand().Get(bid)
}

// GetBrands gets all brands from the db
func (a *App) GetBrands(limit, offset int) ([]*model.Brand, *model.AppErr) {
	return a.Srv().Store.Brand().GetAll(limit, offset)
}

// DeleteBrand hard deletes the brand from the db
func (a *App) DeleteBrand(bid int64) *model.AppErr {
	old, e := a.Srv().Store.Brand().Get(bid)
	if e != nil {
		return e
	}

	err := a.Srv().Store.Brand().Delete(bid)
	if err != nil {
		return err
	}

	defer func() {
		if err := a.DeleteImage(old.LogoPublicID); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
		}
	}()

	return nil
}
