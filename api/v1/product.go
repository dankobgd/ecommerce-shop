package apiv1

import (
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductPatchFromJSON = &i18n.Message{ID: "api.product.patch_product.app_error", Other: "could not decode product patch data"}
	msgProductFileErr       = &i18n.Message{ID: "api.product.create_product.formfile.app_error", Other: "error parsing files"}
	msgProductMultipart     = &i18n.Message{ID: "api.product.create_product.multipart.app_error", Other: "could not decode product multipart data"}
	msgPatchProduct         = &i18n.Message{ID: "api.product.patch_product.app_error", Other: "could not patch product"}
	msgURLParamErr          = &i18n.Message{ID: "api.product.url.params.app_error", Other: "could not parse URL params"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.Routes.Products.Get("/", a.getProducts)
	a.Routes.Products.Post("/", a.createProduct)
	a.Routes.Products.Get("/tags/{tag_id:[A-Za-z0-9]+}", a.getSingleTag)
	a.Routes.Products.Get("/images/{image_id:[A-Za-z0-9]+}", a.getSingleImage)
	a.Routes.Products.Patch("/tags/{tag_id:[A-Za-z0-9]+}", a.patchProductTag)
	a.Routes.Products.Patch("/images/{image_id:[A-Za-z0-9]+}", a.patchProductImage)
	a.Routes.Products.Delete("/tags/{tag_id:[A-Za-z0-9]+}", a.deleteProductTag)
	a.Routes.Products.Delete("/images/{image_id:[A-Za-z0-9]+}", a.deleteProductImage)

	a.Routes.Product.Get("/", a.getProduct)
	a.Routes.Product.Patch("/", a.patchProduct)
	a.Routes.Product.Delete("/", a.deleteProduct)
	a.Routes.Product.Get("/tags", a.getProductTags)
	a.Routes.Product.Get("/images", a.getProductImages)
}

func (a *API) createProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductMultipart, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	var p model.Product
	if err := model.SchemaDecoder.Decode(&p, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductMultipart, http.StatusInternalServerError, nil))
		return
	}

	tagnames := mpf.Value["tags"]
	headers := mpf.File["images"]
	fh := mpf.File["image"][0]

	tags := make([]*model.ProductTag, 0)
	for _, t := range tagnames {
		tags = append(tags, &model.ProductTag{Name: model.NewString(t)})
	}

	product, pErr := a.app.CreateProduct(&p, fh, headers, tags)
	if pErr != nil {
		respondError(w, pErr)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (a *API) patchProduct(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.ProductPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductPatchFromJSON, http.StatusInternalServerError, nil))
		return
	}

	uprod, pErr := a.app.PatchProduct(pid, patch)
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
	products, err := a.app.GetProducts()
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, products)
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

func (a *API) getSingleTag(w http.ResponseWriter, r *http.Request) {
	tid, e := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getSingleTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	tag, err := a.app.GetTag(tid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, tag)
}

func (a *API) getSingleImage(w http.ResponseWriter, r *http.Request) {
	imgID, e := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getSingleImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	img, err := a.app.GetImage(imgID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, img)
}

func (a *API) patchProductTag(w http.ResponseWriter, r *http.Request) {
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

	utag, pErr := a.app.PatchProductTag(tid, patch)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, utag)
}

func (a *API) patchProductImage(w http.ResponseWriter, r *http.Request) {
	imgID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.ProductImagePatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductPatchFromJSON, http.StatusInternalServerError, nil))
		return
	}

	uprod, pErr := a.app.PatchProductImage(imgID, patch)
	if pErr != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, uprod)
}

func (a *API) deleteProductTag(w http.ResponseWriter, r *http.Request) {
	tid, err := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteProductTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteProductTag(tid); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}

func (a *API) deleteProductImage(w http.ResponseWriter, r *http.Request) {
	imgID, err := strconv.ParseInt(chi.URLParam(r, "image_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteProductImage", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteProductImage(imgID); err != nil {
		respondError(w, err)
		return
	}
	respondOK(w)
}
