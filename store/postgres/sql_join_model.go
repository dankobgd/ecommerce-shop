package postgres

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
)

// productJoin is temp join type
type productJoin struct {
	model.Product
	*BrandJoin
	*CategoryJoin
}

// BrandJoin is temp join type
type BrandJoin struct {
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

// CategoryJoin is temp join type
type CategoryJoin struct {
	CID          int64  `db:"category_id"`
	CProductID   int64  `db:"category_product_id"`
	CName        string `db:"category_name"`
	CSlug        string `db:"category_slug"`
	CDescription string `db:"category_description"`
}

func (pj *productJoin) ToProduct() *model.Product {
	return &model.Product{
		ID:          pj.ID,
		Name:        pj.Name,
		Slug:        pj.Slug,
		ImageURL:    pj.ImageURL,
		Description: pj.Description,
		Price:       pj.Price,
		Stock:       pj.Stock,
		SKU:         pj.SKU,
		IsFeatured:  pj.IsFeatured,
		CreatedAt:   pj.CreatedAt,
		UpdatedAt:   pj.UpdatedAt,
		Brand: &model.ProductBrand{
			ID:          pj.BID,
			ProductID:   pj.BProductID,
			Name:        pj.BName,
			Slug:        pj.BSlug,
			Type:        pj.BType,
			Description: pj.BDescription,
			Email:       pj.BEmail,
			WebsiteURL:  pj.BWebsiteURL,
			CreatedAt:   pj.BCreatedAt,
			UpdatedAt:   pj.BUpdatedAt,
		},
		Category: &model.ProductCategory{
			ID:          pj.CID,
			ProductID:   pj.CProductID,
			Name:        pj.CName,
			Slug:        pj.CSlug,
			Description: pj.CDescription,
		},
	}
}
