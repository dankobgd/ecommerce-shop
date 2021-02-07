package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgOrderDetailStore is the postgres implementation
type PgOrderDetailStore struct {
	PgStore
}

// NewPgOrderDetailStore creates the new order detail store
func NewPgOrderDetailStore(pgst *PgStore) store.OrderDetailStore {
	return &PgOrderDetailStore{*pgst}
}

var (
	msgBulkInsertOrderDetails = &i18n.Message{ID: "store.postgres.order_detail.bulk_insert.app_error", Other: "could not bulk insert order details"}
	msgGetOrderDetails        = &i18n.Message{ID: "store.postgres.order_details.get.app_error", Other: "could not get order details"}
	msgSaveOrderDetail        = &i18n.Message{ID: "store.postgres.order_detail.create.app_error", Other: "could not create order detail"}
)

// BulkInsert inserts multiple order details into the db
func (s *PgOrderDetailStore) BulkInsert(items []*model.OrderDetail) *model.AppErr {
	if _, err := s.db.NamedExec(`INSERT INTO public.order_detail (order_id, product_id, quantity, history_price, history_sku) VALUES (:order_id, :product_id, :quantity, :history_price, :history_sku)`, items); err != nil {
		return model.NewAppErr("PgOrderDetailStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertOrderDetails, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save creates the new order detail
func (s *PgOrderDetailStore) Save(o *model.OrderDetail) (*model.OrderDetail, *model.AppErr) {
	if _, err := s.db.NamedExec(`INSERT INTO public.order_detail (order_id, product_id, quantity, history_price, history_sku) VALUES (:order_id, :product_id, :quantity, :history_price, :history_sku)`, o); err != nil {
		return nil, model.NewAppErr("PgOrderDetailStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveOrderDetail, http.StatusInternalServerError, nil)
	}
	return o, nil
}

// GetAll gets the all order details
func (s *PgOrderDetailStore) GetAll(orderID int64) ([]*model.OrderInfo, *model.AppErr) {
	var ods = make([]*model.OrderInfo, 0)

	q := `SELECT * FROM public.order_detail od LEFT JOIN public.product p ON od.product_id = p.id WHERE od.order_id = $1`
	if err := s.db.Select(&ods, q, orderID); err != nil {
		return nil, model.NewAppErr("PgBrandStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetOrderDetails, http.StatusInternalServerError, nil)
	}
	return ods, nil
}
