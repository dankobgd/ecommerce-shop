package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgBrandStore is the postgres implementation
type PgBrandStore struct {
	PgStore
}

// NewPgBrandStore creates the new brand store
func NewPgBrandStore(pgst *PgStore) store.BrandStore {
	return &PgBrandStore{*pgst}
}

var (
	msgUniqueConstraintBrand = &i18n.Message{ID: "store.postgres.brand.save.unique_constraint.app_error", Other: "invalid brand foreign key"}
	msgSaveBrand             = &i18n.Message{ID: "store.postgres.brand.save.app_error", Other: "could not save brand"}
	msgUpdateBrand           = &i18n.Message{ID: "store.postgres.brand.update.app_error", Other: "could not update brand"}
	msgBulkInsertBrands      = &i18n.Message{ID: "store.postgres.brand.bulk.insert.app_error", Other: "could not bulk insert brands"}
	msgGetBrand              = &i18n.Message{ID: "store.postgres.brand.get.app_error", Other: "could not get the brand"}
	msgGetBrands             = &i18n.Message{ID: "store.postgres.brand.get.app_error", Other: "could not get the brand"}
	msgDeleteBrand           = &i18n.Message{ID: "store.postgres.brand.delete.app_error", Other: "could not delete brand"}
)

// BulkInsert inserts multiple brands in the db
func (s PgBrandStore) BulkInsert(brands []*model.Brand) *model.AppErr {
	q := `INSERT INTO public.brand(name, slug, type, logo, description, email, website_url, created_at, updated_at) VALUES(:name, :slug, :type, :logo, :description, :email, :website_url, :created_at, :updated_at) RETURNING id`

	if _, err := s.db.NamedExec(q, brands); err != nil {
		return model.NewAppErr("PgBrandStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertBrands, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new brand in the db
func (s PgBrandStore) Save(brand *model.Brand) (*model.Brand, *model.AppErr) {
	q := `INSERT INTO public.brand(name, slug, type, logo, description, email, website_url, created_at, updated_at) VALUES(:name, :slug, :type, :logo, :description, :email, :website_url, :created_at, :updated_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, brand)
	if err != nil {
		return nil, model.NewAppErr("PgBrandStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveBrand, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgBrandStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintBrand, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgBrandStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveBrand, http.StatusInternalServerError, nil)
	}

	brand.ID = id
	return brand, nil
}

// Update updates the brand
func (s PgBrandStore) Update(id int64, brand *model.Brand) (*model.Brand, *model.AppErr) {
	q := `UPDATE public.brand SET name=:name, slug=:slug, type=:type, description=:description, email=:email, website_url=:website_url, logo=:logo, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, brand); err != nil {
		return nil, model.NewAppErr("PgBrandStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateBrand, http.StatusInternalServerError, nil)
	}
	return brand, nil
}

// Get gets one brand by id
func (s PgBrandStore) Get(id int64) (*model.Brand, *model.AppErr) {
	var brand model.Brand
	if err := s.db.Get(&brand, "SELECT * FROM public.brand WHERE id = $1", id); err != nil {
		return nil, model.NewAppErr("PgBrandStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetBrand, http.StatusInternalServerError, nil)
	}
	return &brand, nil
}

// GetAll returns all brands
func (s PgBrandStore) GetAll(limit, offset int) ([]*model.Brand, *model.AppErr) {
	var brands = make([]*model.Brand, 0)
	if err := s.db.Select(&brands, `SELECT COUNT(*) OVER() AS total_count, * FROM public.brand LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgBrandStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetBrands, http.StatusInternalServerError, nil)
	}

	return brands, nil
}

// Delete deletes the brand
func (s PgBrandStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.brand WHERE id = :id", map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgBrandStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteBrand, http.StatusInternalServerError, nil)
	}
	return nil
}
