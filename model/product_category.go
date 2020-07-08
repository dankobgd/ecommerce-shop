package model

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidCategory             = &i18n.Message{ID: "model.product_category.validate.app_error", Other: "invalid category data"}
	msgValidateCategoryID          = &i18n.Message{ID: "model.product_category.validate.id.app_error", Other: "invalid category id"}
	msgValidateCategoryProductID   = &i18n.Message{ID: "model.product_category.validate.product_id.app_error", Other: "invalid category product id"}
	msgValidateCategoryName        = &i18n.Message{ID: "model.product_category.validate.name.app_error", Other: "invalid category name"}
	msgValidateCategorySlug        = &i18n.Message{ID: "model.product_category.validate.created_at.app_error", Other: "invalid category created_at timestamp"}
	msgValidateCategoryDescription = &i18n.Message{ID: "model.product_category.validate.updated_at.app_error", Other: "invalid category updated_at timestamp"}
)

// ProductCategory is the category of the product
type ProductCategory struct {
	ID          int64  `json:"id" db:"category_id" schema:"-"`
	ProductID   int64  `json:"-" db:"category_product_id" schema:"-"`
	Name        string `json:"name" db:"category_name" schema:"name"`
	Slug        string `json:"slug" db:"category_slug" schema:"slug"`
	Description string `json:"description" db:"category_description" schema:"description"`
}

func isValidCategoryName(s string) bool {
	var categories []string
	f, _ := os.Open("./migrations/data/category.csv")
	reader := csv.NewReader(f)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		categories = append(categories, line[0])
	}
	for _, cat := range categories {
		if cat == s {
			return true
		}
	}
	return false
}

// Validate validates the user and returns an error if it doesn't pass criteria
func (pc *ProductCategory) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pc.ID != 0 {
		errs.Add(Invalid("category.id", l, msgValidateCategoryID))
	}
	if pc.ProductID != 0 {
		errs.Add(Invalid("category.product_id", l, msgValidateCategoryProductID))
	}
	if pc.Name == "" || !isValidCategoryName(pc.Name) {
		errs.Add(Invalid("category.name", l, msgValidateCategoryName))
	}
	if pc.Slug == "" {
		errs.Add(Invalid("category.slug", l, msgValidateCategorySlug))
	}
	if pc.Description == "" {
		errs.Add(Invalid("category.description", l, msgValidateCategoryDescription))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductCategory", msgInvalidCategory, "", errs)
	}
	return nil
}
