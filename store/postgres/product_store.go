package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx"
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
	msgSaveProduct        = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product"}
	msgGetProduct         = &i18n.Message{ID: "store.postgres.product.get.app_error", Other: "could not get product"}
	msgGetProducts        = &i18n.Message{ID: "store.postgres.product.get_all.app_error", Other: "could not get products"}
	msgUpdateProduct      = &i18n.Message{ID: "store.postgres.product.update.app_error", Other: "could not update product"}
	msgDeleteProduct      = &i18n.Message{ID: "store.postgres.product.delete.app_error", Other: "could not delete product"}
	msgInvalidColumn      = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product, invalid foreign key value"}
)

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	q := `INSERT INTO public.product (name, brand_id, category_id, slug, image_url, description, price, in_stock, sku, is_featured, created_at, updated_at) 
	VALUES (:name, :brand_id, :category_id, :slug, :image_url, :description, :price, :in_stock, :sku, :is_featured, :created_at, :updated_at)`

	if _, err := s.db.NamedExec(q, products); err != nil {
		return model.NewAppErr("PgProductStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertProducts, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new product in the db
func (s PgProductStore) Save(p *model.Product) (*model.Product, *model.AppErr) {
	q := `INSERT INTO public.product (name, brand_id, category_id, slug, image_url, description, price, in_stock, sku, is_featured, created_at, updated_at, properties)
		VALUES (:name, :brand_id, :category_id, :slug, :image_url, :description, :price, :in_stock, :sku, :is_featured, :created_at, :updated_at, :properties) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, p)
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}

	if err := rows.Err(); err != nil {
		if IsForeignKeyConstraintViolationError(err) {
			return nil, model.NewAppErr("PgProductStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgInvalidColumn, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}

	product, pErr := s.Get(id)
	if pErr != nil {
		return nil, model.NewAppErr("PgProductStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}

	return product, nil
}

// Get gets one product by id
func (s PgProductStore) Get(id int64) (*model.Product, *model.AppErr) {
	q := `SELECT p.*,
   b.name AS brand_name,
   b.slug AS brand_slug,
   b.type AS brand_type,
   b.description AS brand_description,
   b.email AS brand_email,
   b.logo AS brand_logo,
   b.website_url AS brand_website_url,
   b.created_at AS brand_created_at,
   b.updated_at AS brand_updated_at,
   c.name AS category_name,
   c.slug AS category_slug,
	 c.description AS category_description,
	 c.logo AS category_logo,
	 c.created_at AS category_created_at,
   c.updated_at AS category_updated_at
   FROM public.product p
   LEFT JOIN brand b ON p.brand_id = b.id
	 LEFT JOIN category c ON p.category_id = c.id
	 WHERE p.id = $1
   GROUP BY p.id, b.id, c.id`

	var pj productJoin
	if err := s.db.Get(&pj, q, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	return pj.ToProduct(), nil
}

// GetAll returns all products
func (s PgProductStore) GetAll(limit, offset int) ([]*model.Product, *model.AppErr) {
	q := `SELECT 
	 COUNT(*) OVER() AS total_count,
	 p.*,
   b.name AS brand_name,
   b.slug AS brand_slug,
   b.type AS brand_type,
   b.description AS brand_description,
	 b.email AS brand_email,
	 b.logo AS brand_logo,
   b.website_url AS brand_website_url,
   b.created_at AS brand_created_at,
   b.updated_at AS brand_updated_at,
   c.name AS category_name,
   c.slug AS category_slug,
	 c.description AS category_description,
	 c.logo AS category_logo,
	 c.created_at AS category_created_at,
   c.updated_at AS category_updated_at
   FROM public.product p
   LEFT JOIN brand b ON p.brand_id = b.id
   LEFT JOIN category c ON p.category_id = c.id
	 GROUP BY p.id, b.id, c.id
	 LIMIT $1 OFFSET $2`

	var pj []productJoin
	if err := s.db.Select(&pj, q, limit, offset); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProducts, http.StatusInternalServerError, nil)
	}

	products := make([]*model.Product, 0)
	for _, x := range pj {
		products = append(products, x.ToProduct())
	}

	return products, nil
}

// ListByIDS returns all products where ids are in slice
func (s PgProductStore) ListByIDS(ids []int64) ([]*model.Product, *model.AppErr) {
	q, args, err := sqlx.In(`
	 SELECT p.*,
   b.name AS brand_name,
   b.slug AS brand_slug,
   b.type AS brand_type,
   b.description AS brand_description,
	 b.email AS brand_email,
	 b.logo AS brand_logo,
   b.website_url AS brand_website_url,
   b.created_at AS brand_created_at,
   b.updated_at AS brand_updated_at,
   c.name AS category_name,
   c.slug AS category_slug,
	 c.description AS category_description,
	 c.logo AS category_logo,
	 c.created_at AS category_created_at,
   c.updated_at AS category_updated_at
   FROM public.product p
   LEFT JOIN brand b ON p.brand_id = b.id
   LEFT JOIN category c ON p.category_id = c.id
   WHERE p.id IN (?)`, ids)

	if err != nil {
		return nil, model.NewAppErr("PgProductStore.ListByIDS", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProducts, http.StatusInternalServerError, nil)
	}

	var pj []productJoin
	if err := s.db.Select(&pj, s.db.Rebind(q), args...); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProducts, http.StatusInternalServerError, nil)
	}

	products := make([]*model.Product, 0)
	for _, x := range pj {
		products = append(products, x.ToProduct())
	}

	return products, nil
}

// Update updates the product
func (s PgProductStore) Update(id int64, p *model.Product) (*model.Product, *model.AppErr) {
	q := `UPDATE public.product SET name=:name, brand_id=:brand_id, category_id=:category_id, slug=:slug, image_url=:image_url, description=:description, price=:price, in_stock=:in_stock, sku=:sku, is_featured=:is_featured, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, p); err != nil {
		return nil, model.NewAppErr("PgProductStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProduct, http.StatusInternalServerError, nil)
	}
	return p, nil
}

// Delete hard deletes the product from db
func (s PgProductStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.product WHERE id = :id`, map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgProductStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProduct, http.StatusInternalServerError, nil)
	}
	return nil
}

// GetFeatured returns featured products
func (s PgProductStore) GetFeatured(limit, offset int) ([]*model.Product, *model.AppErr) {
	q := `SELECT 
	COUNT(*) OVER() AS total_count,
	p.*,
	b.name AS brand_name,
	b.slug AS brand_slug,
	b.type AS brand_type,
	b.description AS brand_description,
	b.email AS brand_email,
	b.logo AS brand_logo,
	b.website_url AS brand_website_url,
	b.created_at AS brand_created_at,
	b.updated_at AS brand_updated_at,
	c.name AS category_name,
	c.slug AS category_slug,
	c.description AS category_description,
	c.logo AS category_logo,
	c.created_at AS category_created_at,
	c.updated_at AS category_updated_at
	FROM public.product p
	LEFT JOIN brand b ON p.brand_id = b.id
	LEFT JOIN category c ON p.category_id = c.id
	WHERE p.is_featured = true
	GROUP BY p.id, b.id, c.id
	LIMIT $1 OFFSET $2`

	var pj []productJoin
	if err := s.db.Select(&pj, q, limit, offset); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProducts, http.StatusInternalServerError, nil)
	}

	products := make([]*model.Product, 0)
	for _, x := range pj {
		products = append(products, x.ToProduct())
	}

	return products, nil
}

// GetReviews returns all reviews
func (s PgProductStore) GetReviews(id int64) ([]*model.Review, *model.AppErr) {
	var reviews = make([]*model.Review, 0)
	if err := s.db.Select(&reviews, `SELECT * FROM public.product_review WHERE product_id = $1`, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetReviews", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReviews, http.StatusInternalServerError, nil)
	}

	return reviews, nil
}
