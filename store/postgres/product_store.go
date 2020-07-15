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
	msgBulkInsertProducts = &i18n.Message{ID: "store.postgres.product.bulk_insert.app_error", Other: "could not bulk insert products"}
	msgSaveProduct        = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product to db"}
	msgInvalidColumn      = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product to db, tried to insert non existing column name"}
	msgGetProduct         = &i18n.Message{ID: "store.postgres.product.get.app_error", Other: "could not get product from db"}
	msgUpdateProduct      = &i18n.Message{ID: "store.postgres.product.update.app_error", Other: "could not update product"}
	msgDeleteProduct      = &i18n.Message{ID: "store.postgres.product.delete.app_error", Other: "could not delete product"}
)

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	q := `INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at) VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at)`

	if _, err := s.db.NamedExec(q, products); err != nil {
		return model.NewAppErr("PgProductStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertProducts, http.StatusInternalServerError, nil)
	}
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
		INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at)
		VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at)
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
		if IsForeignKeyConstraintViolationError(err) {
			return nil, model.NewAppErr("PgProductStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgInvalidColumn, http.StatusInternalServerError, nil)
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
	q := `SELECT p.*,
	b.id AS brand_id,
	b.product_id AS brand_product_id,
	b.name AS brand_name,
	b.slug AS brand_slug,
	b.type AS brand_type,
	b.description AS brand_description,
	b.email AS brand_email,
	b.website_url AS brand_website_url,
	b.created_at AS brand_created_at,
	b.updated_at AS brand_updated_at,
	c.id AS category_id,
	c.product_id AS category_product_id,
	c.name AS category_name,
	c.slug AS category_slug,
	c.description AS category_description
	FROM public.product p
	LEFT JOIN product_brand b ON p.id = b.product_id
	LEFT JOIN product_category c ON p.id = c.product_id
	WHERE p.id = $1
	GROUP BY p.id, b.id, c.id`

	var pj productJoin
	if err := s.db.Get(&pj, q, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	return pj.ToProduct(), nil
}

// GetAll returns all products
func (s PgProductStore) GetAll() ([]*model.Product, *model.AppErr) {
	q := `SELECT p.*,
	b.id AS brand_id,
	b.product_id AS brand_product_id,
	b.name AS brand_name,
	b.slug AS brand_slug,
	b.type AS brand_type,
	b.description AS brand_description,
	b.email AS brand_email,
	b.website_url AS brand_website_url,
	b.created_at AS brand_created_at,
	b.updated_at AS brand_updated_at,
	c.id AS category_id,
	c.product_id AS category_product_id,
	c.name AS category_name,
	c.slug AS category_slug,
	c.description AS category_description	
	FROM public.product p
	LEFT JOIN product_brand b ON p.id = b.product_id
	LEFT JOIN product_category c ON p.id = c.product_id`

	var pj []productJoin
	if err := s.db.Select(&pj, q); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	products := make([]*model.Product, 0)
	for _, x := range pj {
		products = append(products, x.ToProduct())
	}

	return products, nil
}

// Update updates the product
func (s PgProductStore) Update(id int64, p *model.Product) (*model.Product, *model.AppErr) {
	q := `UPDATE public.product SET name=:name, slug=:slug, image_url=:image_url, description=:description, price=:price, stock=:stock, sku=:sku, is_featured=:is_featured, created_at=:created_at, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, p); err != nil {
		return nil, model.NewAppErr("PgProductStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProduct, http.StatusInternalServerError, nil)
	}
	return p, nil
}

// Delete ...
func (s PgProductStore) Delete(id int64) *model.AppErr {
	m := map[string]interface{}{"id": id}
	if _, err := s.db.NamedExec(`DELETE FROM public.product WHERE id = :id`, m); err != nil {
		return model.NewAppErr("PgProductStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProduct, http.StatusInternalServerError, nil)
	}
	return nil
}
