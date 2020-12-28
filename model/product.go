package model

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/utils/is"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/random"
	"github.com/jmoiron/sqlx/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// FutureSaleEndsTime is the products sale end time
var FutureSaleEndsTime = time.Date(2050, 01, 01, 00, 00, 00, 000000000, time.UTC)

// FileUploadSizeLimit for image upload - 3MB
const FileUploadSizeLimit int64 = 3 << 20

// error msgs
var (
	msgInvalidProduct            = &i18n.Message{ID: "model.product.validate.app_error", Other: "invalid product data"}
	msgValidateProductID         = &i18n.Message{ID: "model.product.validate.id.app_error", Other: "invalid product id"}
	msgValidateProductBrandID    = &i18n.Message{ID: "model.product.validate.brand_id.app_error", Other: "invalid product brand id"}
	msgValidateProductCategoryID = &i18n.Message{ID: "model.product.validate.category_id.app_error", Other: "invalid product category id"}
	msgValidateProductName       = &i18n.Message{ID: "model.product.validate.name.app_error", Other: "invalid product name"}
	msgValidateProductSlug       = &i18n.Message{ID: "model.product.validate.slug.app_error", Other: "invalid product slug"}
	msgValidateProductPrice      = &i18n.Message{ID: "model.product.validate.price.app_error", Other: "invalid product price"}
	msgValidateProductSKU        = &i18n.Message{ID: "model.product.validate.sku.app_error", Other: "invalid product sku"}
	msgValidateProductCrAt       = &i18n.Message{ID: "model.product.validate.created_at.app_error", Other: "invalid created_at timestamp"}
	msgValidateProductUpAt       = &i18n.Message{ID: "model.product.validate.updated_at.app_error", Other: "invalid updated_at timestamp"}

	msgInvalidProductPricing            = &i18n.Message{ID: "model.product_price.validate.app_error", Other: "invalid product price data"}
	msgValidateProductPricingID         = &i18n.Message{ID: "model.product_price.validate.id.app_error", Other: "invalid product price id"}
	msgValidateProductPricingProductID  = &i18n.Message{ID: "model.product_price.validate.product_id.app_error", Other: "invalid product price product_id"}
	msgValidateProductPricingPrice      = &i18n.Message{ID: "model.product_price.validate.price_app_error", Other: "invalid product price amount"}
	msgValidateProductPricingSaleStarts = &i18n.Message{ID: "model.product_price.validate.sale_starts.app_error", Other: "invalid product price sale starts"}
	msgValidateProductPricingSaleEnds   = &i18n.Message{ID: "model.product_price.validate.sale_ends.app_error", Other: "invalid product price sale ends"}
	msgValidateProductProperties        = &i18n.Message{ID: "model.product.validate.properties.app_error", Other: "invalid json provided as properties"}
)

// Product represents the shop product model
type Product struct {
	TotalRecordsCount
	ID             int64           `json:"id" db:"id" schema:"-"`
	BrandID        int64           `json:"-" db:"brand_id" schema:"brand_id"`
	CategoryID     int64           `json:"-" db:"category_id" schema:"category_id"`
	Name           string          `json:"name" db:"name" schema:"name"`
	Slug           string          `json:"slug" db:"slug" schema:"slug"`
	ImageURL       string          `json:"image_url" db:"image_url" schema:"-"`
	ImagePublicID  string          `json:"image_public_id" db:"image_public_id" schema:"-"`
	Description    string          `json:"description,omitempty" db:"description" schema:"description"`
	InStock        bool            `json:"in_stock" db:"in_stock" schema:"in_stock"`
	SKU            string          `json:"sku" db:"sku" schema:"-"`
	IsFeatured     bool            `json:"is_featured" db:"is_featured" schema:"is_featured"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at" schema:"-"`
	Properties     *types.JSONText `json:"properties" db:"properties" schema:"-"`
	PropertiesText *string         `json:"-" schema:"properties"`

	*ProductPricing `schema:"-"`
	Brand           *Brand    `json:"brand" schema:"-"`
	Category        *Category `json:"category" schema:"-"`
}

// ProductPatch is the product patch model
type ProductPatch struct {
	BrandID        *int64          `json:"brand_id,omitempty" schema:"brand_id"`
	CategoryID     *int64          `json:"category_id,omitempty" schema:"category_id"`
	Name           *string         `json:"name,omitempty" schema:"name"`
	Slug           *string         `json:"slug,omitempty" schema:"slug"`
	ImageURL       *string         `json:"image_url,omitempty" schema:"-"`
	ImagePublicID  *string         `json:"image_public_id,omitempty" schema:"-"`
	Description    *string         `json:"description,omitempty" schema:"description"`
	InStock        *bool           `json:"in_stock,omitempty" schema:"in_stock"`
	IsFeatured     *bool           `json:"is_featured,omitempty" schema:"is_featured"`
	Properties     *types.JSONText `json:"properties,omitempty" schema:"-"`
	PropertiesText *string         `json:"-" schema:"properties"`
}

