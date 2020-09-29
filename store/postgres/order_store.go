package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgOrderStore is the postgres implementation
type PgOrderStore struct {
	PgStore
}

// NewPgOrderStore creates the new order store
func NewPgOrderStore(pgst *PgStore) store.OrderStore {
	return &PgOrderStore{*pgst}
}

var (
	msgSaveOrder   = &i18n.Message{ID: "store.postgres.order.save.app_error", Other: "could not save order"}
	msgUpdateOrder = &i18n.Message{ID: "store.postgres.order.update.app_error", Other: "could not update order"}
	msgGetOrder    = &i18n.Message{ID: "store.postgres.order.get.app_error", Other: "could not get order"}
	msgGetOrders   = &i18n.Message{ID: "store.postgres.orders.get.app_error", Other: "could not get orders"}
)

// Save creates the new order
func (s PgOrderStore) Save(o *model.Order) (*model.Order, *model.AppErr) {
	q := `INSERT INTO public.order (user_id, status, total, shipped_at, created_at, billing_address_line_1, billing_address_line_2, billing_address_city, billing_address_country, billing_address_state, billing_address_zip, billing_address_latitude, billing_address_longitude, shipping_address_line_1, shipping_address_line_2, shipping_address_city, shipping_address_country, shipping_address_state, shipping_address_zip, shipping_address_latitude, shipping_address_longitude) 
	VALUES (:user_id, :status, :total, :shipped_at, :created_at, :billing_address_line_1, :billing_address_line_2, :billing_address_city, :billing_address_country, :billing_address_state, :billing_address_zip, :billing_address_latitude, :billing_address_longitude, :shipping_address_line_1, :shipping_address_line_2, :shipping_address_city, :shipping_address_country, :shipping_address_state, :shipping_address_zip, :shipping_address_latitude, :shipping_address_longitude) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, o)
	if err != nil {
		return nil, model.NewAppErr("PgOrderStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveOrder, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgOrderStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintUser, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgOrderStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveUser, http.StatusInternalServerError, nil)
	}

	o.ID = id
	return o, nil
}

// Update updates the product
func (s PgOrderStore) Update(id int64, o *model.Order) (*model.Order, *model.AppErr) {
	if _, err := s.db.NamedExec(`UPDATE public.order SET status=:status, total=:total, shipped_at=:shipped_at WHERE id=:id`, o); err != nil {
		return nil, model.NewAppErr("PgOrderStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateOrder, http.StatusInternalServerError, nil)
	}
	return o, nil
}

// Get gets the order by id
func (s PgOrderStore) Get(id int64) (*model.Order, *model.AppErr) {
	var o model.Order
	if err := s.db.Get(&o, `SELECT * FROM public.order WHERE id = $1`, id); err != nil {
		return nil, model.NewAppErr("PgOrderStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetOrder, http.StatusInternalServerError, nil)
	}
	return &o, nil
}

// GetAll returns all orders
func (s PgOrderStore) GetAll(limit, offset int) ([]*model.Order, *model.AppErr) {
	var orders = make([]*model.Order, 0)
	if err := s.db.Select(&orders, `SELECT COUNT(*) OVER() AS total_count, * FROM public.order LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgOrderStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetOrders, http.StatusInternalServerError, nil)
	}

	return orders, nil
}

// Delete deletes the order
func (s PgOrderStore) Delete(id int64) *model.AppErr {
	return nil
}
