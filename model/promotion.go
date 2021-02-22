package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidPromotion           = &i18n.Message{ID: "model.promotion.validate.app_error", Other: "invalid promotion data"}
	msgValidatePromotionPromoCode = &i18n.Message{ID: "model.promotion.validate.promo_code.app_error", Other: "invalid promo code"}
	msgValidatePromotionType      = &i18n.Message{ID: "model.promotion.validate.type.app_error", Other: "invalid promotion type"}
	msgValidatePromotionAmount    = &i18n.Message{ID: "model.promotion.validate.amount.app_error", Other: "invalid promotion amount value"}
	msgValidatePromotionStartsAt  = &i18n.Message{ID: "model.promotion.validate.starts_at.app_error", Other: "invalid promotion starts_at timestamp"}
	msgValidatePromotionEndsAt    = &i18n.Message{ID: "model.promotion.validate.ends_at.app_error", Other: "invalid promotion ends_at timestamp"}
	msgValidatePromotionCreatedAt = &i18n.Message{ID: "model.promotion.validate.created_at.app_error", Other: "invalid promotion created_at timestamp"}
	msgValidatePromotionUpdatedAt = &i18n.Message{ID: "model.promotion.validate.updated_at.app_error", Other: "invalid promotion updated_at timestamp"}
)

// Promotion is the promotion model (discount for order)
type Promotion struct {
	TotalRecordsCount
	PromoCode   string    `json:"promo_code" db:"promo_code"`
	Type        string    `json:"type" db:"type"`
	Amount      int       `json:"amount" db:"amount"`
	Description string    `json:"description,omitempty" db:"description"`
	StartsAt    time.Time `json:"starts_at" db:"starts_at"`
	EndsAt      time.Time `json:"ends_at" db:"ends_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PromotionDetail is is the promotion association
type PromotionDetail struct {
	UserID    int64  `json:"user_id" db:"user_id"`
	PromoCode string `json:"promo_code" db:"promo_code"`
}

// PreSave will fill timestamps and other defaults
func (p *Promotion) PreSave() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
}

// PreUpdate sets the update timestamp
func (p *Promotion) PreUpdate() {
	p.UpdatedAt = time.Now()
}

// Validate validates the category and returns an error if it doesn't pass criteria
func (p *Promotion) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if p.PromoCode == "" {
		errs.Add(Invalid("promo_code", l, msgValidatePromotionPromoCode))
	}
	if p.Type == "" {
		errs.Add(Invalid("type", l, msgValidatePromotionType))
	}
	if p.Amount == 0 {
		errs.Add(Invalid("amount", l, msgValidatePromotionAmount))
	}
	if p.StartsAt.IsZero() {
		errs.Add(Invalid("starts_at", l, msgValidatePromotionStartsAt))
	}
	if p.EndsAt.IsZero() || p.EndsAt.Before(p.StartsAt) {
		errs.Add(Invalid("ends_at", l, msgValidatePromotionEndsAt))
	}
	if p.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidatePromotionCreatedAt))
	}
	if p.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidatePromotionUpdatedAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Promotion", msgInvalidPromotion, "", errs)
	}
	return nil
}

// PromotionPatch is the category patch model
type PromotionPatch struct {
	Type        *string    `json:"type,omitempty"`
	Amount      *int       `json:"amount,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
}

// Patch patches the category fields that are provided
func (p *Promotion) Patch(patch *PromotionPatch) {
	if patch.Type != nil {
		p.Type = *patch.Type
	}
	if patch.Amount != nil {
		p.Amount = *patch.Amount
	}
	if patch.Description != nil {
		p.Description = *patch.Description
	}
	if patch.StartsAt != nil {
		p.StartsAt = *patch.StartsAt
	}
	if patch.EndsAt != nil {
		p.EndsAt = *patch.EndsAt
	}
}

// Validate validates the promotion patch and returns an error if it doesn't pass criteria
func (patch *PromotionPatch) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if patch.Type != nil && *patch.Type == "" {
		errs.Add(Invalid("type", l, msgValidatePromotionType))
	}
	if patch.Amount != nil && *patch.Amount == 0 {
		errs.Add(Invalid("amount", l, msgValidatePromotionAmount))
	}
	if patch.StartsAt != nil && patch.StartsAt.IsZero() {
		errs.Add(Invalid("starts_at", l, msgValidatePromotionStartsAt))
	}
	if patch.EndsAt != nil && patch.EndsAt.IsZero() || patch.EndsAt.Before(*patch.StartsAt) {
		errs.Add(Invalid("ends_at", l, msgValidatePromotionEndsAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Promotion", msgInvalidPromotion, "", errs)
	}
	return nil
}

// PromotionPatchFromJSON decodes the input and returns the PromotionPatch
func PromotionPatchFromJSON(data io.Reader) (*PromotionPatch, error) {
	var patch *PromotionPatch
	err := json.NewDecoder(data).Decode(&patch)
	return patch, err
}

// PromotionFromJSON decodes the input and returns the Promotion
func PromotionFromJSON(data io.Reader) (*Promotion, error) {
	var p *Promotion
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// ToJSON converts Category to json string
func (p *Promotion) ToJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

// IsActive checks if the promotion is currently active
func (p *Promotion) IsActive(t time.Time) bool {
	return t.After(p.StartsAt) && t.Before(p.EndsAt)
}
