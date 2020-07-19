package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgGetAddressGeocodeResult = &i18n.Message{ID: "app.order.get_address_geocode_result.app_error", Other: "could not get geocoding result on given address"}
)

// CreateOrder creates the new order
func (a *App) CreateOrder(o *model.Order, shipAddr *model.Address, billAddr *model.Address) (*model.Order, *model.AppErr) {
	o.PreSave()
	return a.Srv().Store.Order().Save(o, shipAddr, billAddr)
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

// GetAddressGeocodeResult gets the lat, lng etc...
func (a *App) GetAddressGeocodeResult(addr *model.Address) (*model.GeocodingResult, *model.AppErr) {
	client := &http.Client{Timeout: 5 * time.Second}
	baseURL, _ := url.Parse("https://us1.locationiq.com/v1/search.php")

	q := baseURL.Query()
	q.Set("format", "json")
	q.Set("key", a.cfg.GeocodingSettings.ApiKey)
	q.Set("city", addr.City)
	q.Set("country", addr.Country)
	if addr.ZIP != nil {
		q.Set("postalcode", *addr.ZIP)
	}
	baseURL.RawQuery = q.Encode()

	var result model.GeocodingResultList
	resp, e := client.Get(baseURL.String())
	if e != nil {
		return nil, model.NewAppErr("GetAddressGeocodeResult", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetAddressGeocodeResult, http.StatusInternalServerError, nil)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&result)

	// maybe return the one with highest importance points...
	return result[0], nil
}
