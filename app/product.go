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
	msgProductSizeExceeded = &i18n.Message{ID: "app.product.create_product.image_size.app_error", Other: "upload image size exceeded"}
	msgProductFileErr      = &i18n.Message{ID: "app.product.create_product.formfile.app_error", Other: "error parsing files"}
	msgErrPropsJSONFile    = &i18n.Message{ID: "app.product.get_product_properties.app_error", Other: "error parsing properties json file"}
	msgProductImageFileErr = &i18n.Message{ID: "app.product.create_product_image.formfile.app_error", Other: "error parsing product image"}
	msgProductImagesErr    = &i18n.Message{ID: "app.product.create_product_images.formfile.app_error", Other: "No images provided"}
)

// GetProductsCount gets all products count
func (a *App) GetProductsCount() int {
	return a.Srv().Store.Product().Count()
}

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
		if oldPublicID != "" {
			if err := a.DeleteImage(oldPublicID); err != nil {
				a.Log().Error(err.Error(), zlog.Err(err))
			}
		}
	}()

	return uprod, nil
}

// DeleteProduct hard deletes the product from the db
func (a *App) DeleteProduct(pid int64) *model.AppErr {
	old, e := a.Srv().Store.Product().Get(pid)
	if e != nil {
		return e
	}

	err := a.Srv().Store.Product().Delete(pid)
	if err != nil {
		return err
	}

	defer func() {
		if old.ImageURL != "" {
			if err := a.DeleteImage(old.ImagePublicID); err != nil {
				a.Log().Error(err.Error(), zlog.Err(err))
			}
		}
	}()

	return nil
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

// GetMostSoldProducts returns most sold products
func (a *App) GetMostSoldProducts(limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetMostSold(limit, offset)
}

// GetBestDealsProducts returns the most discounted products
func (a *App) GetBestDealsProducts(limit, offset int) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().GetBestDeals(limit, offset)
}

// GetProductLatestPricing creates the discount
func (a *App) GetProductLatestPricing(pid int64) (*model.ProductPricing, *model.AppErr) {
	return a.Srv().Store.Product().GetLatestPricing(pid)
}

// DeleteProducts creates the discount
func (a *App) DeleteProducts(ids []int) *model.AppErr {
	return a.Srv().Store.Product().BulkDelete(ids)
}

// AddProductPricing adds the new pricing (updates the prev val and creates 2 new entries)
func (a *App) AddProductPricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr) {
	old, err := a.GetProductLatestPricing(pricing.ProductID)
	if err != nil {
		return nil, err
	}

	t := pricing.SaleStarts.Add(time.Millisecond - 1)

	patch := &model.ProductPricingPatch{SaleEnds: &t}
	old.Patch(patch)

	if _, err := a.UpdateProductPricing(old); err != nil {
		return nil, err
	}

	pricing.OriginalPrice = old.Price

	discount, err := a.InsertProductPricing(pricing)
	if err != nil {
		return nil, err
	}

	pricingAfter := &model.ProductPricing{
		ProductID:     pricing.ProductID,
		Price:         old.Price,
		OriginalPrice: discount.Price,
		SaleStarts:    pricing.SaleEnds.Add(time.Millisecond),
		SaleEnds:      model.FutureSaleEndsTime,
	}

	if _, err := a.InsertProductPricing(pricingAfter); err != nil {
		return nil, err
	}
	return pricing, nil
}

// InsertProductPricing creates the discount
func (a *App) InsertProductPricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr) {
	if err := pricing.Validate(); err != nil {
		return nil, err
	}
	return a.Srv().Store.Product().InsertPricing(pricing)
}

// UpdateProductPricing updates the discount
func (a *App) UpdateProductPricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr) {
	return a.Srv().Store.Product().UpdatePricing(pricing)
}

// CreateProductTag gets all tags for the product
func (a *App) CreateProductTag(pid int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().Save(pid, pt)
}

// GetProductTags gets all tags for the product
func (a *App) GetProductTags(pid int64) ([]*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().GetAll(pid)
}

