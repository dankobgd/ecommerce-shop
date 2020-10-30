package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
)

// CreatePromotion creates the new promotion in the system
func (a *App) CreatePromotion(p *model.Promotion) (*model.Promotion, *model.AppErr) {
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
