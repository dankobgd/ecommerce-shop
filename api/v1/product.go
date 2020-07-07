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
	msgProductURLParams     = &i18n.Message{ID: "api.product.url.params.app_error", Other: "could not parse url params"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.BaseRoutes.Products.Post("/", a.createProduct)
	a.BaseRoutes.Products.Get("/", a.getProducts)
	a.BaseRoutes.Product.Patch("/", a.patchProduct)
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

	fh := mpf.File["image"][0]
	headers := mpf.File["images"]

	product, productError := a.app.CreateProduct(&p, fh, headers)
	if productError != nil {
		respondError(w, productError)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) {}

func (a *API) patchProduct(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductURLParams, http.StatusInternalServerError, nil))
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
