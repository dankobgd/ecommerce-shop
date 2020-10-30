package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgPromotionStore is the postgres implementation
type PgPromotionStore struct {
	PgStore
}

// NewPgPromotionStore creates the new promotion store
func NewPgPromotionStore(pgst *PgStore) store.PromotionStore {
	return &PgPromotionStore{*pgst}
}

var (
	msgUniqueConstraintPromotion = &i18n.Message{ID: "store.postgres.promotion.save.unique_constraint.app_error", Other: "promotion with given promo_code already exists"}
	msgSavePromotion             = &i18n.Message{ID: "store.postgres.promotion.save.app_error", Other: "could not save promotion"}
	msgUpdatePromotion           = &i18n.Message{ID: "store.postgres.promotion.update.app_error", Other: "could not update promotion"}
	msgBulkInsertPromotions      = &i18n.Message{ID: "store.postgres.promotion.bulk.insert.app_error", Other: "could not bulk insert promotions"}
	msgGetPromotion              = &i18n.Message{ID: "store.postgres.promotion.get.app_error", Other: "could not get the promotion"}
	msgGetPromotions             = &i18n.Message{ID: "store.postgres.promotion.get.app_error", Other: "could not get the promotion"}
	msgDeletePromotion           = &i18n.Message{ID: "store.postgres.promotion.delete.app_error", Other: "could not delete promotion"}
)

// BulkInsert inserts multiple promotions in the db
func (s PgPromotionStore) BulkInsert(promotions []*model.Promotion) *model.AppErr {
	q := `INSERT INTO public.promotion(promo_code, type, amount, description, starts_at, ends_at) VALUES(:promo_code, :type, :amount, :description, :starts_at, :ends_at) RETURNING promo_code`

	if _, err := s.db.NamedExec(q, promotions); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return model.NewAppErr("PgPromotionStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgUniqueConstraintPromotion, http.StatusInternalServerError, nil)
		}
		return model.NewAppErr("PgPromotionStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertPromotions, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new promotion in the db
func (s PgPromotionStore) Save(promotion *model.Promotion) (*model.Promotion, *model.AppErr) {
	q := `INSERT INTO public.promotion(promo_code, type, amount, description, starts_at, ends_at) VALUES(:promo_code, :type, :amount, :description, :starts_at, :ends_at) RETURNING promo_code`
	if _, err := s.db.NamedExec(q, promotion); err != nil {

		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgPromotionStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgUniqueConstraintPromotion, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgPromotionStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSavePromotion, http.StatusInternalServerError, nil)
	}
	return promotion, nil
}

// Update updates the promotion
func (s PgPromotionStore) Update(code string, promotion *model.Promotion) (*model.Promotion, *model.AppErr) {
	m := map[string]interface{}{
		"code":        code,
		"promo_code":  promotion.PromoCode,
		"type":        promotion.Type,
		"amount":      promotion.Amount,
		"description": promotion.Description,
		"starts_at":   promotion.StartsAt,
		"ends_at":     promotion.EndsAt,
	}

	q := `UPDATE public.promotion SET promo_code=:promo_code, type=:type, amount=:amount, description=:description, starts_at=:starts_at, ends_at=:ends_at WHERE promo_code=:code`
	if _, err := s.db.NamedExec(q, m); err != nil {
		return nil, model.NewAppErr("PgPromotionStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdatePromotion, http.StatusInternalServerError, nil)
	}
	return promotion, nil
}

// Get gets one promotion by id
func (s PgPromotionStore) Get(code string) (*model.Promotion, *model.AppErr) {
	var promotion model.Promotion
	if err := s.db.Get(&promotion, "SELECT * FROM public.promotion WHERE promo_code = $1", code); err != nil {
		return nil, model.NewAppErr("PgPromotionStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetPromotion, http.StatusInternalServerError, nil)
	}
	return &promotion, nil
}

// GetAll returns all promotions
func (s PgPromotionStore) GetAll(limit, offset int) ([]*model.Promotion, *model.AppErr) {
	var promotions = make([]*model.Promotion, 0)
	if err := s.db.Select(&promotions, `SELECT COUNT(*) OVER() AS total_count, * FROM public.promotion LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgPromotionStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetPromotions, http.StatusInternalServerError, nil)
	}

	return promotions, nil
}

// Delete hard deletes the promotion
func (s PgPromotionStore) Delete(code string) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.promotion WHERE promo_code = :code", map[string]interface{}{"code": code}); err != nil {
		return model.NewAppErr("PgPromotionStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeletePromotion, http.StatusInternalServerError, nil)
	}
	return nil
}
