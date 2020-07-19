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
)

// Save creates the new order
func (s *PgOrderStore) Save(o *model.Order, shipAddr *model.Address, billAddr *model.Address) (*model.Order, *model.AppErr) {
	q := `INSERT INTO public.order (user_id, status, total, shipped_at, created_at, billing_address_line_1, billing_address_line_2, billing_address_city, billing_address_country, billing_address_state, billing_address_zip, billing_address_latitude, billing_address_longitude, shipping_address_line_1, shipping_address_line_2, shipping_address_city, shipping_address_country, shipping_address_state, shipping_address_zip, shipping_address_latitude, shipping_address_longitude) 
	VALUES (:user_id, :status, :total, :shipped_at, :created_at, :billing_address_line_1, :billing_address_line_2, :billing_address_city, :billing_address_country, :billing_address_state, :billing_address_zip, :billing_address_latitude, :billing_address_longitude, :shipping_address_line_1, :shipping_address_line_2, :shipping_address_city, :shipping_address_country, :shipping_address_state, :shipping_address_zip, :shipping_address_latitude, :shipping_address_longitude) RETURNING id`

	m := map[string]interface{}{
		"user_id":                    o.UserID,
		"status":                     o.Status,
		"total":                      o.Total,
		"shipped_at":                 o.ShippedAt,
		"created_at":                 o.CreatedAt,
		"billing_address_line_1":     billAddr.Line1,
		"billing_address_line_2":     billAddr.Line2,
		"billing_address_city":       billAddr.City,
		"billing_address_country":    billAddr.Country,
		"billing_address_state":      billAddr.State,
		"billing_address_zip":        billAddr.ZIP,
		"billing_address_latitude":   billAddr.Latitude,
		"billing_address_longitude":  billAddr.Longitude,
		"shipping_address_line_1":    shipAddr.Line1,
		"shipping_address_line_2":    shipAddr.Line2,
		"shipping_address_city":      shipAddr.City,
		"shipping_address_country":   shipAddr.Country,
		"shipping_address_state":     shipAddr.State,
		"shipping_address_zip":       shipAddr.ZIP,
		"shipping_address_latitude":  shipAddr.Latitude,
		"shipping_address_longitude": shipAddr.Longitude,
	}

	var id int64
	rows, err := s.db.NamedQuery(q, m)
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
	o.BillingAddressLine1 = billAddr.Line1
	o.BillingAddressLine2 = billAddr.Line2
	o.BillingAddressCity = billAddr.City
	o.BillingAddressCountry = billAddr.Country
	o.BillingAddressState = billAddr.State
	o.BillingAddressZIP = billAddr.ZIP
	o.ShippingAddressLine1 = shipAddr.Line1
	o.ShippingAddressLine2 = shipAddr.Line2
	o.ShippingAddressCity = shipAddr.City
	o.ShippingAddressCountry = shipAddr.Country
	o.ShippingAddressState = shipAddr.State
	o.ShippingAddressZIP = shipAddr.ZIP
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
func (s *PgOrderStore) Get(id int64) (*model.Order, *model.AppErr) {
	var o model.Order
	if err := s.db.Get(&o, `SELECT * FROM public.order WHERE id = $1`, id); err != nil {
		return nil, model.NewAppErr("PgOrderStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetOrder, http.StatusInternalServerError, nil)
	}
	return &o, nil
}
