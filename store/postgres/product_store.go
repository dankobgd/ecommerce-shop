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

// BulkInsertTags inserts multiple tags in the db
func (s PgProductStore) BulkInsertTags(tags []*model.ProductTag) *model.AppErr {
	q := `INSERT INTO public.product_tag (product_id, name, created_at, updated_at) VALUES(:product_id, :name, :created_at, :updated_at)`

	if _, err := s.db.NamedExec(q, tags); err != nil {
		return model.NewAppErr("PgUserStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertUsers, http.StatusInternalServerError, nil)
	}
	return nil
}

// BulkInsertImages inserts multiple images in the db
func (s PgProductStore) BulkInsertImages(images []*model.ProductImage) *model.AppErr {
	q := `INSERT INTO public.product_image (product_id, url) VALUES(:product_id, :url)`

	if _, err := s.db.NamedExec(q, images); err != nil {
		return model.NewAppErr("PgUserStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertUsers, http.StatusInternalServerError, nil)
	}
	return nil
}

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	return nil
}

// Save inserts the new product in the db
func (s PgProductStore) Save(p *model.Product, pb *model.ProductBrand, pc *model.ProductCategory) (*model.Product, *model.AppErr) {
	m := map[string]interface{}{
		"name":        p.Name,
		"slug":        p.Slug,
		"image_url":   p.ImageURL,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"sku":         p.SKU,
		"is_featured": p.IsFeatured,
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
		"deleted_at":  p.DeletedAt,

		"cat_name":        pc.Name,
		"cat_slug":        pc.Slug,
		"cat_description": pc.Description,

		"brand_name":        pb.Name,
		"brand_slug":        pb.Slug,
		"brand_type":        pb.Type,
		"brand_description": pb.Description,
		"brand_website_url": pb.WebsiteURL,
		"brand_email":       pb.Email,
		"brand_created_at":  pb.CreatedAt,
		"brand_updated_at":  pb.UpdatedAt,
	}

	q := `WITH product_insert AS (
		INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at, deleted_at)
		VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at, :deleted_at)
		RETURNING id as product_id
		),
		category_insert AS (
		INSERT INTO public.product_category (product_id, name, slug, description)
		VALUES ((SELECT product_id FROM product_insert), :cat_name, :cat_slug, :cat_description)
		RETURNING id
		)
		INSERT INTO public.product_brand (product_id, name, slug, type, description, email, website_url, created_at, updated_at)
		VALUES ((SELECT product_id FROM product_insert), :brand_name, :brand_slug, :brand_type, :brand_description, :brand_email, :brand_website_url, :brand_created_at, :brand_updated_at)
		RETURNING product_id`

	var id int64
	rows, err := s.db.NamedQuery(q, m)
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

	p.ID = id
	return p, nil
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