// Patch patches the product fields that are provided
func (p *Product) Patch(patch *ProductPatch) {
	if patch.BrandID != nil {
		p.BrandID = *patch.BrandID
	}
	if patch.CategoryID != nil {
		p.CategoryID = *patch.CategoryID
	}
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
	if patch.InStock != nil {
		p.InStock = *patch.InStock
	}
	if patch.IsFeatured != nil {
		p.IsFeatured = *patch.IsFeatured
	}
	if patch.Properties != nil {
		p.Properties = patch.Properties
	}
}

// ProductPatchFromJSON decodes the input and returns the ProductPatch
func ProductPatchFromJSON(data io.Reader) (*ProductPatch, error) {
	var pp *ProductPatch
	err := json.NewDecoder(data).Decode(&pp)
	return pp, err
}

// SetImageDetails sets the product image_url and public_id
func (p *Product) SetImageDetails(details *gocloudinary.ResourceDetails) {
	p.ImageURL = details.SecureURL
	p.ImagePublicID = details.PublicID
}

// SetProperties sets the product properties
func (p *Product) SetProperties(properties *string) {
	if properties != nil {
		props := types.JSONText(*properties)
		p.Properties = &props
	}
}

// SetProperties sets the ProductPatch properties
func (patch *ProductPatch) SetProperties(properties *string) {
	if properties != nil {
		if len(*properties) == 0 {
			patch.Properties = nil
		} else {
			props := types.JSONText(*properties)
			patch.Properties = &props
		}
	}
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
	p.SetProperties(p.PropertiesText)
}

// PreUpdate sets the update timestamp
func (p *Product) PreUpdate() {
	p.UpdatedAt = time.Now()
}

// Validate validates the product and returns an error if it doesn't pass criteria
func (p *Product) Validate(fh *multipart.FileHeader) *AppErr {
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
	if p.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateProductCrAt))
	}
	if p.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateProductUpAt))
	}
	if fh == nil {
		errs.Add(Invalid("image", l, msgValidateProductImage))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("image", l, msgValidateProductImageSize))
	}

	// ideally validate properties against json schema to check for the right keys, values and structure...
	if p.PropertiesText != nil && !is.ValidJSON(*p.PropertiesText) {
		errs.Add(Invalid("properties", l, msgValidateProductProperties))
	}

	if !errs.IsZero() {
		return NewValidationError("Product", msgInvalidProduct, "", errs)
	}

	return nil
}

// Validate validates the product patch and returns an error if it doesn't pass criteria
func (patch *ProductPatch) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if patch.BrandID != nil && *patch.BrandID == 0 {
		errs.Add(Invalid("brand_id", l, msgValidateProductBrandID))
	}
	if patch.CategoryID != nil && *patch.CategoryID == 0 {
		errs.Add(Invalid("category_id", l, msgValidateProductCategoryID))
	}
	if patch.Name != nil && *patch.Name == "" {
		errs.Add(Invalid("name", l, msgValidateProductName))
	}
	if patch.Slug != nil && *patch.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateProductSlug))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("image", l, msgValidateProductImageSize))
	}
	// ideally validate properties against json schema to check for the right keys, values and structure...
	if patch.PropertiesText != nil && len(*patch.PropertiesText) != 0 && !is.ValidJSON(*patch.PropertiesText) {
		errs.Add(Invalid("properties", l, msgValidateProductProperties))
	}

	if !errs.IsZero() {
		return NewValidationError("Product", msgInvalidProduct, "", errs)
	}

	return nil
}

// ProductPricing has info about price discounts
type ProductPricing struct {
	PriceID    int64     `json:"-" db:"price_id"`
	ProductID  int64     `json:"-" db:"product_id"`
	Price      int       `json:"price" db:"price"`
	SaleStarts time.Time `json:"sale_starts" db:"sale_starts"`
	SaleEnds   time.Time `json:"sale_ends" db:"sale_ends"`
}

// Validate validates the ProductPricing and returns an error if it doesn't pass criteria
func (pp *ProductPricing) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pp.PriceID != 0 {
		errs.Add(Invalid("id", l, msgValidateProductPricingID))
	}
	if pp.ProductID == 0 {
		errs.Add(Invalid("product_id", l, msgValidateProductPricingProductID))
	}
	if pp.Price <= 0 {
		errs.Add(Invalid("price", l, msgValidateProductPricingPrice))
	}
	if pp.SaleStarts.IsZero() {
		errs.Add(Invalid("sale_starts", l, msgValidateProductPricingSaleStarts))
	}
	if pp.SaleEnds.IsZero() {
		errs.Add(Invalid("sale_ends", l, msgValidateProductPricingSaleEnds))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductPricing", msgInvalidProductPricing, "", errs)
	}

	return nil
}
