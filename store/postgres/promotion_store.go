package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx"
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
	msgUniqueConstraintPromotion       = &i18n.Message{ID: "store.postgres.promotion.save.unique_constraint.app_error", Other: "promotion with given promo_code already exists"}
	msgSavePromotion                   = &i18n.Message{ID: "store.postgres.promotion.save.app_error", Other: "could not save promotion"}
	msgUpdatePromotion                 = &i18n.Message{ID: "store.postgres.promotion.update.app_error", Other: "could not update promotion"}
	msgBulkInsertPromotions            = &i18n.Message{ID: "store.postgres.promotion.bulk.insert.app_error", Other: "could not bulk insert promotions"}
	msgGetPromotion                    = &i18n.Message{ID: "store.postgres.promotion.get.app_error", Other: "could not get the promotion"}
	msgGetPromotions                   = &i18n.Message{ID: "store.postgres.promotion.get.app_error", Other: "could not get the promotion"}
	msgDeletePromotion                 = &i18n.Message{ID: "store.postgres.promotion.delete.app_error", Other: "could not delete promotion"}
	msgBulkDeletePromotions            = &i18n.Message{ID: "store.postgres.promotion.bulk_delete.app_error", Other: "could not bulk delete promotions"}
	msgPromoStatus                     = &i18n.Message{ID: "store.postgres.promotion.status.app_error", Other: "could not get promo_code status"}
	msgPromoCodeUsed                   = &i18n.Message{ID: "store.postgres.promotion.is_used.app_error", Other: "you have already used this promo code"}
	msgPromoCodeInvalid                = &i18n.Message{ID: "store.postgres.promotion.is_valid.app_error", Other: "promo code is invalid or is no longer active"}
	msgInsertPromotionDetail           = &i18n.Message{ID: "store.postgres.promotion.insert_detail.app_error", Other: "could not save promotion detail"}
	msgUniqueConstraintPromotionDetail = &i18n.Message{ID: "store.postgres.promotion.insert_detail.unique_constraint.app_error", Other: "promotion already used by the same user"}
)

// BulkInsert inserts multiple promotions in the db
func (s PgPromotionStore) BulkInsert(promotions []*model.Promotion) *model.AppErr {
	q := `INSERT INTO public.promotion(promo_code, type, amount, description, starts_at, ends_at, created_at, updated_at) VALUES(:promo_code, :type, :amount, :description, :starts_at, :ends_at, :created_at, :updated_at) RETURNING promo_code`

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
	q := `INSERT INTO public.promotion(promo_code, type, amount, description, starts_at, ends_at, created_at, updated_at) VALUES(:promo_code, :type, :amount, :description, :starts_at, :ends_at, :created_at, :updated_at) RETURNING promo_code`
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
		"updated_at":  promotion.UpdatedAt,
	}

	q := `UPDATE public.promotion SET promo_code=:promo_code, type=:type, amount=:amount, description=:description, starts_at=:starts_at, ends_at=:ends_at, updated_at=:updated_at WHERE promo_code=:code`
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
	if err := s.db.Select(&promotions, `SELECT COUNT(*) OVER() AS total_count, * FROM public.promotion ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset); err != nil {
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

// InsertDetail inserts the new promotion in the db
func (s PgPromotionStore) InsertDetail(pdetail *model.PromotionDetail) (*model.PromotionDetail, *model.AppErr) {
	q := `INSERT INTO public.promotion_detail(user_id, promo_code) VALUES(:user_id, :promo_code) RETURNING *`
	if _, err := s.db.NamedExec(q, pdetail); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgPromotionStore.InsertDetail", model.ErrInternal, locale.GetUserLocalizer("en"), msgUniqueConstraintPromotionDetail, http.StatusInternalServerError, nil)
		}

		return nil, model.NewAppErr("PgPromotionStore.InsertDetail", model.ErrInternal, locale.GetUserLocalizer("en"), msgInsertPromotionDetail, http.StatusInternalServerError, nil)
	}
	return pdetail, nil
}

// BulkDelete deletes tags with given ids
func (s PgPromotionStore) BulkDelete(codes []string) *model.AppErr {
	q, args, err := sqlx.In(`DELETE FROM public.promotion WHERE promo_code IN (?)`, codes)
	if err != nil {
		return model.NewAppErr("PgPromotionStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeletePromotions, http.StatusInternalServerError, nil)
	}

	if _, err := s.db.Exec(s.db.Rebind(q), args...); err != nil {
		return model.NewAppErr("PgPromotionStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeletePromotions, http.StatusInternalServerError, nil)
	}

	return nil
}

// IsValid checks if the promo code exists and it is active and valid
func (s PgPromotionStore) IsValid(code string) *model.AppErr {
	var valid bool
	q := `SELECT EXISTS (SELECT 1 FROM promotion p LEFT JOIN promotion_detail pd ON p.promo_code = pd.promo_code WHERE p.promo_code = $1 AND CURRENT_TIMESTAMP BETWEEN p.starts_at AND p.ends_at)`
	if err := s.db.Get(&valid, q, code); err != nil {
		return model.NewAppErr("PgPromotionStore.IsValid", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromoStatus, http.StatusInternalServerError, nil)
	}
	if valid == false {
		return model.NewAppErr("PgPromotionStore.IsValid", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromoCodeInvalid, http.StatusInternalServerError, nil)
	}
	return nil
}

// IsUsed checks if promo code has been used by the user already
func (s PgPromotionStore) IsUsed(code string, userID int64) *model.AppErr {
	var used bool
	q := `SELECT EXISTS (SELECT 1 FROM promotion_detail pd WHERE pd.promo_code = $1 AND pd.user_id = $2)`
	if err := s.db.Get(&used, q, code, userID); err != nil {
		return model.NewAppErr("PgPromotionStore.IsUsed", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromoStatus, http.StatusInternalServerError, nil)
	}
	if used == true {
		return model.NewAppErr("PgPromotionStore.IsUsed", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromoCodeUsed, http.StatusInternalServerError, nil)
	}
	return nil
}
