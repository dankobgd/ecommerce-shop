package postgres

import (
	"fmt"
	"net/http"
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
	msgUniqueConstraintProduct = &i18n.Message{ID: "store.postgres.product.save.unique_constraint.app_error", Other: "invalid product foreign key"}
	msgBulkInsertProducts      = &i18n.Message{ID: "store.postgres.product.bulk_insert.app_error", Other: "could not bulk insert products"}
	msgSaveProduct             = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product"}
	msgGetProduct              = &i18n.Message{ID: "store.postgres.product.get.app_error", Other: "could not get product"}
	msgGetProducts             = &i18n.Message{ID: "store.postgres.product.get_all.app_error", Other: "could not get products"}
	msgUpdateProduct           = &i18n.Message{ID: "store.postgres.product.update.app_error", Other: "could not update product"}
	msgDeleteProduct           = &i18n.Message{ID: "store.postgres.product.delete.app_error", Other: "could not delete product"}
	msgInvalidColumn           = &i18n.Message{ID: "store.postgres.product.save.app_error", Other: "could not save product, invalid foreign key value"}
	msgGetPricing              = &i18n.Message{ID: "store.postgres.product.get_pricing.app_error", Other: "could not get latest pricing"}
	msgSavePricing             = &i18n.Message{ID: "store.postgres.product.insert_pricing.app_error", Other: "could not insert pricing data"}
	msgUpdatePricing           = &i18n.Message{ID: "store.postgres.product.update_pricing.app_error", Other: "could not update pricing data"}
	msgBulkDeleteProducts      = &i18n.Message{ID: "store.postgres.product.bulk_delete.app_error", Other: "could not bulk delete products"}
)

