package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/jmoiron/sqlx/types"
	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"
)

var seedCmd = &cobra.Command{
	Use:    "seed",
	Short:  "seed database",
	RunE:   seedDatabaseFn,
	PreRun: loadApp,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seedDatabaseFn(command *cobra.Command, args []string) error {
	if err := seedUsers(); err != nil {
		return err
	}
	if err := seedCategories(); err != nil {
		return err
	}
	if err := seedBrands(); err != nil {
		return err
	}
	if err := seedTags(); err != nil {
		return err
	}
	if err := seedPromotions(); err != nil {
		return err
	}
	if err := seedProducts(); err != nil {
		return err
	}
	if err := seedDiscounts(); err != nil {
		return err
	}
	if err := seedReviews(); err != nil {
		return err
	}
	if err := seedOrders(); err != nil {
		return err
	}
	cmdApp.Log().Info("database seed completed successfully")
	return nil
}

// seedUsers seeds the user tables
func seedUsers() error {
	var users []*model.User

	data, err := ioutil.ReadFile("./data/seeds/users.json")
	if err != nil {
		cmdApp.Log().Error("could not read users.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &users); err != nil {
		cmdApp.Log().Error("could not unmarshal users.json", zlog.String("err: ", err.Error()))
		return err
	}

	// for performance is the same hash for all...
	mockHash := model.HashPassword("Test_123")

	for _, u := range users {
		u.PreSave(true)
		u.Password = mockHash
	}
	if err := cmdApp.Srv().Store.User().BulkInsert(users); err != nil {
		cmdApp.Log().Error("could not seed users", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("users seeded")
	return nil
}

// seedCategories seeds the categories table
func seedCategories() error {
	var categories []*model.Category

	data, err := ioutil.ReadFile("./data/seeds/categories.json")
	if err != nil {
		cmdApp.Log().Error("could not read categories.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &categories); err != nil {
		cmdApp.Log().Error("could not unmarshal categories.json", zlog.String("err: ", err.Error()))
		return err
	}

	for _, u := range categories {
		u.PreSave()
	}
	if err := cmdApp.Srv().Store.Category().BulkInsert(categories); err != nil {
		cmdApp.Log().Error("could not seed categories", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("categories seeded")
	return nil
}

// seedBrands seeds the brand table
func seedBrands() error {
	var brands []*model.Brand

	data, err := ioutil.ReadFile("./data/seeds/brands.json")
	if err != nil {
		cmdApp.Log().Error("could not read brands.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &brands); err != nil {
		cmdApp.Log().Error("could not unmarshal brands.json", zlog.String("err: ", err.Error()))
		return err
	}

	for _, u := range brands {
		u.PreSave()
	}
	if err := cmdApp.Srv().Store.Brand().BulkInsert(brands); err != nil {
		cmdApp.Log().Error("could not seed brands", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("brands seeded")
	return nil
}

// seedTags seeds the tag table
func seedTags() error {
	var tags []*model.Tag

	data, err := ioutil.ReadFile("./data/seeds/tags.json")
	if err != nil {
		cmdApp.Log().Error("could not read tags.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &tags); err != nil {
		cmdApp.Log().Error("could not unmarshal tags.json", zlog.String("err: ", err.Error()))
		return err
	}

	for _, t := range tags {
		t.PreSave()
	}
	if err := cmdApp.Srv().Store.Tag().BulkInsert(tags); err != nil {
		cmdApp.Log().Error("could not seed tags", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("tags seeded")
	return nil
}

// seedReviews seeds the product_review table
func seedReviews() error {
	var reviews []*model.ProductReview

	data, err := ioutil.ReadFile("./data/seeds/reviews.json")
	if err != nil {
		cmdApp.Log().Error("could not read reviews.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &reviews); err != nil {
		cmdApp.Log().Error("could not unmarshal reviews.json", zlog.String("err: ", err.Error()))
		return err
	}

	for _, t := range reviews {
		t.PreSave()
	}
	if err := cmdApp.Srv().Store.ProductReview().BulkInsert(reviews); err != nil {
		cmdApp.Log().Error("could not seed reviews", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("reviews seeded")
	return nil
}

// seedOrders seeds the orders / order_details table
func seedOrders() error {
	var orderDataList []*model.OrderRequestData

	data, err := ioutil.ReadFile("./data/seeds/orders.json")
	if err != nil {
		cmdApp.Log().Error("could not read orders.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &orderDataList); err != nil {
		cmdApp.Log().Error("could not unmarshal orders.json", zlog.String("err: ", err.Error()))
		return err
	}

	for _, orderData := range orderDataList {
		ids := make([]int64, 0)
		for _, x := range orderData.Items {
			ids = append(ids, x.ProductID)
		}

		products, err := cmdApp.GetProductsbyIDS(ids)
		if err != nil {
			cmdApp.Log().Error("could not get products by id", zlog.String("err: ", err.Error()))
			return err
		}

		var total int
		for i, p := range products {
			total += p.Price * orderData.Items[i].Quantity
		}

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		params := &stripe.PaymentMethodParams{
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String("4242424242424242"),
				ExpMonth: stripe.String("2"),
				ExpYear:  stripe.String("2022"),
				CVC:      stripe.String("314"),
			},
			Type: stripe.String("card"),
		}
		pm, e := paymentmethod.New(params)
		if e != nil {
			cmdApp.Log().Error("stripe payment method error", zlog.String("err: ", e.Error()))
			return e
		}

		rand.Seed(time.Now().UnixNano())
		userID := rand.Intn(1000) + 1
		user, _ := cmdApp.GetUserByID(int64(userID))

		o := &model.Order{
			UserID:                   int64(userID),
			Subtotal:                 total,
			Total:                    total,
			Status:                   model.OrderStatusSuccess.String(),
			PaymentMethodID:          pm.ID,
			BillingAddressLine1:      orderData.BillingAddress.Line1,
			BillingAddressLine2:      orderData.BillingAddress.Line2,
			BillingAddressCity:       orderData.BillingAddress.City,
			BillingAddressCountry:    orderData.BillingAddress.Country,
			BillingAddressState:      orderData.BillingAddress.State,
			BillingAddressZIP:        orderData.BillingAddress.ZIP,
			BillingAddressLatitude:   orderData.BillingAddress.Latitude,
			BillingAddressLongitude:  orderData.BillingAddress.Longitude,
			ShippingAddressLine1:     orderData.ShippingAddress.Line1,
			ShippingAddressLine2:     orderData.ShippingAddress.Line2,
			ShippingAddressCity:      orderData.ShippingAddress.City,
			ShippingAddressCountry:   orderData.ShippingAddress.Country,
			ShippingAddressState:     orderData.ShippingAddress.State,
			ShippingAddressZIP:       orderData.ShippingAddress.ZIP,
			ShippingAddressLatitude:  orderData.ShippingAddress.Latitude,
			ShippingAddressLongitude: orderData.ShippingAddress.Longitude,
		}

		pi, cErr := cmdApp.PaymentProvider().Charge(pm.ID, o, user, uint64(total), "usd")
		if cErr != nil {
			cmdApp.Log().Error("stripe charge err", zlog.String("err: ", cErr.Error()))
			return cErr
		}

		o.PaymentIntentID = pi.ID
		o.ReceiptURL = pi.Charges.Data[0].ReceiptURL

		o.PreSave()
		order, err := cmdApp.Srv().Store.Order().Save(o)
		if err != nil {
			cmdApp.Log().Error("seed save order err", zlog.String("err: ", err.Error()))
			return err
		}

		orderDetails := make([]*model.OrderDetail, 0)
		for i, p := range products {
			detail := &model.OrderDetail{
				OrderID:      order.ID,
				ProductID:    p.ID,
				Quantity:     orderData.Items[i].Quantity,
				HistoryPrice: p.Price,
				HistorySKU:   p.SKU,
			}
			orderDetails = append(orderDetails, detail)
		}

		if err := cmdApp.InsertOrderDetails(orderDetails); err != nil {
			cmdApp.Log().Error("seed insert order details err", zlog.String("err: ", err.Error()))
			return err
		}
	}

	cmdApp.Log().Info("orders seeded")
	return nil
}

// seedPromotions seeds the product_promotions table
func seedPromotions() error {
	var promotions []*model.Promotion

	data, err := ioutil.ReadFile("./data/seeds/promotions.json")
	if err != nil {
		cmdApp.Log().Error("could not read promotions.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &promotions); err != nil {
		cmdApp.Log().Error("could not unmarshal promotions.json", zlog.String("err: ", err.Error()))
		return err
	}

	if err := cmdApp.Srv().Store.Promotion().BulkInsert(promotions); err != nil {
		cmdApp.Log().Error("could not seed promotions", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("promotions seeded")
	return nil
}

// seedDiscounts seeds price discounts
func seedDiscounts() error {
	var discounts []*model.ProductPricing

	data, err := ioutil.ReadFile("./data/seeds/discounts.json")
	if err != nil {
		cmdApp.Log().Error("could not read discounts.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &discounts); err != nil {
		cmdApp.Log().Error("could not unmarshal discounts.json", zlog.String("err: ", err.Error()))
		return err
	}

	ids := make([]int64, 0)
	for _, x := range discounts {
		ids = append(ids, x.ProductID)
	}

	discountPercents := []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
	products, e := cmdApp.GetProductsbyIDS(ids)
	if e != nil {
		cmdApp.Log().Error("could not get products for discount ids.json", zlog.String("err: ", e.Error()))
		return e
	}

	rand.Seed(time.Now().UnixNano())

	findIndex := func(arr []*model.Product, x int) int {
		for i, elm := range arr {
			if int64(x) == elm.ID {
				return i
			}
		}
		return -1
	}

	for i := 0; i < len(discounts); i++ {
		pcv := discountPercents[rand.Intn(len(discountPercents))]

		if idx := findIndex(products, int(discounts[i].ProductID)); idx != -1 {
			dprice := float64(products[idx].Price) - float64(pcv)/100*float64(products[idx].Price)
			discounts[i].Price = int(math.Round(dprice*100) / 100)
			discounts[i].OriginalPrice = products[i].Price
		}
	}

	for _, d := range discounts {
		if _, err := cmdApp.AddProductPricing(d); err != nil {
			cmdApp.Log().Error("could not add pricing for seed", zlog.String("err: ", e.Error()))
			return errors.New("could not add pricing")
		}
	}

	cmdApp.Log().Info("discounts seeded")
	return nil
}

type productSeed struct {
	ID            int64          `json:"id"`
	BrandID       int64          `json:"brand_id"`
	CategoryID    int64          `json:"category_id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	ImageURL      string         `json:"image_url"`
	ImagePublicID string         `json:"image_public_id"`
	Description   string         `json:"description"`
	Price         int            `json:"price"`
	OriginalPrice int            `json:"original_price"`
	InStock       bool           `json:"in_stock"`
	SKU           string         `json:"sku"`
	IsFeatured    bool           `json:"is_featured"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Tags          []int64        `json:"tags"`
	Images        []string       `json:"images"`
	Properties    types.JSONText `json:"properties"`
}

// seedProducts seeds the product table
func seedProducts() error {
	var ps []*productSeed

	data, err := ioutil.ReadFile("./data/seeds/products.json")
	if err != nil {
		cmdApp.Log().Error("could not read products.json seed", zlog.String("err: ", err.Error()))
		return err
	}
	if err := json.Unmarshal(data, &ps); err != nil {
		cmdApp.Log().Error("could not unmarshal products.json", zlog.String("err: ", err.Error()))
		return err
	}

	products := make([]*model.Product, 0)
	pricings := make([]*model.ProductPricing, 0)
	productTags := make([]*model.ProductTag, 0)
	productImgs := make([]*model.ProductImage, 0)

	for i, x := range ps {
		p := &model.Product{
			BrandID:       x.BrandID,
			CategoryID:    x.CategoryID,
			Name:          x.Name,
			Slug:          x.Slug,
			ImageURL:      x.ImageURL,
			ImagePublicID: x.ImagePublicID,
			Description:   x.Description,
			InStock:       x.InStock,
			SKU:           x.SKU,
			IsFeatured:    x.IsFeatured,
			Properties:    &x.Properties,
			ProductPricing: &model.ProductPricing{
				Price:         x.Price,
				OriginalPrice: x.Price,
			},
		}

		for _, tagID := range x.Tags {
			productTags = append(productTags, &model.ProductTag{
				TagID:     model.NewInt64(tagID),
				ProductID: model.NewInt64(int64(i + 1)),
			})
		}
		for _, img := range x.Images {
			now := time.Now()
			productImgs = append(productImgs, &model.ProductImage{
				ProductID: model.NewInt64(int64(i + 1)),
				URL:       model.NewString(img),
				PublicID:  model.NewString(""),
				CreatedAt: &now,
				UpdatedAt: &now,
			})
		}

		price := &model.ProductPricing{
			ProductID:     int64(i + 1),
			Price:         x.Price,
			OriginalPrice: x.Price,
			SaleStarts:    time.Now(),
			SaleEnds:      model.FutureSaleEndsTime,
		}

		products = append(products, p)
		pricings = append(pricings, price)
	}

	for _, p := range products {
		p.PreSave()
	}
	if err := cmdApp.Srv().Store.Product().BulkInsert(products); err != nil {
		cmdApp.Log().Error("could not seed products", zlog.String("err: ", err.Message))
		return err
	}

	if err := cmdApp.Srv().Store.Product().InsertPricingBulk(pricings); err != nil {
		cmdApp.Log().Error("could not seed product pricings", zlog.String("err: ", err.Message))
		return err
	}

	if err := cmdApp.Srv().Store.ProductTag().BulkInsert(productTags); err != nil {
		cmdApp.Log().Error("could not seed bulk insert product tags", zlog.String("err: ", err.Message))
		return err
	}

	for _, img := range productImgs {
		img.PreSave()
	}
	if err := cmdApp.Srv().Store.ProductImage().BulkInsert(productImgs); err != nil {
		cmdApp.Log().Error("could not seed bulk insert product images", zlog.String("err: ", err.Message))
		return err
	}

	cmdApp.Log().Info("products seeded")
	return nil
}
