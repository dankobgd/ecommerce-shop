package postgres

import (
	"fmt"
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
	q := `INSERT INTO public.product (name, slug, image_url, description, price, stock, sku, is_featured, created_at, updated_at)
	VALUES (:name, :slug, :image_url, :description, :price, :stock, :sku, :is_featured, :created_at, :updated_at)`

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
	q := `SELECT 
	  p.*,
		br.id AS brand_id,
		br.product_id AS brand_product_id,
		br.name AS brand_name,
		br.slug AS brand_slug,
		br.type AS brand_type,
		br.description AS brand_description,
		br.email AS brand_email,
		br.website_url AS brand_website_url,
		br.created_at AS brand_created_at,
		br.updated_at AS brand_updated_at,	
		cat.id AS category_id,
		cat.product_id AS category_product_id,
		cat.name AS category_name,
		cat.slug AS category_slug,
		cat.description AS category_description		
		FROM public.product AS p
		LEFT JOIN product_brand AS br ON p.id = br.product_id
		LEFT JOIN product_category AS cat ON p.id = cat.product_id
		WHERE p.id = $1
		GROUP BY p.id, br.id, cat.id`

	var psql model.ProductSQL
	if err := s.db.Get(&psql, q, id); err != nil {
		fmt.Println(err)
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	qtags := `SELECT * FROM public.product_tag WHERE public.product_tag.product_id = $1`
	var tags []*model.ProductTag
	if err := s.db.Select(&tags, qtags, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	qimgs := `SELECT * FROM public.product_image WHERE public.product_image.product_id = $1`
	var imgs []*model.ProductImage
	if err := s.db.Select(&imgs, qimgs, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	p := psql.ToProduct()
	p.Tags = tags
	p.Images = imgs

	return p, nil
}

// GetAll returns all products
func (s PgProductStore) GetAll() ([]*model.Product, *model.AppErr) {
	return []*model.Product{}, nil
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