// BulkInsert inserts multiple products into db
func (s PgProductStore) BulkInsert(products []*model.Product) *model.AppErr {
	q := `INSERT INTO public.product (name, brand_id, category_id, slug, image_url, image_public_id, description, in_stock, sku, is_featured, created_at, updated_at) 
	VALUES (:name, :brand_id, :category_id, :slug, :image_url, :image_public_id, :description, :in_stock, :sku, :is_featured, :created_at, :updated_at)`

	if _, err := s.db.NamedExec(q, products); err != nil {
		return model.NewAppErr("PgProductStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertProducts, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new product in the db
func (s PgProductStore) Save(p *model.Product) (*model.Product, *model.AppErr) {
	q := `INSERT INTO public.product (name, brand_id, category_id, slug, image_url, image_public_id, description, in_stock, sku, is_featured, created_at, updated_at, properties)
		VALUES (:name, :brand_id, :category_id, :slug, :image_url, :image_public_id, :description, :in_stock, :sku, :is_featured, :created_at, :updated_at, :properties) RETURNING id`

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

	pricing := &model.ProductPricing{
		ProductID:     id,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		SaleStarts:    p.CreatedAt,
		SaleEnds:      model.FutureSaleEndsTime,
	}
	if _, err := s.InsertPricing(pricing); err != nil {
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
	q := `SELECT DISTINCT ON (p.id) p.*,
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
	 c.logo_public_id AS category_logo_public_id,
	 c.properties AS category_properties,
	 c.created_at AS category_created_at,
	 c.updated_at AS category_updated_at,
	 pp.id AS pricing_id,
   pp.product_id AS pricing_product_id,
   pp.price AS pricing_price,
   pp.original_price AS pricing_original_price,
   pp.sale_starts AS pricing_sale_starts,
   pp.sale_ends AS pricing_sale_ends
	 FROM public.product p
	 LEFT JOIN product_pricing pp ON p.id = pp.product_id
   LEFT JOIN brand b ON p.brand_id = b.id
	 LEFT JOIN category c ON p.category_id = c.id
	 WHERE CURRENT_TIMESTAMP BETWEEN pp.sale_starts AND pp.sale_ends
	 AND p.id = $1
	 GROUP BY p.id, b.id, c.id, pp.id
	 ORDER BY p.id, pp.id DESC`

	var pj productJoin
	if err := s.db.Get(&pj, q, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProduct, http.StatusInternalServerError, nil)
	}

	return pj.ToProduct(), nil
}

// GetAll returns all products
func (s PgProductStore) GetAll(filters map[string][]string, limit, offset int) ([]*model.Product, *model.AppErr) {
	baseQuery := `SELECT DISTINCT ON (p.id)
	(SELECT COUNT(DISTINCT product.id) FROM product LEFT JOIN product_pricing on product.id = product_pricing.product_id) AS total_count,
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
	c.properties AS category_properties,
	c.created_at AS category_created_at,
	c.updated_at AS category_updated_at,
	pp.id AS pricing_id,
  pp.product_id AS pricing_product_id,
	pp.price AS pricing_price,
	pp.original_price AS pricing_original_price,
  pp.sale_starts AS pricing_sale_starts,
  pp.sale_ends AS pricing_sale_ends
	FROM public.product p
	LEFT JOIN product_pricing pp ON p.id = pp.product_id
	LEFT JOIN brand b ON p.brand_id = b.id
	LEFT JOIN category c ON p.category_id = c.id	
	LEFT JOIN product_tag pt ON p.id = pt.product_id 
	LEFT JOIN tag t on t.id = pt.tag_id
	WHERE CURRENT_TIMESTAMP BETWEEN pp.sale_starts AND pp.sale_ends`

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
	q, args, err := sqlx.In(`SELECT DISTINCT ON (p.id)
	 (SELECT COUNT(DISTINCT product.id) FROM product LEFT JOIN product_pricing on product.id = product_pricing.product_id WHERE product.id IN (?)) AS total_count,
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
	 c.properties AS category_properties,
	 c.created_at AS category_created_at,
	 c.updated_at AS category_updated_at,
	 pp.id AS pricing_id,
   pp.product_id AS pricing_product_id,
	 pp.price AS pricing_price,
	 pp.original_price AS pricing_original_price,
   pp.sale_starts AS pricing_sale_starts,
   pp.sale_ends AS pricing_sale_ends
	 FROM public.product p
	 LEFT JOIN product_pricing pp ON p.id = pp.product_id
   LEFT JOIN brand b ON p.brand_id = b.id
	 LEFT JOIN category c ON p.category_id = c.id
	 WHERE CURRENT_TIMESTAMP BETWEEN pp.sale_starts AND pp.sale_ends
	 AND p.id IN (?)`, ids, ids)

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
	q := `UPDATE public.product SET brand_id=:brand_id, category_id=:category_id, name=:name, slug=:slug, image_url=:image_url, image_public_id=:image_public_id, description=:description, in_stock=:in_stock, sku=:sku, is_featured=:is_featured, updated_at=:updated_at, properties=:properties WHERE id=:id`
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
	q := `SELECT DISTINCT ON (p.id)
	(SELECT COUNT(DISTINCT product.id) FROM product LEFT JOIN product_pricing on product.id = product_pricing.product_id WHERE product.is_featured = true) AS total_count,
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
	c.properties AS category_properties,
	c.created_at AS category_created_at,
	c.updated_at AS category_updated_at,
	pp.id AS pricing_id,
	pp.product_id AS pricing_product_id,
	pp.price AS pricing_price,
	pp.original_price AS pricing_original_price,
	pp.sale_starts AS pricing_sale_starts,
	pp.sale_ends AS pricing_sale_ends
	FROM public.product p
	LEFT JOIN product_pricing pp ON p.id = pp.product_id
	LEFT JOIN brand b ON p.brand_id = b.id
	LEFT JOIN category c ON p.category_id = c.id
	WHERE CURRENT_TIMESTAMP BETWEEN pp.sale_starts AND pp.sale_ends
	AND p.is_featured = true
	GROUP BY p.id, b.id, c.id, pp.id
	ORDER BY p.id DESC, pp.id DESC
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
func (s PgProductStore) GetReviews(id int64) ([]*model.ProductReview, *model.AppErr) {
	q := `SELECT r.*,
	u.id AS user_id,
	u.first_name AS user_first_name,
	u.last_name AS user_last_name,
	u.username AS user_username,
	u.avatar_url AS user_avatar_url,
	u.avatar_public_id AS user_avatar_public_id
	FROM product_review r LEFT JOIN public.user u ON r.user_id = u.id 
	WHERE r.product_id = $1`

	var rj []reviewJoin
	if err := s.db.Select(&rj, q, id); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetReviews", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReviews, http.StatusInternalServerError, nil)
	}

	var reviews = make([]*model.ProductReview, 0)
	for _, x := range rj {
		reviews = append(reviews, x.ToReview())
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

// GetLatestPricing gets latest pricing record
func (s PgProductStore) GetLatestPricing(pid int64) (*model.ProductPricing, *model.AppErr) {
	q := `SELECT pp.id AS price_id, pp.product_id, pp.price, pp.original_price, pp.sale_starts, pp.sale_ends FROM product_pricing pp WHERE product_id = $1 ORDER BY id DESC LIMIT 1`
	var pricing model.ProductPricing
	if err := s.db.Get(&pricing, q, pid); err != nil {
		return nil, model.NewAppErr("PgProductStore.GetLatestPricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetPricing, http.StatusInternalServerError, nil)
	}

	return &pricing, nil
}

// BulkDelete deletes products with given ids
func (s PgProductStore) BulkDelete(ids []int) *model.AppErr {
	q, args, err := sqlx.In(`DELETE FROM product WHERE id IN (?)`, ids)
	if err != nil {
		return model.NewAppErr("PgProductStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteProducts, http.StatusInternalServerError, nil)
	}

	if _, err := s.db.Exec(s.db.Rebind(q), args...); err != nil {
		return model.NewAppErr("PgProductStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteProducts, http.StatusInternalServerError, nil)
	}

	return nil
}

// InsertPricing inserts the price info into product_pricing
func (s PgProductStore) InsertPricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr) {
	q := `INSERT INTO product_pricing(product_id, price, original_price, sale_starts, sale_ends) VALUES(:product_id, :price, :original_price, :sale_starts, :sale_ends) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, pricing)
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.InsertPricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgSavePricing, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgProductStore.InsertPricing", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintProduct, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgProductStore.InsertPricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgSavePricing, http.StatusInternalServerError, nil)
	}

	pricing.PriceID = id
	return pricing, nil
}

// InsertPricingBulk bulk inserts the price info into product_pricing
func (s PgProductStore) InsertPricingBulk(pricing []*model.ProductPricing) *model.AppErr {
	q := `INSERT INTO product_pricing(product_id, price, original_price, sale_starts, sale_ends) VALUES(:product_id, :price, :original_price, :sale_starts, :sale_ends) RETURNING id`
	if _, err := s.db.NamedExec(q, pricing); err != nil {
		return model.NewAppErr("PgProductStore.InsertPricingBulk", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintProduct, http.StatusInternalServerError, nil)
	}
	return nil
}

// UpdatePricing updates the ProductPricing
func (s PgProductStore) UpdatePricing(pricing *model.ProductPricing) (*model.ProductPricing, *model.AppErr) {
	q := `UPDATE product_pricing SET price=:price, original_price=:original_price, sale_starts=:sale_starts, sale_ends=:sale_ends WHERE id=:price_id`
	if _, err := s.db.NamedExec(q, pricing); err != nil {
		return nil, model.NewAppErr("PgProductStore.UpdatePricing", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdatePricing, http.StatusInternalServerError, nil)
	}
	return pricing, nil
}

// generate dynamic search query for get all products (shop search page with all sidebar filters...)

type prop struct {
	name   string
	values []string
}
type categoryFilter struct {
	category string
	props    []prop
}

func containsCat(arr []categoryFilter, f categoryFilter) (bool, int) {
	for i, x := range arr {
		if x.category == f.category {
			return true, i
		}
	}
	return false, -1
}
func containsProp(arr []prop, p prop) (bool, int) {
	for i, x := range arr {
		if x.name == p.name {
			return true, i
		}
	}
	return false, -1
}

func buildProductsFilterSearchQuery(queryString string, filters map[string][]string, limit, offset int) (string, []interface{}, error) {
	basic := make(map[string][]string, 0)
	specific := make(map[string][]string, 0)

	for filter, val := range filters {
		if filter == "page" || filter == "per_page" || filter == "category" || filter == "brand" || filter == "tag" || filter == "price_min" || filter == "price_max" {
			basic[filter] = val
		} else {
			specific[filter] = val
		}
	}

	query := fmt.Sprintf("%s", queryString)
	var args []interface{}

	// handle price range filters
	min, minOk := basic["price_min"]
	max, maxOk := basic["price_max"]
	if minOk && maxOk {
		query += " AND (pp.price >= ? AND pp.price <= ?)\n"
		args = append(args, min[0], max[0])
	} else if minOk && !maxOk {
		query += " AND pp.price >= ?\n"
		args = append(args, min[0])
	} else if !minOk && maxOk {
		query += " AND pp.price <= ?\n"
		args = append(args, max[0])
	}

	// handle brand filters
	if brand, ok := basic["brand"]; ok {
		query += " AND b.slug IN (?)\n"
		args = append(args, brand)
	}

	// handle tag filters
	if tag, ok := basic["tag"]; ok {
		query += " AND t.slug IN (?)\n"
		args = append(args, tag)
	}

	// handle product category specific property filters
	if category, ok := basic["category"]; ok {
		var filtersList = make([]categoryFilter, 0)

		for _, cat := range category {
			filtersList = append(filtersList, categoryFilter{category: cat, props: make([]prop, 0)})
		}

		propsList := make([]prop, 0)
		for filter, values := range specific {
			full := strings.SplitN(filter, "_", 2)
			catName, propName := full[0], full[1]

			p := prop{name: propName, values: values}
			ok, _ := containsProp(propsList, p)
			if !ok {
				propsList = append(propsList, p)
			}

			f := categoryFilter{category: catName, props: propsList}
			ok2, idx := containsCat(filtersList, f)
			if ok2 {
				filtersList[idx].props = append(filtersList[idx].props, p)
			}
		}

		for i, cat := range filtersList {
			str := ""
			outerClause := ""
			outerParensClose := ""

			if i == 0 {
				outerClause = "AND (\n"
			} else {
				outerClause = "OR\n"
			}
			if i == len(filtersList)-1 {
				outerParensClose = "\n)"
			}

			for _, prop := range cat.props {
				innerClause := ""

				for k, val := range prop.values {
					groupStart := ""
					groupEnd := ""
					if len(prop.values) > 1 {
						if k == 0 {
							groupStart = "("
						}
						if k == len(prop.values)-1 {
							groupEnd = ")"
						}
					}

					if k == 0 {
						innerClause = "AND"
					} else {
						innerClause = "OR"
					}

					str += fmt.Sprintf(" %s %sp.properties->>'%s' = '%s'%s", innerClause, groupStart, prop.name, val, groupEnd)
				}
			}
			query += fmt.Sprintf(" %s (c.slug = '%s'%s)%s", outerClause, cat.category, str, outerParensClose)
		}
	}

	query += " \nGROUP BY p.id, b.id, c.id, pp.id"
	query += " \nORDER BY p.id DESC, pp.id DESC"
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
