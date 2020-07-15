package apiv1

import (
	"github.com/dankobgd/ecommerce-shop/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// API is the v1 wrapper
type API struct {
	app    *app.App
	Routes *Routes
}

// Routes contains all api route definitions
type Routes struct {
	Root     chi.Router // ''
	API      chi.Router // 'api/v1'
	Users    chi.Router // 'api/v1/users'
	User     chi.Router // 'api/v1/users/{user_id:[A-Za-z0-9]+}'
	Products chi.Router // 'api/v1/products'
	Product  chi.Router // 'api/v1/products/{product_id:[A-Za-z0-9]+}'
}

// Init inits the API
func Init(a *app.App, r *chi.Mux) {
	api := &API{
		app:    a,
		Routes: &Routes{},
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	api.Routes.Root = r
	api.Routes.API = api.Routes.Root.Route("/api/v1", nil)
	api.Routes.Users = api.Routes.API.Route("/users", nil)
	api.Routes.User = api.Routes.Users.Route("/{user_id:[A-Za-z0-9]+}", nil)
	api.Routes.Products = api.Routes.API.Route("/products", nil)
	api.Routes.Product = api.Routes.Products.Route("/{product_id:[A-Za-z0-9]+}", nil)

	InitUser(api)
	InitProducts(api)
}
