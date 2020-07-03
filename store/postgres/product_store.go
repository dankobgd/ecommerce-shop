package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgProductStore is the postgres implementation
type PgProductStore struct {
	PgStore
}

// NewPgProductStore creates the new user store
func NewPgProductStore(pgst *PgStore) store.ProductStore {
	return &PgProductStore{*pgst}
}

var (
	msgSaveProduct = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product to db"}
)

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	return nil
}

// Save inserts the new user in the db
func (s PgProductStore) Save(product *model.Product) (*model.Product, *model.AppErr) {
	q := `INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at, deleted_at) 
		VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at, :deleted_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, product)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintError(err) {
			return nil, model.NewAppErr("PgProductStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraint, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	product.ID = id
	return product, nil
}

// Get gets one product by id
func (s PgProductStore) Get(id int64) (*model.Product, *model.AppErr) {
	return nil, nil
}

// GetAll returns all products
func (s PgProductStore) GetAll() ([]*model.Product, *model.AppErr) {
	return []*model.Product{}, nil
}

// Update ...
func (s PgProductStore) Update(id int64, u *model.Product) (*model.Product, *model.AppErr) {
	return &model.Product{}, nil
}

// Delete ...
func (s PgProductStore) Delete(id int64) (*model.Product, *model.AppErr) {
	return &model.Product{}, nil
}

// InsertTag saves the product tag info
func (s PgProductStore) InsertTag(tag *model.ProductTag) *model.AppErr {
	q := `INSERT INTO public.product_tag (product_id, name, created_at, updated_at) VALUES (:product_id, :name, :created_at, :updated_at)`
	if _, err := s.db.NamedExec(q, tag); err != nil {
		return model.NewAppErr("PgProduct.InsertTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	return nil
}

// InsertCategory saves the product category info
func (s PgProductStore) InsertCategory(category *model.ProductCategory) *model.AppErr {
	q := `INSERT INTO public.product_category (product_id, name, slug, description) VALUES (:product_id, :name, :slug, :description)`
	if _, err := s.db.NamedExec(q, category); err != nil {
		return model.NewAppErr("PgProduct.InsertCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	return nil
}

// InsertBrand saves the product brand info
func (s PgProductStore) InsertBrand(brand *model.ProductBrand) *model.AppErr {
	q := `INSERT INTO public.product_brand (product_id, name, slug, type, description, email, website_url) VALUES (:product_id, :name, :slug, :type, :description, :email, :website_url)`
	if _, err := s.db.NamedExec(q, brand); err != nil {
		return model.NewAppErr("PgProduct.InsertBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	return nil
}
