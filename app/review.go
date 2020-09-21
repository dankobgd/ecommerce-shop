package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
)

// CreateReview creates the new review in the system
func (a *App) CreateReview(rev *model.Review) (*model.Review, *model.AppErr) {
	rev.PreSave()
	if err := rev.Validate(); err != nil {
		return nil, err
	}

	review, rErr := a.Srv().Store.Review().Save(rev)
	if rErr != nil {
		a.Log().Error(rErr.Error(), zlog.Err(rErr))
		return nil, rErr
	}

	return review, nil
}

// PatchReview patches the review
func (a *App) PatchReview(rid int64, patch *model.ReviewPatch) (*model.Review, *model.AppErr) {
	old, err := a.Srv().Store.Review().Get(rid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	urev, err := a.Srv().Store.Review().Update(rid, old)
	if err != nil {
		return nil, err
	}

	return urev, nil
}

// GetReview gets the review by the id
func (a *App) GetReview(tid int64) (*model.Review, *model.AppErr) {
	return a.Srv().Store.Review().Get(tid)
}

// GetReviews gets all reviews from the db
func (a *App) GetReviews(limit, offset int) ([]*model.Review, *model.AppErr) {
	return a.Srv().Store.Review().GetAll(limit, offset)
}

// DeleteReview hard deletes the review from the db
func (a *App) DeleteReview(tid int64) *model.AppErr {
	return a.Srv().Store.Review().Delete(tid)
}
