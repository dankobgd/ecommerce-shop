package postgres

import (
	"fmt"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgAddressStore is the postgres implementation
type PgAddressStore struct {
	PgStore
}

// NewPgAddressStore creates the new order store
func NewPgAddressStore(pgst *PgStore) store.AddressStore {
	return &PgAddressStore{*pgst}
}

var (
	msgSaveAddress   = &i18n.Message{ID: "store.postgres.address.save.app_error", Other: "could not save address"}
	msgGetAddress    = &i18n.Message{ID: "store.postgres.address.get.app_error", Other: "could not get address"}
	msgUpdateAddress = &i18n.Message{ID: "store.postgres.address.update.app_error", Other: "could not update address"}
	msgDeleteAddress = &i18n.Message{ID: "store.postgres.address.save.app_error", Other: "could not delete address"}
)

// Save creates the new address
func (s *PgAddressStore) Save(addr *model.Address, userID int64, addrType model.AddrType) (*model.Address, *model.AppErr) {
	q := `WITH addr_ins AS (
		INSERT INTO public.address (line_1, line_2, city, country, state, zip, latitude, longitude, phone, created_at, updated_at) 
		VALUES (:line_1, :line_2, :city, :country, :state, :zip, :latitude, :longitude, :phone, :created_at, :updated_at) 
		RETURNING id AS addr_id
	),
	addr_type_ins AS (
		INSERT INTO public.address_type (name, address_id) VALUES (:address_type, (SELECT addr_id FROM addr_ins)) 
		RETURNING id AS addr_type_id
	)
	INSERT INTO public.user_address (user_id, address_id, address_type_id) 
	VALUES (:user_id, (SELECT addr_id FROM addr_ins), (SELECT addr_type_id FROM addr_type_ins)) 
	RETURNING (SELECT addr_id FROM addr_ins)
	`

	m := map[string]interface{}{
		"line_1":       addr.Line1,
		"line_2":       addr.Line2,
		"city":         addr.City,
		"country":      addr.Country,
		"state":        addr.State,
		"zip":          addr.ZIP,
		"latitude":     addr.Latitude,
		"longitude":    addr.Longitude,
		"phone":        addr.Phone,
		"created_at":   addr.CreatedAt,
		"updated_at":   addr.UpdatedAt,
		"address_type": addrType.String(),
		"user_id":      userID,
	}

	var id int64
	rows, err := s.db.NamedQuery(q, m)
	if err != nil {
		return nil, model.NewAppErr("PgAddressStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveAddress, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		if IsForeignKeyConstraintViolationError(err) {
			return nil, model.NewAppErr("PgAddressStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgInvalidColumn, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgAddressStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveAddress, http.StatusInternalServerError, nil)
	}

	addr.ID = id
	return addr, nil
}

// Get gets the address
func (s *PgAddressStore) Get(id int64) (*model.Address, *model.AppErr) {
	q := `SELECT * FROM public.address WHERE id = $1`
	var addr model.Address
	if err := s.db.Get(&addr, q, id); err != nil {
		return nil, model.NewAppErr("PgAddressStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetAddress, http.StatusInternalServerError, nil)
	}

	return &addr, nil
}

// Update updates the address
func (s *PgAddressStore) Update(id int64, addr *model.Address) (*model.Address, *model.AppErr) {
	q := `UPDATE public.address SET line_1=:line_1, line_2=:line_2, city=:city, country=:country, state=:state, zip=:zip, latitude=:latitude, longitude=:longitude, phone=:phone, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, addr); err != nil {
		return nil, model.NewAppErr("PgAddressStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateAddress, http.StatusInternalServerError, nil)
	}
	return addr, nil
}

// Delete hard deletes the address
func (s *PgAddressStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.address WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgAddressStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteAddress, http.StatusInternalServerError, nil)
	}
	return nil
}
