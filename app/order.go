package app

import (
	"github.com/dankobgd/ecommerce-shop/model"
)

// CreateOrder creates the new order
func (a *App) CreateOrder(o *model.Order) (*model.Order, *model.AppErr) {
	o.PreSave()
	return a.Srv().Store.Order().Save(o)
}

// GetOrder gets the order by id
func (a *App) GetOrder(id int64) (*model.Order, *model.AppErr) {
	return a.Srv().Store.Order().Get(id)
}

// UpdateOrder updates the order
func (a *App) UpdateOrder(id int64, o *model.Order) (*model.Order, *model.AppErr) {
	return a.Srv().Store.Order().Update(id, o)
}

// CreateOrderDetails inserts new order details
func (a *App) CreateOrderDetails(items []*model.OrderDetail) *model.AppErr {
	return a.Srv().Store.OrderDetail().BulkInsert(items)
}

// GetOrderDetail gets the order detail by id
func (a *App) GetOrderDetail(id int64) (*model.OrderDetail, *model.AppErr) {
	return a.Srv().Store.OrderDetail().Get(id)
}
