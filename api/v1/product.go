package apiv1

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductFileErr   = &i18n.Message{ID: "api.product.create_product.formfile.app_error", Other: "error parsing files"}
	msgProductMultipart = &i18n.Message{ID: "api.product.create_product.multipart.app_error", Other: "could not decode product multipart data"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.BaseRoutes.Products.Post("/", a.createProduct)
	a.BaseRoutes.Products.Get("/", a.getProducts)
}

func (a *API) createProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductMultipart, http.StatusInternalServerError, nil))
		return
	}

	var p model.Product
	var pTag model.ProductTag
	var pCat model.ProductCategory
	var pBrand model.ProductBrand

	model.SchemaDecoder.IgnoreUnknownKeys(true)
	hasError := false
	vals := r.MultipartForm.Value

	if err := model.SchemaDecoder.Decode(&p, vals); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pTag, vals); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pCat, vals); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pBrand, vals); err != nil {
		hasError = true
	}

	if hasError {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductMultipart, http.StatusInternalServerError, nil))
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFileErr, http.StatusInternalServerError, nil))
		return
	}
	defer file.Close()

	product, productError := a.app.CreateProduct(&p, &pTag, &pCat, &pBrand, file, header)
	if productError != nil {
		respondError(w, productError)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) {}