// ReplaceProductTags patches the product tag
func (a *App) ReplaceProductTags(pid int64, tagIDs []int) ([]*model.ProductTag, *model.AppErr) {
	return a.Srv().Store.ProductTag().Replace(pid, tagIDs)
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

// DeleteProductTags bulk deletes tags
func (a *App) DeleteProductTags(pid int64, ids []int) *model.AppErr {
	return a.Srv().Store.ProductTag().BulkDelete(pid, ids)
}

// CreateProductImages bulk inserts product images
func (a *App) CreateProductImages(pid int64, imagesFHs []*multipart.FileHeader) *model.AppErr {
	if len(imagesFHs) == 0 {
		return model.NewAppErr("CreateProductImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductImagesErr, http.StatusInternalServerError, nil)
	}

	images := make([]*model.ProductImage, 0)

	tmp := &model.ProductImage{}
	for _, fh := range imagesFHs {
		if err := tmp.Validate(fh); err != nil {
			return err
		}
	}

	for _, fh := range imagesFHs {
		f, err := fh.Open()
		if err != nil {
			return model.NewAppErr("CreateProductImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductImageFileErr, http.StatusInternalServerError, nil)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return model.NewAppErr("CreateProductImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductImageFileErr, http.StatusInternalServerError, nil)
		}
		// TODO: upload in parallel...
		details, uErr := a.UploadImage(bytes.NewBuffer(b), fh.Filename)
		if uErr != nil {
			return uErr
		}

		images = append(images, &model.ProductImage{
			ProductID: model.NewInt64(pid),
			URL:       model.NewString(details.SecureURL),
			PublicID:  model.NewString(details.PublicID),
		})
	}

	for _, img := range images {
		img.PreSave()
	}

	if err := a.Srv().Store.ProductImage().BulkInsert(images); err != nil {
		a.Log().Error(err.Error(), zlog.Err(err))
		return err
	}

	return nil
}

// CreateProductImage creates the image for product
func (a *App) CreateProductImage(pid int64, img *model.ProductImage, fh *multipart.FileHeader) (*model.ProductImage, *model.AppErr) {
	if err := img.Validate(fh); err != nil {
		return nil, err
	}

	thumbnail, err := fh.Open()
	if err != nil {
		return nil, model.NewAppErr("CreateProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductImageFileErr, http.StatusInternalServerError, nil)
	}
	b, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		return nil, model.NewAppErr("CreateProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductImageFileErr, http.StatusInternalServerError, nil)
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

// CreateProductReview creates new review for the product
func (a *App) CreateProductReview(pid int64, rev *model.ProductReview) (*model.ProductReview, *model.AppErr) {
	rev.PreSave()
	if err := rev.Validate(); err != nil {
		return nil, err
	}

	return a.Srv().Store.ProductReview().Save(pid, rev)
}

// GetProductReviews gets all reviews for the product
func (a *App) GetProductReviews(pid int64) ([]*model.ProductReview, *model.AppErr) {
	return a.Srv().Store.ProductReview().GetAll(pid)
}

// GetProductReview gets all reviews for the product
func (a *App) GetProductReview(pid, rid int64) (*model.ProductReview, *model.AppErr) {
	return a.Srv().Store.ProductReview().Get(pid, rid)
}

// PatchProductReview patches the product review
func (a *App) PatchProductReview(pid, rid int64, patch *model.ProductReviewPatch) (*model.ProductReview, *model.AppErr) {
	if err := patch.Validate(); err != nil {
		return nil, err
	}

	old, err := a.Srv().Store.ProductReview().Get(pid, rid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	urev, err := a.Srv().Store.ProductReview().Update(pid, rid, old)
	if err != nil {
		return nil, err
	}

	return urev, nil
}

// DeleteProductReview deletes the product review
func (a *App) DeleteProductReview(pid, rid int64) *model.AppErr {
	return a.Srv().Store.ProductReview().Delete(pid, rid)
}

// DeleteProductReviews bulk deletes reviews
func (a *App) DeleteProductReviews(pid int64, ids []int) *model.AppErr {
	return a.Srv().Store.ProductReview().BulkDelete(pid, ids)
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
		if *oldPublicID != "" {
			if err := a.DeleteImage(*oldPublicID); err != nil {
				a.Log().Error(err.Error(), zlog.Err(err))
			}
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
		if *old.PublicID != "" {
			if err := a.DeleteImage(*old.PublicID); err != nil {
				a.Log().Error(err.Error(), zlog.Err(err))
			}
		}
	}()

	return nil
}

// DeleteProductImages bulk deletes images
func (a *App) DeleteProductImages(pid int64, ids []int) *model.AppErr {
	return a.Srv().Store.ProductImage().BulkDelete(pid, ids)
}

// SearchProducts performs the full text search on products
func (a *App) SearchProducts(query string) ([]*model.Product, *model.AppErr) {
	return a.Srv().Store.Product().Search(query)
}
