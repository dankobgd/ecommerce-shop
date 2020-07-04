package model

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// FileUploadSizeLimit for image upload
const FileUploadSizeLimit int64 = 10 << 20

// error msgs
var (
	msgInvalidProduct             = &i18n.Message{ID: "model.product.validate.app_error", Other: "invalid product data"}
	msgValidateProductID          = &i18n.Message{ID: "model.product.validate.id.app_error", Other: "invalid product id"}
	msgValidateProductBrandID     = &i18n.Message{ID: "model.product.validate.brand.app_error", Other: "invalid product brand id"}
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
	ID          int64      `json:"id" db:"id" schema:"-"`
	Name        string     `json:"name" db:"name" schema:"name"`
	Slug        string     `json:"slug" db:"slug" schema:"slug"`
	ImageURL    string     `json:"image_url" db:"image_url" schema:"-"`
	Description string     `json:"description" db:"description" schema:"description"`
	Price       int        `json:"price" db:"price" schema:"price"`
	Stock       int        `json:"stock" db:"stock" schema:"stock"`
	SKU         string     `json:"sku" db:"sku" schema:"sku"`
	IsFeatured  bool       `json:"is_featured" db:"is_featured" schema:"is_featured"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" schema:"-"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at" schema:"-"`
}

// ProductCreateData is used when creating new products
type ProductCreateData struct {
	P            *Product
	Cat          *ProductCategory
	Brand        *ProductBrand
	Tag          *ProductTag
	TagNames     []string
	ImgFH        *multipart.FileHeader
	ImageHeaders []*multipart.FileHeader
}

// SetImageURL sets the product image url
func (p *Product) SetImageURL(url string) {
	p.ImageURL = url
}

// ProductFromJSON decodes the input and return the Product
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
	if p.SKU == "" {
		errs.Add(Invalid("sku", l, msgValidateProductSKU))
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
