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
	Root       chi.Router // ''
	API        chi.Router // 'api/v1'
	Users      chi.Router // 'api/v1/users'
	User       chi.Router // 'api/v1/users/{user_id:[A-Za-z0-9]+}'
	Products   chi.Router // 'api/v1/products'
	Product    chi.Router // 'api/v1/products/{product_id:[A-Za-z0-9]+}'
	Orders     chi.Router // 'api/v1/orders'
	Order      chi.Router // 'api/v1/orders/{order_id:[A-Za-z0-9]+}'
	Categories chi.Router // 'api/v1/categories'
	Category   chi.Router // 'api/v1/categories/{category_id:[A-Za-z0-9]+}'
	Brands     chi.Router // 'api/v1/brands'
	Brand      chi.Router // 'api/v1/brands/{brand_id:[A-Za-z0-9]+}'
	Tags       chi.Router // 'api/v1/tags'
	Tag        chi.Router // 'api/v1/tags/{tag_id:[A-Za-z0-9]+}'
	Reviews    chi.Router // 'api/v1/reviews'
	Review     chi.Router // 'api/v1/reviews/{review_id:[A-Za-z0-9]+}'
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
	api.Routes.Orders = api.Routes.API.Route("/orders", nil)
	api.Routes.Order = api.Routes.Orders.Route("/{order_id:[A-Za-z0-9]+}", nil)
	api.Routes.Categories = api.Routes.API.Route("/categories", nil)
	api.Routes.Category = api.Routes.Categories.Route("/{category_id:[A-Za-z0-9]+}", nil)
	api.Routes.Brands = api.Routes.API.Route("/brands", nil)
	api.Routes.Brand = api.Routes.Brands.Route("/{brand_id:[A-Za-z0-9]+}", nil)
	api.Routes.Tags = api.Routes.API.Route("/tags", nil)
	api.Routes.Tag = api.Routes.Tags.Route("/{tag_id:[A-Za-z0-9]+}", nil)
	api.Routes.Reviews = api.Routes.API.Route("/reviews", nil)
	api.Routes.Review = api.Routes.Reviews.Route("/{review_id:[A-Za-z0-9]+}", nil)

	InitUser(api)
	InitProducts(api)
	InitOrder(api)
	InitCategories(api)
	InitBrands(api)
	InitTags(api)
	InitReviews(api)
}
