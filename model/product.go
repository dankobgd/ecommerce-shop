package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/random"
	"github.com/jmoiron/sqlx/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// FileUploadSizeLimit for image upload
const FileUploadSizeLimit int64 = 10 << 20

// error msgs
var (
	msgInvalidProduct             = &i18n.Message{ID: "model.product.validate.app_error", Other: "invalid product data"}
	msgValidateProductID          = &i18n.Message{ID: "model.product.validate.id.app_error", Other: "invalid product id"}
	msgValidateProductBrandID     = &i18n.Message{ID: "model.product.validate.brand_id.app_error", Other: "invalid product brand id"}
	msgValidateProductCategoryID  = &i18n.Message{ID: "model.product.validate.category_id.app_error", Other: "invalid product category id"}
	msgValidateProductName        = &i18n.Message{ID: "model.product.validate.name.app_error", Other: "invalid product name"}
	msgValidateProductSlug        = &i18n.Message{ID: "model.product.validate.slug.app_error", Other: "invalid product slug"}
	msgValidateProductDescription = &i18n.Message{ID: "model.product.validate.description.app_error", Other: "invalid product description"}
	msgValidateProductPrice       = &i18n.Message{ID: "model.product.validate.price.app_error", Other: "invalid product price"}
	msgValidateProductSKU         = &i18n.Message{ID: "model.product.validate.sku.app_error", Other: "invalid product sku"}
	msgValidateProductCrAt        = &i18n.Message{ID: "model.product.validate.created_at.app_error", Other: "invalid created_at timestamp"}
	msgValidateProductUpAt        = &i18n.Message{ID: "model.product.validate.updated_at.app_error", Other: "invalid updated_at timestamp"}
)

// Product represents the shop product model
type Product struct {
	TotalRecordsCount
	ID          int64          `json:"id" db:"id" schema:"-"`
	BrandID     int64          `json:"-" db:"brand_id" schema:"brand_id"`
	CategoryID  int64          `json:"-" db:"category_id" schema:"category_id"`
	Name        string         `json:"name" db:"name" schema:"name"`
	Slug        string         `json:"slug" db:"slug" schema:"slug"`
	ImageURL    string         `json:"image_url" db:"image_url" schema:"-"`
	Description string         `json:"description" db:"description" schema:"description"`
	Price       int            `json:"price" db:"price" schema:"price"`
	InStock     bool           `json:"in_stock" db:"in_stock" schema:"in_stock"`
	SKU         string         `json:"sku" db:"sku" schema:"-"`
	IsFeatured  bool           `json:"is_featured" db:"is_featured" schema:"is_featured"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at" schema:"-"`
	Properties  types.JSONText `json:"properties"`
	Brand       *Brand         `json:"brand"`
	Category    *Category      `json:"category"`
}

// ProductPatch is the product patch model
type ProductPatch struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	ImageURL    *string `json:"image_url"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
	InStock     *bool   `json:"in_stock"`
	IsFeatured  *bool   `json:"is_featured"`
}

// Patch patches the product fields that are provided
func (p *Product) Patch(patch *ProductPatch) {
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Slug != nil {
		p.Slug = *patch.Slug
	}
	if patch.ImageURL != nil {
		p.ImageURL = *patch.ImageURL
	}
	if patch.Description != nil {
		p.Description = *patch.Description
	}
	if patch.Price != nil {
		p.Price = *patch.Price
	}
	if patch.InStock != nil {
		p.InStock = *patch.InStock
	}
	if patch.IsFeatured != nil {
		p.IsFeatured = *patch.IsFeatured
	}
}

// ProductPatchFromJSON decodes the input and returns the ProductPatch
func ProductPatchFromJSON(data io.Reader) (*ProductPatch, error) {
	var pp *ProductPatch
	err := json.NewDecoder(data).Decode(&pp)
	return pp, err
}

// SetImageURL sets the product image url
func (p *Product) SetImageURL(url string) {
	p.ImageURL = url
}

// ProductFromJSON decodes the input and returns the Product
func ProductFromJSON(data io.Reader) (*Product, error) {
	var p *Product
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// ToJSON converts Product to json string
func (p *Product) ToJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

// PreSave will set missing defaults and fill CreatedAt and UpdatedAt times
func (p *Product) PreSave() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	p.SKU = random.AlphaNumeric(64)
}

// PreUpdate sets the update timestamp
func (p *Product) PreUpdate() {
	p.UpdatedAt = time.Now()
}

// Validate validates the product and returns an error if it doesn't pass criteria
func (p *Product) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if p.ID != 0 {
		errs.Add(Invalid("id", l, msgValidateProductID))
	}
	if p.BrandID == 0 {
		errs.Add(Invalid("brand_id", l, msgValidateProductBrandID))
	}
	if p.CategoryID == 0 {
		errs.Add(Invalid("category_id", l, msgValidateProductCategoryID))
	}
	if p.Name == "" {
		errs.Add(Invalid("name", l, msgValidateProductName))
	}
	if p.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateProductSlug))
	}
	if p.Description == "" {
		errs.Add(Invalid("description", l, msgValidateProductDescription))
	}
	if p.Price == 0 {
		errs.Add(Invalid("price", l, msgValidateProductPrice))
	}
	if p.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateProductCrAt))
	}
	if p.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateProductUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Product", msgInvalidProduct, "", errs)
	}

	return nil
}
