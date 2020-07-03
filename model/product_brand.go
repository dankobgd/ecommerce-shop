package model

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidBrand             = &i18n.Message{ID: "model.product_brand.validate.app_error", Other: "invalid brand data"}
	msgValidateBrandID          = &i18n.Message{ID: "model.product_brand.validate.id.app_error", Other: "invalid brand id"}
	msgValidateBrandProductID   = &i18n.Message{ID: "model.product_brand.validate.product_id.app_error", Other: "invalid brand product id"}
	msgValidateBrandName        = &i18n.Message{ID: "model.product_brand.validate.name.app_error", Other: "invalid brand name"}
	msgValidateBrandSlug        = &i18n.Message{ID: "model.product_brand.validate.slug.app_error", Other: "invalid brand slug"}
	msgValidateBrandDescription = &i18n.Message{ID: "model.product_brand.validate.description.app_error", Other: "invalid brand description"}
	msgValidateBrandType        = &i18n.Message{ID: "model.product_brand.validate.type.app_error", Other: "invalid brand type"}
	msgValidateBrandEmail       = &i18n.Message{ID: "model.product_brand.validate.email.app_error", Other: "invalid brand email"}
	msgValidateBrandWebsiteURL  = &i18n.Message{ID: "model.product_brand.validate.website_url.app_error", Other: "invalid brand website URL"}
	msgValidateBrandCrAt        = &i18n.Message{ID: "model.product_brand.validate.created_at.app_error", Other: "invalid brand created_at timestamp"}
	msgValidateBrandUpAt        = &i18n.Message{ID: "model.product_brand.validate.updated_at.app_error", Other: "invalid brand updated_at timestamp"}
)

// ProductBrand is the brand of the product
type ProductBrand struct {
	ID          int64     `json:"id" db:"id" schema:"-"`
	ProductID   int64     `json:"product_id" db:"product_id" schema:"-"`
	Name        string    `json:"name" db:"name" schema:"brand_name"`
	Slug        string    `json:"slug" db:"slug" schema:"brand_slug"`
	Type        string    `json:"type" db:"type" schema:"brand_type"`
	Description string    `json:"description" db:"description" schema:"brand_description"`
	Email       string    `json:"email" db:"email" schema:"brand_email"`
	WebsiteURL  string    `json:"website_url" db:"website_url" schema:"brand_website_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" schema:"-"`
}

// SetProductID sets product id
func (pb *ProductBrand) SetProductID(pid int64) {
	pb.ProductID = pid
}

// PreSave will fill timestamps and will set productID
func (pb *ProductBrand) PreSave() {
	pb.CreatedAt = time.Now()
	pb.UpdatedAt = pb.CreatedAt
	pb.Email = NormalizeEmail(pb.Email)
}

// PreUpdate sets the update timestamp
func (pb *ProductBrand) PreUpdate() {
	pb.UpdatedAt = time.Now()
	pb.Email = NormalizeEmail(pb.Email)
}

// Validate validates the brand and returns an error if it doesn't pass criteria
func (pb *ProductBrand) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pb.ID != 0 {
		errs.Add(Invalid("brand.id", l, msgValidateBrandID))
	}
	if pb.ProductID != 0 {
		errs.Add(Invalid("brand.product_id", l, msgValidateBrandProductID))
	}
	if pb.Name == "" {
		errs.Add(Invalid("brand.name", l, msgValidateBrandName))
	}
	if pb.Slug == "" {
		errs.Add(Invalid("brand.slug", l, msgValidateBrandSlug))
	}
	if pb.Type == "" {
		errs.Add(Invalid("brand.type", l, msgValidateBrandType))
	}
	if pb.Description == "" {
		errs.Add(Invalid("brand.description", l, msgValidateBrandDescription))
	}
	if pb.Email == "" {
		errs.Add(Invalid("brand.email", l, msgValidateBrandEmail))
	}
	if pb.WebsiteURL == "" {
		errs.Add(Invalid("brand.website_url", l, msgValidateBrandWebsiteURL))
	}
	if pb.CreatedAt.IsZero() {
		errs.Add(Invalid("brand.created_at", l, msgValidateBrandCrAt))
	}
	if pb.UpdatedAt.IsZero() {
		errs.Add(Invalid("brand.updated_at", l, msgValidateBrandUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductBrand", msgInvalidBrand, "", errs)
	}
	return nil
}
