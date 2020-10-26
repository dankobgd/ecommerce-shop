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

	if data.UseExistingBillingAddress == nil && data.SaveAddress != nil {
		if _, err := a.CreateUserAddress(data.BillingAddress, userID); *data.SaveAddress == true && err != nil {
			return nil, err
		}
	}

	billAddrInfo := &model.Address{}

	if data.UseExistingBillingAddress != nil {
		ua, err := a.GetUserAddress(userID, *data.BillingAddressID)
		if err != nil {
			return nil, err
		}
		billAddrInfo = ua
	} else {
		billAddrInfo = data.BillingAddress
	}

	o := &model.Order{
		UserID: userID,
		Total:  total,
		Status: model.OrderStatusSuccess.String(),
	}

	o.BillingAddressLine1 = billAddrInfo.Line1
	o.BillingAddressLine2 = billAddrInfo.Line2
	o.BillingAddressCity = billAddrInfo.City
	o.BillingAddressCountry = billAddrInfo.Country
	o.BillingAddressState = billAddrInfo.State
	o.BillingAddressZIP = billAddrInfo.ZIP
	o.BillingAddressLatitude = billAddrInfo.Latitude
	o.BillingAddressLongitude = billAddrInfo.Longitude

	o.PreSave()

	if data.SameShippingAsBilling {
		o.ShippingAddressLine1 = billAddrInfo.Line1
		o.ShippingAddressLine2 = billAddrInfo.Line2
		o.ShippingAddressCity = billAddrInfo.City
		o.ShippingAddressCountry = billAddrInfo.Country
		o.ShippingAddressState = billAddrInfo.State
		o.ShippingAddressZIP = billAddrInfo.ZIP
		o.ShippingAddressLatitude = billAddrInfo.Latitude
		o.ShippingAddressLongitude = billAddrInfo.Longitude
	} else {
		o.ShippingAddressLine1 = data.ShippingAddress.Line1
		o.ShippingAddressLine2 = data.ShippingAddress.Line2
		o.ShippingAddressCity = data.ShippingAddress.City
		o.ShippingAddressCountry = data.ShippingAddress.Country
		o.ShippingAddressState = data.ShippingAddress.State
		o.ShippingAddressZIP = data.ShippingAddress.ZIP
		o.ShippingAddressLatitude = data.ShippingAddress.Latitude
		o.ShippingAddressLongitude = data.ShippingAddress.Longitude
	}

	if data.UseExistingBillingAddress == nil || (data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == false) {
		bGeocode, err := a.GetAddressGeocodeResult(data.BillingAddress)
		if err != nil {
			return nil, err
		}

		var sGeocode *model.GeocodingResult
		if data.SameShippingAsBilling {
			sGeocode = bGeocode
		} else {
			sGeocode, err = a.GetAddressGeocodeResult(data.ShippingAddress)
			if err != nil {
				return nil, err
			}
		}

		bLat, _ := strconv.ParseFloat(bGeocode.Lat, 64)
		bLon, _ := strconv.ParseFloat(bGeocode.Lon, 64)
		sLat, _ := strconv.ParseFloat(sGeocode.Lat, 64)
		sLon, _ := strconv.ParseFloat(sGeocode.Lon, 64)

		o.BillingAddressLatitude = &bLat
		o.BillingAddressLongitude = &bLon
		o.ShippingAddressLatitude = &sLat
		o.ShippingAddressLongitude = &sLon
	}

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
