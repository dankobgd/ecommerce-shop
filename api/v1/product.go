package apiv1

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/pagination"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductPatchFromJSON   = &i18n.Message{ID: "api.product.patch_product.app_error", Other: "could not decode product patch data"}
	msgProductFileErr         = &i18n.Message{ID: "api.product.create_product.formfile.app_error", Other: "error parsing files"}
	msgProductAvatarMultipart = &i18n.Message{ID: "api.product.create_product.multipart.app_error", Other: "could not decode product multipart data"}
	msgProductPriceErr        = &i18n.Message{ID: "api.product.create_product.price.app_error", Other: "could not decode product price"}
	msgPatchProduct           = &i18n.Message{ID: "api.product.patch_product.app_error", Other: "could not patch product"}
	msgURLParamErr            = &i18n.Message{ID: "api.product.url.params.app_error", Other: "could not parse URL params"}
	msgDiscountFromJSON       = &i18n.Message{ID: "api.product.create_product_discount.app_error", Other: "could not parse discount pricing from json"}
	msgReviewFromJSON         = &i18n.Message{ID: "api.product.create_product_review.app_error", Other: "could not parse product review from json"}
	msgReviewURLParamErr      = &i18n.Message{ID: "api.product.create_product_review.app_error", Other: "invalid product review url param"}
	msgReviewPatchFromJSONErr = &i18n.Message{ID: "api.product.patch_product_review.app_error", Other: "could not decode product review patch data"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.Routes.Products.Post("/", a.AdminSessionRequired(a.createProduct))
	a.Routes.Products.Get("/", a.getProducts)
	a.Routes.Products.Get("/featured", a.getFeaturedProducts)
	a.Routes.Products.Get("/search", a.searchProducts)
	a.Routes.Products.Delete("/bulk", a.deleteProducts)

	a.Routes.Product.Get("/", a.getProduct)
	a.Routes.Product.Patch("/", a.AdminSessionRequired(a.patchProduct))
	a.Routes.Product.Delete("/", a.AdminSessionRequired(a.deleteProduct))

	// discount
	a.Routes.Product.Get("/pricing/latest", a.AdminSessionRequired(a.getProductLatestPricing))
	a.Routes.Product.Post("/pricing", a.AdminSessionRequired(a.createProductPricing))

	// product tags
	a.Routes.Product.Post("/tags", a.AdminSessionRequired(a.createProductTag))
	a.Routes.Product.Get("/tags", a.getProductTags)
	a.Routes.Product.Put("/tags/replace", a.AdminSessionRequired(a.replaceProductTags))
	a.Routes.Product.Patch("/tags/{tag_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.patchProductTag))
	a.Routes.Product.Delete("/tags/{tag_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.deleteProductTag))
	a.Routes.Product.Delete("/tags/bulk", a.AdminSessionRequired(a.deleteProductTags))

	// product images
	a.Routes.Product.Post("/images/bulk", a.AdminSessionRequired(a.createProductImages))
	a.Routes.Product.Post("/images", a.AdminSessionRequired(a.createProductImage))
	a.Routes.Product.Get("/images", a.getProductImages)
	a.Routes.Product.Patch("/images/{image_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.patchProductImage))
	a.Routes.Product.Delete("/images/{image_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.deleteProductImage))
	a.Routes.Product.Delete("/images/bulk", a.AdminSessionRequired(a.deleteProductImages))

	// product reviews
	a.Routes.Product.Post("/reviews", a.SessionRequired(a.createProductReview))
	a.Routes.Product.Get("/reviews", a.getProductReviews)
	a.Routes.Product.Get("/reviews/{review_id:[A-Za-z0-9]+}", a.getProductReview)
	a.Routes.Product.Patch("/reviews/{review_id:[A-Za-z0-9]+}", a.SessionRequired(a.patchProductReview))
	a.Routes.Product.Delete("/reviews/{review_id:[A-Za-z0-9]+}", a.SessionRequired(a.deleteProductReview))
	a.Routes.Product.Delete("/reviews/bulk", a.AdminSessionRequired(a.deleteProductReviews))
}

func (a *API) createProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	p := &model.Product{}
	if err := model.SchemaDecoder.Decode(p, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	// or just use gorilla/schema embeded struct with: `schema:"pricing"` and then send formData: Pricing.price
	price := mpf.Value["price"]
	tagids := mpf.Value["tags"]
	images := mpf.File["images"]

	if len(price) == 0 {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductPriceErr, http.StatusInternalServerError, nil))
		return
	}
	priceValue, err := strconv.Atoi(price[0])
	if err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductPriceErr, http.StatusInternalServerError, nil))
		return
	}
	p.ProductPricing = &model.ProductPricing{
		Price:         priceValue,
		OriginalPrice: priceValue,
	}

	var thumbnail *multipart.FileHeader
	if len(mpf.File["image"]) > 0 {
		thumbnail = mpf.File["image"][0]
	}

	product, pErr := a.app.CreateProduct(p, thumbnail, images, tagids)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

func (a *API) patchProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	patch := &model.ProductPatch{}
	if err := model.SchemaDecoder.Decode(patch, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}
	patch.SetProperties(patch.PropertiesText)

	var image *multipart.FileHeader
	if len(mpf.File["image"]) > 0 {
		image = mpf.File["image"][0]
	}

	pid, err := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	uprod, pErr := a.app.PatchProduct(pid, patch, image)
	if err != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, uprod)
}

