package apiv1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgOrderItemsDataFromJSON = &i18n.Message{ID: "api.order.create_order.json.app_error", Other: "could not parse order item json data"}
)

// InitOrder inits the order routes
func InitOrder(a *API) {
	a.Routes.Orders.Post("/", a.SessionRequired(a.createOrder))
	a.Routes.Order.Get("/", a.SessionRequired(a.getOrder))
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	items, e := model.OrderItemsDataFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createOrder", model.ErrInternal, locale.GetUserLocalizer("en"), msgOrderItemsDataFromJSON, http.StatusInternalServerError, nil))
		return
	}

	o := &model.Order{UserID: uid, CreatedAt: time.Now()}
	order, err := a.app.CreateOrder(o)
	if err != nil {
		respondError(w, err)
		return
	}

	ids := make([]int64, 0)
	for _, x := range items {
		ids = append(ids, x.ProductID)
	}
	products, err := a.app.GetProductsbyIDS(ids)
	if err != nil {
		respondError(w, err)
		return
	}

	totalPrice := 0
	orderDetails := make([]*model.OrderDetail, 0)

	for i, p := range products {
		dtl := &model.OrderDetail{
			OrderID:       order.ID,
			ProductID:     p.ID,
			Quantity:      items[i].Quantity,
			OriginalPrice: p.Price,
			OriginalSKU:   p.SKU,
		}

		totalPrice += p.Price * items[i].Quantity
		orderDetails = append(orderDetails, dtl)
	}

	if err := a.app.CreateOrderDetails(orderDetails); err != nil {
		respondError(w, err)
		return
	}

	order.Total = totalPrice
	updatedOrder, err := a.app.UpdateOrder(order.ID, order)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, updatedOrder)
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
