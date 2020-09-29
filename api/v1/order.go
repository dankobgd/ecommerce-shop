package apiv1

import (
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/pagination"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgOrderItemsDataFromJSON = &i18n.Message{ID: "api.order.create_order.json.app_error", Other: "could not parse order item json data"}
)

// InitOrder inits the order routes
func InitOrder(a *API) {
	a.Routes.Orders.Post("/", a.SessionRequired(a.createOrder))
	a.Routes.Orders.Get("/", a.SessionRequired(a.getOrders))
	a.Routes.Order.Get("/", a.SessionRequired(a.getOrder))
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	orderData, e := model.OrderRequestDataFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createOrder", model.ErrInternal, locale.GetUserLocalizer("en"), msgOrderItemsDataFromJSON, http.StatusInternalServerError, nil))
		return
	}

	ids := make([]int64, 0)
	for _, x := range orderData.Items {
		ids = append(ids, x.ProductID)
	}
	products, err := a.app.GetProductsbyIDS(ids)
	if err != nil {
		respondError(w, err)
		return
	}

	totalPrice := 0
	for i, p := range products {
		totalPrice += p.Price * orderData.Items[i].Quantity
	}

	order, err := a.app.CreateOrder(uid, orderData, totalPrice)
	if err != nil {
		respondError(w, err)
		return
	}

	orderDetails := make([]*model.OrderDetail, 0)

	for i, p := range products {
		detail := &model.OrderDetail{
			OrderID:       order.ID,
			ProductID:     p.ID,
			Quantity:      orderData.Items[i].Quantity,
			OriginalPrice: p.Price,
			OriginalSKU:   p.SKU,
		}

		orderDetails = append(orderDetails, detail)
	}

	if err := a.app.CreateOrderDetails(orderDetails); err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, order)
}

func (a *API) getOrders(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	orders, err := a.app.GetOrders(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(orders) > 0 {
		totalCount = orders[0].TotalCount
	}
	pages.SetData(orders, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	oid, e := strconv.ParseInt(chi.URLParam(r, "order_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getOrder", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	order, err := a.app.GetOrder(oid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, order)
}
