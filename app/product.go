package app

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

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
func (a *App) CreateProduct(p *model.Product, thumbnailFH *multipart.FileHeader, imagesFHs []*multipart.FileHeader, tagids []string) (*model.Product, *model.AppErr) {
	p.PreSave()
	if err := p.Validate(thumbnailFH); err != nil {
		return nil, err
	}

	thumbnail, err := thumbnailFH.Open()
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}

	details, uErr := a.UploadImage(bytes.NewBuffer(b), thumbnailFH.Filename)
	if uErr != nil {
		return nil, uErr
	}
	p.SetImageDetails(details)

	product, pErr := a.Srv().Store.Product().Save(p)
	if pErr != nil {
		a.Log().Error(pErr.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	// handle optional create product tags
	if len(tagids) > 0 {
		tags := make([]*model.ProductTag, 0)
		for _, tid := range tagids {
			id, err := strconv.ParseInt(tid, 10, 64)
			if err != nil {
				return nil, model.NewAppErr("CreateProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
			}
			tags = append(tags, &model.ProductTag{
				TagID:     model.NewInt64(id),
				ProductID: model.NewInt64(product.ID),
			})
		}

		if err := a.Srv().Store.ProductTag().BulkInsert(tags); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
			return nil, err
		}
	}

	// handle optional create product images
	if len(imagesFHs) > 0 {
		images := make([]*model.ProductImage, 0)

		for _, fh := range imagesFHs {
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
			images = append(images, &model.ProductImage{
				ProductID: model.NewInt64(product.ID),
				URL:       model.NewString(details.SecureURL),
				PublicID:  model.NewString(details.PublicID),
			})
		}

		for _, img := range images {
			img.PreSave()
		}
		if err := a.Srv().Store.ProductImage().BulkInsert(images); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
			return nil, err
		}
	}

	return product, nil
}

// PatchProduct patches the product
func (a *App) PatchProduct(pid int64, patch *model.ProductPatch, fh *multipart.FileHeader) (*model.Product, *model.AppErr) {
	if err := patch.Validate(fh); err != nil {
		return nil, err
	}

	old, err := a.Srv().Store.Product().Get(pid)
	if err != nil {
		return nil, err
	}

	oldPublicID := old.ImagePublicID

	if fh != nil {
		thumbnail, err := fh.Open()
		if err != nil {
			return nil, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
		}
		b, err := ioutil.ReadAll(thumbnail)
		if err != nil {
			return nil, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}

		details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return nil, uErr
		}
		old.SetImageDetails(details)
	}

	old.Patch(patch)
	old.PreUpdate()

	uprod, err := a.Srv().Store.Product().Update(pid, old)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := a.DeleteImage(oldPublicID); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
		}
	}()

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
func (a *App) GetProducts(filters map[string][]string, limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetAll(filters, limit, offset)
}

// GetFeaturedProducts returns featured products
func (a *App) GetFeaturedProducts(limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetFeatured(limit, offset)
}

// CreateProductTag gets all tags for the product
func (a *App) CreateProductTag(pid int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().Save(pid, pt)
}

// GetProductTags gets all tags for the product
func (a *App) GetProductTags(pid int64) ([]*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().GetAll(pid)
}

// PatchProductTag patches the product tag
func (a *App) PatchProductTag(pid, tid int64, patch *model.ProductTagPatch) (*model.ProductTag, *model.AppErr) {
	old, err := a.Srv().Store.ProductTag().Get(pid, tid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	utag, err := a.Srv().Store.ProductTag().Update(pid, tid, old)
	if err != nil {
		return nil, err
	}

	return utag, nil
}

// DeleteProductTag gets all tags for the product
func (a *App) DeleteProductTag(pid, tid int64) *model.AppErr {
	return a.Srv().Store.ProductTag().Delete(pid, tid)
}

// CreateProductImage creates the image for product
func (a *App) CreateProductImage(pid int64, img *model.ProductImage, fh *multipart.FileHeader) (*model.ProductImage, *model.AppErr) {
	if err := img.Validate(fh); err != nil {
		return nil, err
	}

	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("CreateProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("CreateProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
	}

	img.PreSave()

	details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
	if uErr != nil {
		return nil, uErr
	}
	img.SetImageDetails(details)

	return a.Srv().Store.ProductImage().Save(pid, img)
}

// GetProductImages gets all images for the product
func (a *App) GetProductImages(pid int64) ([]*model.ProductImage, *model.AppErr) {
	return a.Srv().Store.ProductImage().GetAll(pid)
}

// GetProductReviews gets all reviews for the product
func (a *App) GetProductReviews(pid int64) ([]*model.Review, *model.AppErr) {
	return a.Srv().Store.Product().GetReviews(pid)
}

// PatchProductImage patches the product image
func (a *App) PatchProductImage(pid, imgID int64, patch *model.ProductImagePatch, fh *multipart.FileHeader) (*model.ProductImage, *model.AppErr) {
	if err := patch.Validate(fh); err != nil {
		return nil, err
	}

	old, err := a.Srv().Store.ProductImage().Get(pid, imgID)
	if err != nil {
		return nil, err
	}

	oldPublicID := old.PublicID

	if fh != nil {
		thumbnail, err := fh.Open()
		if err != nil {
			return nil, model.NewAppErr("PatchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductSizeExceeded, http.StatusInternalServerError, nil)
		}
		b, err := ioutil.ReadAll(thumbnail)
		if err != nil {
			return nil, model.NewAppErr("PatchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil)
		}

		details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return nil, uErr
		}

		old.SetImageDetails(details)
	}

	old.Patch(patch)
	old.PreUpdate()

	uimg, err := a.Srv().Store.ProductImage().Update(pid, imgID, old)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := a.DeleteImage(*oldPublicID); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
		}
	}()

	return uimg, nil
}

// GetProductsbyIDS gets products with specified ids slice
func (a *App) GetProductsbyIDS(ids []int64) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().ListByIDS(ids)
}

// DeleteProductImage deletes the product image
func (a *App) DeleteProductImage(pid, imgID int64) *model.AppErr {
	old, e := a.Srv().Store.ProductImage().Get(pid, imgID)
	if e != nil {
		return e
	}

	err := a.Srv().Store.ProductImage().Delete(pid, imgID)
	if err != nil {
		return err
	}

	defer func() {
		if err := a.DeleteImage(*old.PublicID); err != nil {
			a.Log().Error(err.Error(), zlog.Err(err))
		}
	}()

	return nil
}

// SearchProducts performs the full text search on products
func (a *App) SearchProducts(query string) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().Search(query)
}
