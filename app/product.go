package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductSizeExceeded = &i18n.Message{ID: "app.product.create_product.image_size.app_error", Other: "upload image size exceeded"}
	msgProductFileErr      = &i18n.Message{ID: "app.product.create_product.formfile.app_error", Other: "error parsing files"}
	msgErrPropsJSONFile    = &i18n.Message{ID: "app.product.get_product_properties.app_error", Other: "error parsing properties json file"}
)

// CreateProduct creates the new product in the system
func (a *App) CreateProduct(p *model.Product, fh *multipart.FileHeader, headers []*multipart.FileHeader, tags []*model.ProductTag, properties string) (*model.Product, *model.AppErr) {
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

	details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}

	p.SetImageURL(details.SecureURL)
	p.SetProperties(properties)

	product, pErr := a.Srv().Store.Product().Save(p)
	if pErr != nil {
		a.Log().Error(pErr.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	for _, t := range tags {
		t.ProductID = model.NewInt64(product.ID)
	}

	images := make([]*model.ProductImage, 0)

	for _, fh := range headers {
		f, err := fh.Open()
		if err != nil {
			return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}
		// TODO: upload in parallel...
		details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return nil, uErr
		}
		img := &model.ProductImage{ProductID: model.NewInt64(product.ID), URL: model.NewString(details.SecureURL)}
		images = append(images, img)
	}

	if len(tags) > 0 {
		tagids, err := a.Srv().Store.ProductTag().BulkInsert(tags)
		if err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
			return nil, err
		}
		for i, id := range tagids {
			tags[i].ID = model.NewInt64(id)
		}
	} else {
		tags = make([]*model.ProductTag, 0)
	}

	if len(images) > 0 {
		for _, img := range images {
			img.PreSave()
		}

		imgids, err := a.Srv().Store.ProductImage().BulkInsert(images)
		if err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
			return nil, err
		}
		for i, id := range imgids {
			images[i].ID = model.NewInt64(id)
		}
	} else {
		images = make([]*model.ProductImage, 0)
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
func (a *App) GetProducts(limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetAll(limit, offset)
}

// GetFeaturedProducts returns featured products
func (a *App) GetFeaturedProducts(limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetFeatured(limit, offset)
}

// GetProductTags gets all tags for the product
func (a *App) GetProductTags(pid int64) ([]*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().GetAll(pid)
}

// GetProductImages gets all images for the product
func (a *App) GetProductImages(pid int64) ([]*model.ProductImage, *model.AppErr) {
	return a.Srv().Store.ProductImage().GetAll(pid)
}

// GetProductReviews gets all reviews for the product
func (a *App) GetProductReviews(pid int64) ([]*model.Review, *model.AppErr) {
	return a.Srv().Store.Product().GetReviews(pid)
}

// GetProductTag gets the product tag by id
func (a *App) GetProductTag(id int64) (*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().Get(id)
}

// GetImage gets the image by id
func (a *App) GetImage(id int64) (*model.ProductImage, *model.AppErr) {
	return a.Srv().Store.ProductImage().Get(id)
}

// PatchProductTag patches the product tag
func (a *App) PatchProductTag(tid int64, patch *model.ProductTagPatch) (*model.ProductTag, *model.AppErr) {
	old, err := a.Srv().Store.ProductTag().Get(tid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	utag, err := a.Srv().Store.ProductTag().Update(tid, old)
	if err != nil {
		return nil, err
	}

	return utag, nil
}

// PatchProductImage patches the product image
func (a *App) PatchProductImage(imgID int64, patch *model.ProductImagePatch) (*model.ProductImage, *model.AppErr) {
	old, err := a.Srv().Store.ProductImage().Get(imgID)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	uimg, err := a.Srv().Store.ProductImage().Update(imgID, old)
	if err != nil {
		return nil, err
	}

	return uimg, nil
}

// GetProductsbyIDS gets products with specified ids slice
func (a *App) GetProductsbyIDS(ids []int64) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().ListByIDS(ids)
}

// DeleteProductTag deletes the product tag
func (a *App) DeleteProductTag(tid int64) *model.AppErr {
	// TODO: delete from cloud later
	return a.Srv().Store.ProductTag().Delete(tid)
}

// DeleteProductImage deletes the product image
func (a *App) DeleteProductImage(imgID int64) *model.AppErr {
	// TODO: delete from cloud later
	return a.Srv().Store.ProductImage().Delete(imgID)
}

// GetProductProperties gets the valid products properties (variants for each specific category - size, colors etc...)
func (a *App) GetProductProperties() (*model.ProductProperties, *model.AppErr) {
	file, err := ioutil.ReadFile("./data/variants/variants.json")
	if err != nil {
		return nil, model.NewAppErr("GetProductProperties", model.ErrInternal, locale.GetUserLocalizer("en"), msgErrPropsJSONFile, http.StatusInternalServerError, nil)
	}

	props := &model.ProductProperties{}

	if err := json.Unmarshal([]byte(file), &props); err != nil {
		return nil, model.NewAppErr("GetProductProperties", model.ErrInternal, locale.GetUserLocalizer("en"), msgErrPropsJSONFile, http.StatusInternalServerError, nil)
	}

	return props, nil
}
