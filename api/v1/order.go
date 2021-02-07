package apiv1

import (
	"bytes"
	"fmt"
	"io"
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
	a.Routes.Order.Get("/details", a.SessionRequired(a.getOrderDetails))
	a.Routes.Order.Get("/details/pdf", a.SessionRequired(a.getOrderDetailsPDF))
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	orderData, e := model.OrderRequestDataFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createOrder", model.ErrInternal, locale.GetUserLocalizer("en"), msgOrderItemsDataFromJSON, http.StatusInternalServerError, nil))
		return
	}

	order, err := a.app.CreateOrder(uid, orderData)
	if err != nil {
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

func (a *API) getOrderDetails(w http.ResponseWriter, r *http.Request) {
	oid, e := strconv.ParseInt(chi.URLParam(r, "order_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getOrderDetails", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	details, err := a.app.GetOrderDetails(oid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, details)
}

func (a *API) getOrderDetailsPDF(w http.ResponseWriter, r *http.Request) {
	oid, e := strconv.ParseInt(chi.URLParam(r, "order_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getOrderDetailsPDF", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	order, err := a.app.GetOrder(oid)
	if err != nil {
		respondError(w, err)
		return
	}

	details, err := a.app.GetOrderDetails(oid)
	if err != nil {
		respondError(w, err)
		return
	}

	user, err := a.app.GetUserByID(order.UserID)
	if err != nil {
		respondError(w, err)
		return
	}

	pdf, pdfErr := a.app.GenerateOrderDetailsPDF(order, details, user)
	if pdfErr != nil {
		respondError(w, pdfErr)
	}

	t := order.CreatedAt
	filename := fmt.Sprintf("%02d-%02d-%d-%02d:%02d:%02d-order-invoice.pdf", t.Month(), t.Day(), t.Year(), t.Hour(), t.Minute(), t.Second())

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	io.Copy(w, bytes.NewReader(pdf.Bytes()))
}
