package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// ProductSQL ...
type ProductSQL struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	ImageURL    string     `db:"image_url"`
	Description string     `db:"description"`
	Price       int        `db:"price"`
	Stock       int        `db:"stock"`
	SKU         string     `db:"sku"`
	IsFeatured  bool       `db:"is_featured"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`

	*BrandSQL
	*CategorySQL
	Tags   types.JSONText `db:"tags"`
	images types.JSONText `db:"images"`
}

// BrandSQL ...
type BrandSQL struct {
	BID          int64     `db:"brand_id"`
	BProductID   int64     `db:"brand_product_id"`
	BName        string    `db:"brand_name"`
	BSlug        string    `db:"brand_slug"`
	BType        string    `db:"brand_type"`
	BDescription string    `db:"brand_description"`
	BEmail       string    `db:"brand_email"`
	BWebsiteURL  string    `db:"brand_website_url"`
	BCreatedAt   time.Time `db:"brand_created_at"`
	BUpdatedAt   time.Time `db:"brand_updated_at"`
}

// CategorySQL ...
type CategorySQL struct {
	CID          int64  `db:"category_id"`
	CProductID   int64  `db:"category_product_id"`
	CName        string `db:"category_name"`
	CSlug        string `db:"category_slug"`
	CDescription string `db:"category_description"`
}

// ToProduct converts the sql rewsult struct to product
func (psql *ProductSQL) ToProduct() *Product {
	return &Product{
		ID:          psql.ID,
		Name:        psql.Name,
		Slug:        psql.Slug,
		ImageURL:    psql.ImageURL,
		Description: psql.Description,
		Price:       psql.Price,
		Stock:       psql.Stock,
		SKU:         psql.SKU,
		IsFeatured:  psql.IsFeatured,
		CreatedAt:   psql.CreatedAt,
		UpdatedAt:   psql.UpdatedAt,
		DeletedAt:   psql.DeletedAt,
		Brand: &ProductBrand{
			ID:          psql.BID,
			ProductID:   psql.BProductID,
			Name:        psql.BName,
			Slug:        psql.BSlug,
			Type:        psql.BType,
			Description: psql.BDescription,
			Email:       psql.BEmail,
			WebsiteURL:  psql.BWebsiteURL,
			CreatedAt:   psql.BCreatedAt,
			UpdatedAt:   psql.BUpdatedAt,
		},
		Category: &ProductCategory{
			ID:          psql.CID,
			ProductID:   psql.CProductID,
			Name:        psql.CName,
			Slug:        psql.CSlug,
			Description: psql.CDescription,
		},
	}
}
