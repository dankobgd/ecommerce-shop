package app

import (
	"bytes"
	"io"
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
func (a *App) CreateProduct(p *model.Product, fh *multipart.FileHeader, headers []*multipart.FileHeader) (*model.Product, *model.AppErr) {
	if fh.Size > model.FileUploadSizeLimit {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}
	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}

	p.PreSave()
	if err := p.Validate(); err != nil {
		return nil, err
	}

	url, uErr := a.UploadProductImage(bytes.NewBuffer(b), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}
	p.SetImageURL(url)

	product, pErr := a.Srv().Store.Product().Save(p)
	if pErr != nil {
		a.log.Error(pErr.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	p.Category.ProductID = product.ID
	p.Brand.ProductID = product.ID
	for _, t := range p.Tags {
		t.ProductID = product.ID
	}

	for _, fh := range headers {
		f, err := fh.Open()
		defer f.Close()
		if err != nil {
			return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}
		// TODO: upload in parallel...
		url, uErr := a.UploadProductImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return nil, uErr
		}
		img := &model.ProductImage{ProductID: product.ID, URL: url}
		p.Images = append(p.Images, img)
	}

	if len(p.Tags) > 0 {
		for _, t := range p.Tags {
			t.PreSave()
		}

		tagids, err := a.Srv().Store.ProductTag().BulkInsert(p.Tags)
		if err != nil {
			a.log.Error(err.Error(), zlog.Err(err))
			return nil, err
		}
		for i, id := range tagids {
			p.Tags[i].ID = id
		}
	}

	if len(p.Images) > 0 {
		imgids, err := a.Srv().Store.ProductImage().BulkInsert(p.Images)
		if err != nil {
			a.log.Error(err.Error(), zlog.Err(err))
			return nil, err
		}
		for i, id := range imgids {
			p.Images[i].ID = id
		}
	}

	return product, nil
}

// PatchProduct patches the product
func (a *App) PatchProduct(pid int64, patch *model.ProductPatch) (*model.Product, *model.AppErr) {
	old, err := a.Srv().Store.Product().Get(pid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	uprod, err := a.Srv().Store.Product().Update(pid, old)
	if err != nil {
		return nil, err
	}

	return uprod, nil
}

// DeleteProduct hard deletes the product from the db
func (a *App) DeleteProduct(pid int64) *model.AppErr {
	return a.Srv().Store.Product().Delete(pid)
}

// GetProduct gets the product by the id
func (a *App) GetProduct(pid int64) (*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().Get(pid)
}

// GetProducts gets all products from the db
func (a *App) GetProducts() ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetAll()
}

func (a *App) uploadImageToCloudinary(data io.Reader, filename string) (string, *model.AppErr) {
	return fileupload.UploadImageToCloudinary(data, filename, a.Cfg().CloudinarySettings.EnvURI)
}

// UploadProductImage uploads the image and returns the preview url
func (a *App) UploadProductImage(data io.Reader, filename string) (string, *model.AppErr) {
	return a.uploadImageToCloudinary(data, filename)
}
