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
	msgGetOrderDetail         = &i18n.Message{ID: "store.postgres.order_detail.get.app_error", Other: "could not get order detail"}
	msgCreateOrderDetail      = &i18n.Message{ID: "store.postgres.order_detail.create.app_error", Other: "could not create order detail"}
)

// BulkInsert inserts multiple order details into the db
func (s *PgOrderDetailStore) BulkInsert(items []*model.OrderDetail) *model.AppErr {
	if _, err := s.db.NamedExec(`INSERT INTO public.order_detail (order_id, product_id, quantity, original_price, original_sku) VALUES (:order_id, :product_id, :quantity, :original_price, :original_sku)`, items); err != nil {
		return model.NewAppErr("PgOrderDetailStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertOrderDetails, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save creates the new order detail
func (s *PgOrderDetailStore) Save(o *model.OrderDetail) (*model.OrderDetail, *model.AppErr) {
	return nil, nil
}

// Get gets the order detail by id
func (s *PgOrderDetailStore) Get(id int64) (*model.OrderDetail, *model.AppErr) {
	return nil, nil
}
