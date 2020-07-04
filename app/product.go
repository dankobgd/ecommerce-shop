package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

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
func (a *App) CreateProduct(data *model.ProductCreateData) (*model.Product, *model.AppErr) {
	if data.ImgFH.Size > model.FileUploadSizeLimit {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}
	thumbnail, err := data.ImgFH.Open()
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}

	data.P.PreSave()
	if err := data.P.Validate(); err != nil {
		return nil, err
	}
	data.Tag.PreSave()
	if err := data.Tag.Validate(); err != nil {
		return nil, err
	}
	data.Brand.PreSave()
	if err := data.Brand.Validate(); err != nil {
		return nil, err
	}
	if err := data.Cat.Validate(); err != nil {
		return nil, err
	}

	url, uErr := a.UploadProductImage(bytes.NewBuffer(b), data.ImgFH.Filename)
	if uErr != nil {
		return nil, uErr
	}
	data.P.SetImageURL(url)

	product, pErr := a.Srv().Store.Product().Save(data.P, data.Brand, data.Cat)
	if pErr != nil {
		a.log.Error(err.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	var productTags = make([]*model.ProductTag, 0)
	var productImages = make([]*model.ProductImage, 0)

	for _, tn := range data.TagNames {
		tag := &model.ProductTag{
			ProductID: product.ID,
			Name:      tn,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		productTags = append(productTags, tag)
	}

	for _, fh := range data.ImageHeaders {
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
		productImages = append(productImages, img)
	}

	if len(productTags) > 0 {
		if err := a.Srv().Store.Product().BulkInsertTags(productTags); err != nil {
			a.log.Error(err.Error(), zlog.Err(err))
			return nil, err
		}
	}

	if len(productImages) > 0 {
		if err := a.Srv().Store.Product().BulkInsertImages(productImages); err != nil {
			a.log.Error(err.Error(), zlog.Err(err))
			return nil, err
		}
	}

	return product, nil
}

func (a *App) uploadImageToCloudinary(data io.Reader, filename string) (string, *model.AppErr) {
	return fileupload.UploadImageToCloudinary(data, filename, a.Cfg().CloudinarySettings.EnvURI)
}

// UploadProductImage uploads the image and returns the preview url
func (a *App) UploadProductImage(data io.Reader, filename string) (string, *model.AppErr) {
	return a.uploadImageToCloudinary(data, filename)
}
