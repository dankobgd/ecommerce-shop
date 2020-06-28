package apiv1

import (
	"github.com/dankobgd/ecommerce-shop/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// API is the v1 wrapper
type API struct {
	app        *app.App
	BaseRoutes *Routes
}

// Routes contains all api route definitions
type Routes struct {
	Root     chi.Router // ''
	APIRoot  chi.Router // 'api/v1'
	Users    chi.Router // 'api/v1/users'
	User     chi.Router // 'api/v1/users/{user_id:[A-Za-z0-9]+}'
	Products chi.Router // 'api/v1/products'
	Product  chi.Router // 'api/v1/products/{product_id:[A-Za-z0-9]+}'
}

// Init inits the API
func Init(a *app.App, r *chi.Mux) {
	api := &API{
		app:        a,
		BaseRoutes: &Routes{},
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	api.BaseRoutes.Root = r
	api.BaseRoutes.APIRoot = api.BaseRoutes.Root.Route("/api/v1", nil)
	api.BaseRoutes.Users = api.BaseRoutes.APIRoot.Route("/users", nil)
	api.BaseRoutes.User = api.BaseRoutes.Users.Route("/{user_id:[A-Za-z0-9]+}", nil)
	api.BaseRoutes.Products = api.BaseRoutes.APIRoot.Route("/products", nil)
	api.BaseRoutes.Product = api.BaseRoutes.Products.Route("/{product_id:[A-Za-z0-9]+}", nil)

	InitUser(api)
	InitProducts(api)
}
