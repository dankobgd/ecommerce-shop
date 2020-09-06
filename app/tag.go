package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
)

// CreateTag creates the new tag in the system
func (a *App) CreateTag(t *model.Tag) (*model.Tag, *model.AppErr) {
	t.PreSave()
	if err := t.Validate(); err != nil {
		return nil, err
	}

	tag, tErr := a.Srv().Store.Tag().Save(t)
	if tErr != nil {
		a.Log().Error(tErr.Error(), zlog.Err(tErr))
		return nil, tErr
	}

	return tag, nil
}

// PatchTag patches the tag
func (a *App) PatchTag(tid int64, patch *model.TagPatch) (*model.Tag, *model.AppErr) {
	old, err := a.Srv().Store.Tag().Get(tid)
	if err != nil {
		return nil, err
	}

	old.Patch(patch)
	old.PreUpdate()
	utag, err := a.Srv().Store.Tag().Update(tid, old)
	if err != nil {
		return nil, err
	}

	return utag, nil
}

// GetTag gets the tag by the id
func (a *App) GetTag(tid int64) (*model.Tag, *model.AppErr) {
	return a.Srv().Store.Tag().Get(tid)
}

// GetTags gets all tags from the db
func (a *App) GetTags(limit, offset int) ([]*model.Tag, *model.AppErr) {
	return a.Srv().Store.Tag().GetAll(limit, offset)
}

// DeleteTag hard deletes the tag from the db
func (a *App) DeleteTag(tid int64) *model.AppErr {
	return a.Srv().Store.Tag().Delete(tid)
}
