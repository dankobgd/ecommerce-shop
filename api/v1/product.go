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
	var pc model.ProductCategory
	var pb model.ProductBrand
	var pt model.ProductTag

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	hasError := false
	if err := model.SchemaDecoder.Decode(&p, mpf.Value); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pt, mpf.Value); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pc, mpf.Value); err != nil {
		hasError = true
	}
	if err := model.SchemaDecoder.Decode(&pb, mpf.Value); err != nil {
		hasError = true
	}

	if hasError {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductMultipart, http.StatusInternalServerError, nil))
		return
	}

	fh := mpf.File["image"]
	headers := mpf.File["images"]

	data := &model.ProductCreateData{
		P:            &p,
		Brand:        &pb,
		Cat:          &pc,
		Tag:          &pt,
		ImgFH:        fh[0],
		ImageHeaders: headers,
		TagNames:     mpf.Value["tag_name"],
	}

	product, productError := a.app.CreateProduct(data)
	if productError != nil {
		respondError(w, productError)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) {}
