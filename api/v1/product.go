package apiv1

import (
	"mime/multipart"
	"net/http"
	"strconv"

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
	msgPatchProduct           = &i18n.Message{ID: "api.product.patch_product.app_error", Other: "could not patch product"}
	msgURLParamErr            = &i18n.Message{ID: "api.product.url.params.app_error", Other: "could not parse URL params"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.Routes.Products.Get("/", a.getProducts)
	a.Routes.Products.Post("/", a.AdminSessionRequired(a.createProduct))
	a.Routes.Products.Get("/featured", a.getFeaturedProducts)
	a.Routes.Products.Get("/search", a.searchProducts)

	a.Routes.Product.Get("/", a.getProduct)
	a.Routes.Product.Patch("/", a.AdminSessionRequired(a.patchProduct))
	a.Routes.Product.Delete("/", a.AdminSessionRequired(a.deleteProduct))

	// product reviews
	a.Routes.Product.Get("/reviews", a.getProductReviews)

	// product tags
	a.Routes.Product.Post("/tags", a.createProductTag)
	a.Routes.Product.Get("/tags", a.getProductTags)
	a.Routes.Product.Patch("/tags/{tag_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.patchProductTag))
	a.Routes.Product.Delete("/tags/{tag_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.deleteProductTag))

	// product images
	a.Routes.Product.Post("/images", a.createProductImage)
	a.Routes.Product.Get("/images", a.getProductImages)
	a.Routes.Product.Patch("/images/{image_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.patchProductImage))
	a.Routes.Product.Delete("/images/{image_id:[A-Za-z0-9]+}", a.AdminSessionRequired(a.deleteProductImage))

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

	tagids := mpf.Value["tags"]
	images := mpf.File["images"]

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
