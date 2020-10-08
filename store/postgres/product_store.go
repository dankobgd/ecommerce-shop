package postgres

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
func (s PgProductStore) GetAll(filters map[string][]string, limit, offset int) ([]*model.Product, *model.AppErr) {
	baseQuery := `SELECT 
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
	LEFT JOIN product_tag pt ON p.id = pt.product_id 
	LEFT JOIN tag t on t.id = pt.tag_id`

	q, args, _ := buildProductsFilterSearchQuery(baseQuery, filters, limit, offset)

	var pj []productJoin
	if err := s.db.Select(&pj, q, args...); err != nil {
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

// Search returns all fulltext search product results
func (s PgProductStore) Search(filter string) ([]*model.Product, *model.AppErr) {
	q := `SELECT *, ts_rank(tsv, plainto_tsquery($1)) as rank FROM product_search_view WHERE tsv @@ plainto_tsquery($1) order by rank desc limit 200`

	var pj []productJoin
	if err := s.db.Select(&pj, q, filter); err != nil {
		return nil, model.NewAppErr("PgProductStore.Search", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProducts, http.StatusInternalServerError, nil)
	}

	products := make([]*model.Product, 0)
	for _, x := range pj {
		products = append(products, x.ToProduct())
	}

	return products, nil
}

func buildProductsFilterSearchQuery(queryString string, filters map[string][]string, limit, offset int) (string, []interface{}, error) {
	basic := make(map[string][]string, 0)
	specific := make(map[string][]string, 0)

	for filter, val := range filters {
		if filter == "category" || filter == "brand" || filter == "tag" || filter == "price_min" || filter == "price_max" {
			basic[filter] = val
		} else {
			specific[filter] = val
		}
	}

	query := fmt.Sprintf("%s WHERE 1 = 1", queryString)
	var args []interface{}

	// handle price range filters
	min, minOk := basic["price_min"]
	max, maxOk := basic["price_max"]
	if minOk && maxOk {
		query += " AND p.price >= ? AND p.price <= ?"
		args = append(args, min[0], max[0])
	} else if minOk && !maxOk {
		query += " AND p.price >= ?"
		args = append(args, min[0])
	} else if !minOk && maxOk {
		query += " AND p.price <= ?"
		args = append(args, max[0])
	}

	// handle brand filters
	if brand, ok := basic["brand"]; ok {
		query += " AND b.name IN (?)"
		args = append(args, brand)
	}

	// handle tag filters
	if tag, ok := basic["tag"]; ok {
		query += " AND t.name IN (?)"
		args = append(args, tag)
	}

	// handle product category specific property filters
	filtersByCategory := make(map[string]map[string][]string, 0)

	if category, ok := basic["category"]; ok {
		for _, cat := range category {
			filtersByCategory[cat] = make(map[string][]string, 0)
		}

		for filter, vals := range specific {
			full := strings.Split(filter, "_")
			cat := full[0]
			prop := full[1]

			if _, ok := filtersByCategory[cat][prop]; !ok {
				filtersByCategory[cat][prop] = vals
			}
		}

		keys := make([]string, 0)
		for k := range filtersByCategory {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			str := ""
			outerCond := "AND"
			if i != 0 {
				outerCond = "OR"
			}

			for prop, vals := range filtersByCategory[k] {
				for idx, v := range vals {
					innerCond := "AND"
					if idx != 0 {
						innerCond = "OR"
					}
					str += fmt.Sprintf(" %s p.properties->>'%s' = '%s'", innerCond, prop, v)
				}
			}
			query += fmt.Sprintf(" %s (c.name = '%s'%s)", outerCond, k, str)
		}
	}

	query += " GROUP BY p.id, b.id, c.id"
	if limit != 0 {
		query += " LIMIT ?"
		args = append(args, strconv.Itoa(limit))
	}
	if offset != 0 {
		query += " OFFSET ?"
		args = append(args, strconv.Itoa(offset))
	}

	builtQuery, builtQueryArgs, err := sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}
	builtQuery = sqlx.Rebind(sqlx.DOLLAR, builtQuery)

	return builtQuery, builtQueryArgs, nil
}