func (a *API) deleteProduct(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteProduct(pid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) getProduct(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	p, err := a.app.GetProduct(pid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, p)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()
	pages := pagination.NewFromRequest(r)
	products, err := a.app.GetProducts(filters, pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(products) > 0 {
		totalCount = products[0].TotalCount
	}
	pages.SetData(products, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) getProductLatestPricing(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductLatestPricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	pp, err := a.app.GetProductLatestPricing(pid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, pp)
}

func (a *API) createProductPricing(w http.ResponseWriter, r *http.Request) {
	salePricing, e := model.ProductPricingFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createProductPricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgDiscountFromJSON, http.StatusInternalServerError, nil))
		return
	}

	old, err := a.app.GetProductLatestPricing(salePricing.ProductID)
	if err != nil {
		respondError(w, err)
		return
	}

	t := salePricing.SaleStarts.Add(time.Millisecond - 1)
	patch := &model.ProductPricingPatch{SaleEnds: &t}
	old.Patch(patch)

	if _, err := a.app.UpdateProductPricing(old); err != nil {
		respondError(w, err)
		return
	}

	salePricing.OriginalPrice = old.Price

	discount, err := a.app.CreateProductPricing(salePricing)
	if err != nil {
		respondError(w, err)
		return
	}

	pricingAfter := &model.ProductPricing{
		ProductID:     salePricing.ProductID,
		Price:         old.Price,
		OriginalPrice: discount.Price,
		SaleStarts:    salePricing.SaleEnds.Add(time.Millisecond),
		SaleEnds:      model.FutureSaleEndsTime,
	}

	if _, err := a.app.CreateProductPricing(pricingAfter); err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, discount)
}

func (a *API) createProductTag(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("createProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	pt, e := model.ProductTagFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagFromJSON, http.StatusInternalServerError, nil))
		return
	}

	productTag, err := a.app.CreateProductTag(pid, pt)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, productTag)
}

func (a *API) getProductTags(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	tags, err := a.app.GetProductTags(pid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, tags)
}

func (a *API) replaceProductTags(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("replaceProductTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	tagIDs := model.IntSliceFromJSON(r.Body)

	newTags, pErr := a.app.ReplaceProductTags(pid, tagIDs)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, newTags)
}

func (a *API) patchProductTag(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("patchProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	tid, err := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.ProductTagPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductPatchFromJSON, http.StatusInternalServerError, nil))
		return
	}

	utag, pErr := a.app.PatchProductTag(pid, tid, patch)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, utag)
}

func (a *API) deleteProductTag(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	tid, e := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := a.app.DeleteProductTag(pid, tid); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}

func (a *API) deleteProductTags(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteProductTags(pid, ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) createProductImages(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("createProductImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createProductImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	images := mpf.File["images"]

	if pErr := a.app.CreateProductImages(pid, images); pErr != nil {
		respondError(w, pErr)
		return
	}

	respondOK(w)
}

func (a *API) createProductImage(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("createProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	img := &model.ProductImage{}
	if err := model.SchemaDecoder.Decode(img, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	image := mpf.File["image"][0]

	productImage, pErr := a.app.CreateProductImage(pid, img, image)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, productImage)
}

func (a *API) getProductImages(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	imgs, err := a.app.GetProductImages(pid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, imgs)
}

func (a *API) createProductReview(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("createProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	rev, e := model.ReviewFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewFromJSON, http.StatusInternalServerError, nil))
		return
	}

	rev.UserID = uid
	rev.ProductID = pid

	review, err := a.app.CreateProductReview(pid, rev)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, review)
}

func (a *API) getProductReviews(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductReviews", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	reviews, err := a.app.GetProductReviews(pid)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, reviews)
}

func (a *API) getProductReview(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	rid, e := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	review, err := a.app.GetProductReview(pid, rid)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, review)
}

func (a *API) patchProductReview(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("patchProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	rid, err := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.ReviewPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewPatchFromJSONErr, http.StatusInternalServerError, nil))
		return
	}

	urev, rErr := a.app.PatchProductReview(pid, rid, patch)
	if err != nil {
		respondError(w, rErr)
		return
	}

	respondJSON(w, http.StatusOK, urev)
}

func (a *API) deleteProductReview(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	rid, err := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteProductReview(pid, rid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) deleteProductReviews(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteProductReviews(pid, ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) patchProductImage(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getProductReviews", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	imgID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("patchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)
	image := mpf.File["image"][0]

	patch := &model.ProductImagePatch{}
	if err := model.SchemaDecoder.Decode(patch, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("patchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductAvatarMultipart, http.StatusInternalServerError, nil))
		return
	}

	uimg, pErr := a.app.PatchProductImage(pid, imgID, patch, image)
	if err != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, uimg)
}

func (a *API) deleteProductImage(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	imgID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := a.app.DeleteProductImage(pid, imgID); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}

func (a *API) deleteProductImages(w http.ResponseWriter, r *http.Request) {
	pid, e := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("deleteProductReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteProductImages(pid, ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) getFeaturedProducts(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	featured, err := a.app.GetFeaturedProducts(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(featured) > 0 {
		totalCount = featured[0].TotalCount
	}
	pages.SetData(featured, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	searchResults, err := a.app.SearchProducts(query)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, searchResults)
}

func (a *API) deleteProducts(w http.ResponseWriter, r *http.Request) {
	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteProducts(ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
