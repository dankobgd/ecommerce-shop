package postgres

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
)

// productJoin is temp join type
type productJoin struct {
	model.Product
	Tsv  string `json:"-" db:"tsv"`
	Rank string `json:"-" db:"rank"`
	*BrandJoin
	*CategoryJoin
}

// BrandJoin is temp join type
type BrandJoin struct {
	BID          int64     `db:"brand_id"`
	BName        string    `db:"brand_name"`
	BSlug        string    `db:"brand_slug"`
	BType        string    `db:"brand_type"`
	BDescription string    `db:"brand_description"`
	BEmail       string    `db:"brand_email"`
	BLogo        string    `db:"brand_logo"`
	BWebsiteURL  string    `db:"brand_website_url"`
	BCreatedAt   time.Time `db:"brand_created_at"`
	BUpdatedAt   time.Time `db:"brand_updated_at"`
}

// CategoryJoin is temp join type
type CategoryJoin struct {
	CID          int64     `db:"category_id"`
	CName        string    `db:"category_name"`
	CSlug        string    `db:"category_slug"`
	CLogo        string    `db:"category_logo"`
	CDescription string    `db:"category_description"`
	CIsFeatured  bool      `db:"category_is_featured"`
	CCreatedAt   time.Time `db:"category_created_at"`
	CUpdatedAt   time.Time `db:"category_updated_at"`
}

func (pj *productJoin) ToProduct() *model.Product {
	return &model.Product{
		TotalRecordsCount: pj.TotalRecordsCount,
		ID:                pj.ID,
		Name:              pj.Name,
		Slug:              pj.Slug,
		ImageURL:          pj.ImageURL,
		Description:       pj.Description,
		Price:             pj.Price,
		InStock:           pj.InStock,
		SKU:               pj.SKU,
		IsFeatured:        pj.IsFeatured,
		CreatedAt:         pj.CreatedAt,
		UpdatedAt:         pj.UpdatedAt,
		Properties:        pj.Properties,
		Brand: &model.Brand{
			ID:          pj.BID,
			Name:        pj.BName,
			Slug:        pj.BSlug,
			Type:        pj.BType,
			Description: pj.BDescription,
			Email:       pj.BEmail,
			WebsiteURL:  pj.BWebsiteURL,
			Logo:        pj.BLogo,
			CreatedAt:   pj.BCreatedAt,
			UpdatedAt:   pj.BUpdatedAt,
		},
		Category: &model.Category{
			ID:          pj.CID,
			Name:        pj.CName,
			Slug:        pj.CSlug,
			Logo:        pj.CLogo,
			Description: pj.CDescription,
			IsFeatured:  pj.CIsFeatured,
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
		ID:        0,
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
