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
func (s PgProductStore) BulkInsertTags(tags []*model.ProductTag) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_tag (product_id, name, created_at, updated_at) VALUES(:product_id, :name, :created_at, :updated_at) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, tags)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}

	return ids, nil
}

// BulkInsertImages inserts multiple images in the db
func (s PgProductStore) BulkInsertImages(images []*model.ProductImage) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_image (product_id, url) VALUES(:product_id, :url) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, images)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	return ids, nil
}

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	return nil
}

// Save inserts the new product in the db
func (s PgProductStore) Save(p *model.Product) (*model.Product, *model.AppErr) {
	m := map[string]interface{}{
		"name":              p.Name,
		"slug":              p.Slug,
		"image_url":         p.ImageURL,
		"description":       p.Description,
		"price":             p.Price,
		"stock":             p.Stock,
		"sku":               p.SKU,
		"is_featured":       p.IsFeatured,
		"created_at":        p.CreatedAt,
		"updated_at":        p.UpdatedAt,
		"deleted_at":        p.DeletedAt,
		"cat_name":          p.Category.Name,
		"cat_slug":          p.Category.Slug,
		"cat_description":   p.Category.Description,
		"brand_name":        p.Brand.Name,
		"brand_slug":        p.Brand.Slug,
		"brand_type":        p.Brand.Type,
		"brand_description": p.Brand.Description,
		"brand_website_url": p.Brand.WebsiteURL,
		"brand_email":       p.Brand.Email,
		"brand_created_at":  p.Brand.CreatedAt,
		"brand_updated_at":  p.Brand.UpdatedAt,
	}

	q := `WITH prod_ins AS (
		INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at, deleted_at)
		VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at, :deleted_at)
		RETURNING id as pid
		),
		cat_ins AS (
		INSERT INTO public.product_category (product_id, name, slug, description)
		VALUES ((SELECT pid FROM prod_ins), :cat_name, :cat_slug, :cat_description)
		RETURNING id as cid
		)
		INSERT INTO public.product_brand (product_id, name, slug, type, description, email, website_url, created_at, updated_at)
		VALUES ((SELECT pid FROM prod_ins), :brand_name, :brand_slug, :brand_type, :brand_description, :brand_email, :brand_website_url, :brand_created_at, :brand_updated_at)
		RETURNING (SELECT pid from prod_ins), (SELECT cid from cat_ins), id as bid`

	var pid, cid, bid int64
	rows, err := s.db.NamedQuery(q, m)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		rows.Scan(&pid, &cid, &bid)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintError(err) {
			return nil, model.NewAppErr("PgProductStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraint, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}

	p.ID = pid
	p.Brand.ID = bid
	p.Category.ID = cid

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
