package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
)

// CreateProduct creates the new product in the system
func (a *App) CreateProduct(p *model.Product) (*model.Product, *model.AppErr) {
	p.PreSave()
	if err := p.Validate(); err != nil {
		return nil, err
	}

	product, err := a.Srv().Store.Product().Save(p)
	if err != nil {
		a.log.Error(err.Error(), zlog.Err(err))
		return nil, err
	}

	return product, nil
}
