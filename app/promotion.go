package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgPromoNotExists = &i18n.Message{ID: "app.promotion.create_promotio.status.app_error", Other: "promo_code doesn't exist"}
)

// GetPromotionsCount gets all promotions count
func (a *App) GetPromotionsCount() int {
	return a.Srv().Store.Promotion().Count()
}

// CreatePromotion creates the new promotion in the system
func (a *App) CreatePromotion(p *model.Promotion) (*model.Promotion, *model.AppErr) {
	p.PreSave()
	if err := p.Validate(); err != nil {
		return nil, err
	}

	promotion, pErr := a.Srv().Store.Promotion().Save(p)
	if pErr != nil {
		a.Log().Error(pErr.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	return promotion, nil
}

// PatchPromotion patches the promotion
func (a *App) PatchPromotion(code string, patch *model.PromotionPatch) (*model.Promotion, *model.AppErr) {
	old, err := a.Srv().Store.Promotion().Get(code)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	up, err := a.Srv().Store.Promotion().Update(code, old)
	if err != nil {
		return nil, err
	}

	return up, nil
}

// GetPromotion gets the promotion by the promo_code
func (a *App) GetPromotion(code string) (*model.Promotion, *model.AppErr) {
	return a.Srv().Store.Promotion().Get(code)
}

// GetPromotions gets all promotions from the db
func (a *App) GetPromotions(limit, offset int) ([]*model.Promotion, *model.AppErr) {
	return a.Srv().Store.Promotion().GetAll(limit, offset)
}

// DeletePromotion hard deletes the promotion from the db
func (a *App) DeletePromotion(code string) *model.AppErr {
	return a.Srv().Store.Promotion().Delete(code)
}

// IsValidPromotion checks if promo_code is valid
func (a *App) IsValidPromotion(code string) *model.AppErr {
	return a.Srv().Store.Promotion().IsValid(code)
}

// IsUsedPromotion checks if promo_code was used by user
func (a *App) IsUsedPromotion(code string, uid int64) *model.AppErr {
	return a.Srv().Store.Promotion().IsUsed(code, uid)
}

// GetPromotionStatus checks if the promotion is active and not already used by user
func (a *App) GetPromotionStatus(code string, userID int64) *model.AppErr {
	if err := a.IsValidPromotion(code); err != nil {
		return err
	}
	if err := a.IsUsedPromotion(code, userID); err != nil {
		return err
	}
	return nil
}

// CreatePromotionDetail creates the new promotion detail
func (a *App) CreatePromotionDetail(pd *model.PromotionDetail) (*model.PromotionDetail, *model.AppErr) {
	pdetail, pErr := a.Srv().Store.Promotion().InsertDetail(pd)
	if pErr != nil {
		a.Log().Error(pErr.Error(), zlog.Err(pErr))
		return nil, pErr
	}

	return pdetail, nil
}

// DeletePromotions bulk deletes promotions
func (a *App) DeletePromotions(codes []string) *model.AppErr {
	return a.Srv().Store.Promotion().BulkDelete(codes)
}
