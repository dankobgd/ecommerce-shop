package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgGetAddressGeocodeResult = &i18n.Message{ID: "app.order.get_address_geocode_result.app_error", Other: "could not get geocoding result on given address"}
)

// CreateOrder creates the new order
func (a *App) CreateOrder(userID int64, data *model.OrderRequestData, total int) (*model.Order, *model.AppErr) {
	user, err := a.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	o := &model.Order{
		UserID:                userID,
		Total:                 total,
		Status:                model.OrderStatusSuccess.String(),
		BillingAddressLine1:   data.BillingAddress.Line1,
		BillingAddressLine2:   data.BillingAddress.Line2,
		BillingAddressCity:    data.BillingAddress.City,
		BillingAddressCountry: data.BillingAddress.Country,
		BillingAddressState:   data.BillingAddress.State,
		BillingAddressZIP:     data.BillingAddress.ZIP,
	}
	o.PreSave()

	if data.SameShippingAsBilling {
		o.ShippingAddressLine1 = data.BillingAddress.Line1
		o.ShippingAddressLine2 = data.BillingAddress.Line2
		o.ShippingAddressCity = data.BillingAddress.City
		o.ShippingAddressCountry = data.BillingAddress.Country
		o.ShippingAddressState = data.BillingAddress.State
		o.ShippingAddressZIP = data.BillingAddress.ZIP
	}

	billingAddressGeocode, err := a.GetAddressGeocodeResult(data.BillingAddress)
	if err != nil {
		return nil, err
	}

	var shippingAddressGeocode *model.GeocodingResult

	if data.SameShippingAsBilling {
		shippingAddressGeocode = billingAddressGeocode
	} else {
		shippingAddressGeocode, err = a.GetAddressGeocodeResult(data.ShippingAddress)
		if err != nil {
			return nil, err
		}
	}

	sLat, _ := strconv.ParseFloat(shippingAddressGeocode.Lat, 64)
	sLon, _ := strconv.ParseFloat(shippingAddressGeocode.Lon, 64)
	o.ShippingAddressLatitude = &sLat
	o.ShippingAddressLongitude = &sLon

	bLat, _ := strconv.ParseFloat(billingAddressGeocode.Lat, 64)
	bLon, _ := strconv.ParseFloat(billingAddressGeocode.Lon, 64)
	o.BillingAddressLatitude = &bLat
	o.BillingAddressLongitude = &bLon

	_, cErr := a.PaymentProvider().Charge(data.PaymentMethodID, o, user, uint64(o.Total), "usd")
	if cErr != nil {
		return nil, model.NewAppErr("CreateOrder", model.ErrInternal, locale.GetUserLocalizer("en"), &i18n.Message{ID: "app.order.create_order.app_error", Other: "could not charge card: " + cErr.Error()}, http.StatusInternalServerError, nil)
	}

	return a.Srv().Store.Order().Save(o)
}

// GetOrder gets the order by id
func (a *App) GetOrder(id int64) (*model.Order, *model.AppErr) {
	return a.Srv().Store.Order().Get(id)
}

// GetOrders gets all orders
func (a *App) GetOrders(limit, offset int) ([]*model.Order, *model.AppErr) {
	return a.Srv().Store.Order().GetAll(limit, offset)
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
	q.Set("key", a.cfg.GeocodingSettings.APIKey)
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
