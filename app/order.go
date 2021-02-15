package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/jung-kurt/gofpdf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stripe/stripe-go"
)

var (
	msgGetAddressGeocodeResult = &i18n.Message{ID: "app.order.get_address_geocode_result.app_error", Other: "could not get geocoding result on given address"}
	msgCreatePDF               = &i18n.Message{ID: "app.order.details_pdf.app_error", Other: "could not create order details pdf"}
)

// CreateOrder creates the new order
func (a *App) CreateOrder(userID int64, data *model.OrderRequestData) (*model.Order, *model.AppErr) {
	// validate order request data
	if err := data.Validate(); err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, x := range data.Items {
		ids = append(ids, x.ProductID)
	}

	// get products from given data product ids
	products, err := a.GetProductsbyIDS(ids)
	if err != nil {
		return nil, err
	}

	// get authed user
	user, err := a.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// calc subtotal price (price before discount, possible taxes etc...)
	subtotal := 0
	for i, p := range products {
		subtotal += p.Price * data.Items[i].Quantity
	}

	// calc total price (after possible discount)
	total := 0
	if data.PromoCode == nil {
		total = subtotal
	}

	var promoCodeName *string
	var promoCodeType *string
	var promoCodeAmount *int

	if data.PromoCode != nil && *data.PromoCode != "" {
		if err := a.GetPromotionStatus(*data.PromoCode, userID); err != nil {
			return nil, err
		}

		promo, err := a.GetPromotion(*data.PromoCode)
		if err != nil {
			return nil, err
		}

		promoCodeName = &promo.PromoCode
		promoCodeType = &promo.Type
		promoCodeAmount = &promo.Amount

		if promo.Type == "percentage" {
			t := float64(subtotal) - float64(promo.Amount)/100*float64(subtotal)
			total = int(math.Round(t*100) / 100)
		}

		if promo.Type == "fixed" {
			t := (subtotal - promo.Amount)
			if t < 0 {
				t = 0
			}
			total = t
		}
	}

	billAddrInfo := &model.Address{}

	if data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == true {
		if data.BillingAddressID != nil {
			ua, err := a.GetUserAddress(userID, *data.BillingAddressID)
			if err != nil {
				return nil, err
			}
			billAddrInfo = ua
		}
	} else {
		billAddrInfo = data.BillingAddress
	}

	o := &model.Order{
		UserID:          userID,
		Subtotal:        subtotal,
		Total:           total,
		Status:          model.OrderStatusSuccess.String(),
		PaymentMethodID: data.PaymentMethodID,
		PromoCode:       promoCodeName,
		PromoCodeType:   promoCodeType,
		PromoCodeAmount: promoCodeAmount,
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

	if data.SameShippingAsBilling != nil && *data.SameShippingAsBilling == true {
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
		if data.SameShippingAsBilling != nil && *data.SameShippingAsBilling == true {
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

	pi, cErr := a.PaymentProvider().Charge(data.PaymentMethodID, o, user, uint64(o.Total), "usd")
	if cErr != nil {
		if stripeErr, ok := cErr.(*stripe.Error); ok {
			if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
				dc := ""
				if (string(cardErr.DeclineCode)) != "" {
					dc = "Decline code: " + string(cardErr.DeclineCode)
				}

				return nil, model.NewAppErr("CreateOrder", model.ErrInternal, locale.GetUserLocalizer("en"), &i18n.Message{ID: "app.order.create_order.app_error", Other: fmt.Sprintf("%s\n%s", stripeErr.Msg, dc)}, http.StatusInternalServerError, nil)
			}
			return nil, model.NewAppErr("CreateOrder", model.ErrInternal, locale.GetUserLocalizer("en"), &i18n.Message{ID: "app.order.create_order.app_error", Other: stripeErr.Msg}, http.StatusInternalServerError, nil)

		}
		return nil, model.NewAppErr("CreateOrder", model.ErrInternal, locale.GetUserLocalizer("en"), &i18n.Message{ID: "app.order.create_order.app_error", Other: "could not charge the card"}, http.StatusInternalServerError, nil)
	}

	o.PaymentIntentID = pi.ID
	o.ReceiptURL = pi.Charges.Data[0].ReceiptURL

	// save actual order
	order, err := a.Srv().Store.Order().Save(o)
	if err != nil {
		return nil, err
	}

	orderDetails := make([]*model.OrderDetail, 0)
	for i, p := range products {
		detail := &model.OrderDetail{
			OrderID:      order.ID,
			ProductID:    p.ID,
			Quantity:     data.Items[i].Quantity,
			HistoryPrice: p.Price,
			HistorySKU:   p.SKU,
		}
		orderDetails = append(orderDetails, detail)
	}

	if err := a.InsertOrderDetails(orderDetails); err != nil {
		return nil, err
	}

	if data.PromoCode != nil && *data.PromoCode != "" {
		// insert promo detail to mark the promo_code as used by the specific user
		pd := &model.PromotionDetail{UserID: userID, PromoCode: *data.PromoCode}
		if _, err := a.CreatePromotionDetail(pd); err != nil {
			return nil, err
		}
	}

	defer func() {
		if (data.UseExistingBillingAddress == nil || data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == false) && (data.SaveAddress != nil && *data.SaveAddress == true) {
			if _, err := a.CreateUserAddress(data.BillingAddress, userID); err != nil {
				a.Log().Error(err.Error(), zlog.Err(err))
			}
		}
	}()

	return order, nil
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

// InsertOrderDetails inserts new order details
func (a *App) InsertOrderDetails(items []*model.OrderDetail) *model.AppErr {
	return a.Srv().Store.OrderDetail().BulkInsert(items)
}

// InsertOrderDetail inserts new order detail
func (a *App) InsertOrderDetail(item *model.OrderDetail) (*model.OrderDetail, *model.AppErr) {
	return a.Srv().Store.OrderDetail().Save(item)
}

// GetOrderDetails gets the order details for the order id
func (a *App) GetOrderDetails(orderID int64) ([]*model.OrderInfo, *model.AppErr) {
	return a.Srv().Store.OrderDetail().GetAll(orderID)
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

const (
	logoH   = 94.0
	xIndent = 40.0
)

// GenerateOrderDetailsPDF creates the pdf
func (a *App) GenerateOrderDetailsPDF(o *model.Order, details []*model.OrderInfo, user *model.User) (bytes.Buffer, *model.AppErr) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	w, h := pdf.GetPageSize()
	pdf.AddPage()

	// Top skewed header background
	pdf.SetFillColor(103, 60, 79)
	pdf.Polygon([]gofpdf.PointType{
		{X: 0, Y: 0},
		{X: w, Y: 0},
		{X: w, Y: logoH},
		{X: 0, Y: logoH * 0.9},
	}, "F")
	pdf.Polygon([]gofpdf.PointType{
		{X: 0, Y: h},
		{X: 0, Y: h - (logoH * 0.2)},
		{X: w, Y: h - (logoH * 0.1)},
		{X: w, Y: h},
	}, "F")

	// header invoice
	pdf.SetFont("arial", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt := pdf.GetFontSize()
	pdf.Text(xIndent, logoH-(logoH/2.0)+lineHt/3.1, "INVOICE")

	// logo
	pdf.ImageOptions("./static/images/invoice.jpg", 272.0, 0+(logoH-(logoH/1.5))/2.0, 0, logoH/1.5, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	// user details
	userDetails := make([]string, 0)
	userDetails = append(userDetails, user.FirstName, user.LastName, user.Email)
	formattedUserDetails := ""
	for _, x := range userDetails {
		formattedUserDetails += fmt.Sprintf("%s\n", x)
	}

	pdf.SetFont("arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-2.0*124.0+60, (logoH-(lineHt*1.5*3.0))/2.0)
	pdf.MultiCell(200.0, lineHt*1.5, formattedUserDetails, gofpdf.BorderNone, gofpdf.AlignRight, false)

	// addr details
	billAddrDetails := make([]string, 0)
	billAddrDetails = append(billAddrDetails, o.BillingAddressLine1)
	if o.BillingAddressLine2 != nil && *o.BillingAddressLine2 != "" {
		billAddrDetails = append(billAddrDetails, *o.BillingAddressLine2)
	}
	billAddrDetails = append(billAddrDetails, o.BillingAddressCity, o.BillingAddressCountry)
	if o.BillingAddressZIP != nil && *o.BillingAddressZIP != "" {
		billAddrDetails = append(billAddrDetails, *o.BillingAddressZIP)
	}

	shipAddrDetails := make([]string, 0)
	shipAddrDetails = append(shipAddrDetails, o.ShippingAddressLine1)
	if o.ShippingAddressLine2 != nil && *o.ShippingAddressLine2 != "" {
		shipAddrDetails = append(shipAddrDetails, *o.ShippingAddressLine2)
	}
	shipAddrDetails = append(shipAddrDetails, o.ShippingAddressCity, o.ShippingAddressCountry)
	if o.BillingAddressZIP != nil && *o.BillingAddressZIP != "" {
		shipAddrDetails = append(shipAddrDetails, *o.BillingAddressZIP)
	}

	// Summary - billing / shipping address
	_, sy := summaryBlock(pdf, xIndent, logoH+lineHt*2.0, "Billing Address", billAddrDetails...)
	summaryBlock(pdf, xIndent+175.0, logoH+lineHt*2.0, "Shipping Address", shipAddrDetails...)

	// Summary - Invoice Total
	x, y := w-xIndent-124.0, logoH+lineHt*2.25
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.CellFormat(124.0, lineHt, "Invoice Total", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x+2.0, y+lineHt*1.5
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 42)
	_, lineHt = pdf.GetFontSize()
	alpha := 58
	pdf.SetTextColor(72+alpha, 42+alpha, 55+alpha)
	pdf.CellFormat(124.0, lineHt, toUSD(o.Total), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x-2.0, y+lineHt*1.25

	if sy > y {
		y = sy
	}
	x, y = xIndent-20.0, y+30.0
	pdf.Rect(x, y, w-(xIndent*2.0)+40.0, 3.0, "F")

	// Line Items - headers
	pdf.SetFont("times", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	x, y = xIndent-2.0, y+lineHt
	pdf.MoveTo(x, y)
	pdf.CellFormat(w/2.65+1.5, lineHt, "Name", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHt, "Price", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, "Quantity", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(119.5, lineHt, "Summed Price", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	// Line Items - real data
	y = y + lineHt

	for _, dtl := range details {
		x, y = lineItem(pdf, x, y, dtl)
	}

	// Subtotal etc
	x, y = w/1.75, y+lineHt*2.25
	x, y = trailerLine(pdf, x, y, "Subtotal", toUSD(o.Subtotal))

	if o.PromoCode != nil && *o.PromoCode != "" {
		promoStr := fmt.Sprintf("-%v", toUSD(*o.PromoCodeAmount))
		if *o.PromoCodeType == "percentage" {
			promoStr = fmt.Sprintf("-%v%%", *o.PromoCodeAmount)
		}
		x, y = trailerLine(pdf, x, y, "Promo Code", promoStr)
	}

	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(x+10.0, y, x+220.0, y)
	y = y + lineHt*0.5
	x, y = trailerLine(pdf, x, y, "Total Charge", toUSD(o.Total))

	var buf bytes.Buffer

	if err := pdf.Output(&buf); err != nil {
		return buf, model.NewAppErr("PDF", model.ErrInternal, locale.GetUserLocalizer("en"), msgCreatePDF, http.StatusInternalServerError, nil)
	}

	return buf, nil
}

func trailerLine(pdf *gofpdf.Fpdf, x, y float64, label string, formattedAmount string) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, label, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(119.5, lineHt, formattedAmount, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	y = y + lineHt*1.5
	return origX, y
}

func toUSD(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}
	return fmt.Sprintf("$%d.%s", cents/100, centsStr)
}

func lineItem(pdf *gofpdf.Fpdf, x, y float64, item *model.OrderInfo) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(50, 50, 50)
	pdf.MoveTo(x, y)
	x, y = xIndent-2.0, y+lineHt*.75
	pdf.MoveTo(x, y)
	pdf.MultiCell(w/2.65+1.5, lineHt, item.Product.Name, gofpdf.BorderNone, gofpdf.AlignLeft, false)
	tmp := pdf.SplitLines([]byte(item.Product.Name), w/2.65+1.5)
	maxY := y + float64(len(tmp)-1)*lineHt
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHt, toUSD(item.HistoryPrice), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, fmt.Sprintf("%d", item.Quantity), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(119.5, lineHt, toUSD(item.HistoryPrice*item.Quantity), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	if maxY > y {
		y = maxY
	}
	y = y + lineHt*1.75
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(xIndent-10.0, y, w-xIndent+10.0, y)
	return origX, y
}

func summaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	pdf.SetFont("times", "", 14)
	pdf.SetTextColor(180, 180, 180)
	_, lineHt := pdf.GetFontSize()
	y = y + lineHt
	pdf.Text(x, y, title)
	y = y + lineHt*.25
	pdf.SetTextColor(50, 50, 50)
	for _, str := range data {
		y = y + lineHt*1.25
		pdf.Text(x, y, str)
	}
	return x, y
}
