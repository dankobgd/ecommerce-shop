package app

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/fileupload"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductSizeExceeded = &i18n.Message{ID: "app.product.create_product.image_size.app_error", Other: "upload image size exceeded"}
	msgProductFileErr      = &i18n.Message{ID: "app.product.create_product.formfile.app_error", Other: "error parsing files"}
)

// CreateProduct creates the new product in the system
func (a *App) CreateProduct(p *model.Product, pTag *model.ProductTag, pCat *model.ProductCategory, pBrand *model.ProductBrand, file multipart.File, fh *multipart.FileHeader) (*model.Product, *model.AppErr) {
	if fh.Size > model.FileUploadSizeLimit {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}

	p.PreSave()
	if err := p.Validate(); err != nil {
		return nil, err
	}
	pTag.PreSave()
	if err := pTag.Validate(); err != nil {
		return nil, err
	}
	if err := pCat.Validate(); err != nil {
		return nil, err
	}
	pBrand.PreSave()
	if err := pBrand.Validate(); err != nil {
		return nil, err
	}

	url, uErr := a.UploadProductImage(fileBytes, fh)
	if uErr != nil {
		return nil, uErr
	}
	p.SetImageURL(url)

	product, pErr := a.Srv().Store.Product().Save(p)
	if pErr != nil {
		a.log.Error(err.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	pTag.SetProductID(product.ID)
	if err := a.Srv().Store.Product().InsertTag(pTag); err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}
	pCat.SetProductID(product.ID)
	if err := a.Srv().Store.Product().InsertCategory(pCat); err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}
	pBrand.SetProductID(product.ID)
	if err := a.Srv().Store.Product().InsertBrand(pBrand); err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}

	return product, nil
}

func (a *App) uploadImageToCloudinary(fileBytes []byte, fh *multipart.FileHeader) (string, *model.AppErr) {
	return fileupload.UploadImageToCloudinary(fileBytes, fh, a.Cfg().CloudinarySettings.EnvURI)
}

// UploadProductImage uploads the image and returns the preview url
func (a *App) UploadProductImage(fileBytes []byte, fh *multipart.FileHeader) (string, *model.AppErr) {
	return a.uploadImageToCloudinary(fileBytes, fh)
}
