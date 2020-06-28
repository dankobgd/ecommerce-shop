package apiv1

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgProductFromJSON = &i18n.Message{ID: "api.product.create_product.json.app_error", Other: "could not decode product json data"}
)

// InitProducts inits the product routes
func InitProducts(a *API) {
	a.BaseRoutes.Products.Post("/", a.createProduct)
	a.BaseRoutes.Products.Get("/", a.getProducts)
}

func (a *API) createProduct(w http.ResponseWriter, r *http.Request) {
	p, e := model.ProductFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createProduct", model.ErrInternal, locale.GetUserLocalizer("en"), msgProductFromJSON, http.StatusInternalServerError, nil))
		return
	}

	product, err := a.app.CreateProduct(p)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (a *API) getProducts(w http.ResponseWriter, r *http.Request) {}
