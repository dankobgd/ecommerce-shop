package postgres

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/jmoiron/sqlx/types"
)

// productJoin is temp join type
type productJoin struct {
	model.Product
	Tsv  string `json:"-" db:"tsv"`
	Rank string `json:"-" db:"rank"`
	*PricingJoin
	*BrandJoin
	*CategoryJoin
}

// PricingJoin is temp join type
type PricingJoin struct {
	PID         int64     `db:"pricing_id"`
	PProductID  int64     `db:"pricing_product_id"`
	PPrice      int       `db:"pricing_price"`
	PSaleStarts time.Time `db:"pricing_sale_starts"`
	PSaleEnds   time.Time `db:"pricing_sale_ends"`
}

// BrandJoin is temp join type
type BrandJoin struct {
	BID           int64     `db:"brand_id"`
	BName         string    `db:"brand_name"`
	BSlug         string    `db:"brand_slug"`
	BType         string    `db:"brand_type"`
	BDescription  string    `db:"brand_description"`
	BEmail        string    `db:"brand_email"`
	BLogo         string    `db:"brand_logo"`
	BLogoPublicID string    `db:"brand_logo_public_id"`
	BWebsiteURL   string    `db:"brand_website_url"`
	BCreatedAt    time.Time `db:"brand_created_at"`
	BUpdatedAt    time.Time `db:"brand_updated_at"`
}

// CategoryJoin is temp join type
type CategoryJoin struct {
	CID           int64           `db:"category_id"`
	CName         string          `db:"category_name"`
	CSlug         string          `db:"category_slug"`
	CLogo         string          `db:"category_logo"`
	CLogoPublicID string          `db:"category_logo_public_id"`
	CDescription  string          `db:"category_description"`
	CIsFeatured   bool            `db:"category_is_featured"`
	CProperties   *types.JSONText `db:"category_properties"`
	CCreatedAt    time.Time       `db:"category_created_at"`
	CUpdatedAt    time.Time       `db:"category_updated_at"`
}

func (pj *productJoin) ToProduct() *model.Product {
	return &model.Product{
		BrandID:           pj.BID,
		CategoryID:        pj.CID,
		TotalRecordsCount: pj.TotalRecordsCount,
		ID:                pj.ID,
		Name:              pj.Name,
		Slug:              pj.Slug,
		ImageURL:          pj.ImageURL,
		ImagePublicID:     pj.ImagePublicID,
		Description:       pj.Description,
		InStock:           pj.InStock,
		SKU:               pj.SKU,
		IsFeatured:        pj.IsFeatured,
		CreatedAt:         pj.CreatedAt,
		UpdatedAt:         pj.UpdatedAt,
		Properties:        pj.Properties,
		// Pricing: &model.ProductPricing{
		// 	Price:      pj.PPrice,
		// 	SaleStarts: pj.PSaleStarts,
		// 	SaleEnds:   pj.PSaleEnds,
		// },
		Brand: &model.Brand{
			ID:           pj.BID,
			Name:         pj.BName,
			Slug:         pj.BSlug,
			Type:         pj.BType,
			Description:  pj.BDescription,
			Email:        pj.BEmail,
			WebsiteURL:   pj.BWebsiteURL,
			Logo:         pj.BLogo,
			LogoPublicID: pj.BLogoPublicID,
			CreatedAt:    pj.BCreatedAt,
			UpdatedAt:    pj.BUpdatedAt,
		},
		Category: &model.Category{
			ID:          pj.CID,
			Name:        pj.CName,
			Slug:        pj.CSlug,
			Logo:        pj.CLogo,
			Description: pj.CDescription,
			IsFeatured:  pj.CIsFeatured,
			Properties:  pj.CProperties,
			CreatedAt:   pj.CCreatedAt,
			UpdatedAt:   pj.CUpdatedAt,
		},
	}
}

// addressJoin is temp join type
type addressJoin struct {
	model.Address
	model.UserAddress `db:"user_address"`
}

func (aj *addressJoin) ToAddress() *model.Address {
	return &model.Address{
		Line1:     aj.Line1,
		Line2:     aj.Line2,
		City:      aj.City,
		Country:   aj.Country,
		State:     aj.State,
		ZIP:       aj.ZIP,
		Latitude:  aj.Latitude,
		Longitude: aj.Longitude,
		Phone:     aj.Phone,
		CreatedAt: aj.CreatedAt,
		UpdatedAt: aj.UpdatedAt,
		DeletedAt: aj.DeletedAt,
	}
}
