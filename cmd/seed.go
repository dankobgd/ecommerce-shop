package cmd

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/jmoiron/sqlx/types"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:    "seed",
	Short:  "seed database",
	RunE:   seedDatabaseFn,
	PreRun: loadApp,
}

var seedUsersCmd = &cobra.Command{
	Use:    "users",
	Short:  "seed users",
	RunE:   seedUsersFn,
	PreRun: loadApp,
}

var seedProductsCmd = &cobra.Command{
	Use:    "products",
	Short:  "seed products",
	RunE:   seedProductsFn,
	PreRun: loadApp,
}

var seedCategoriesCmd = &cobra.Command{
	Use:    "categories",
	Short:  "seed categories",
	RunE:   seedCategoriesFn,
	PreRun: loadApp,
}

var seedBrandsCmd = &cobra.Command{
	Use:    "brands",
	Short:  "seed brands",
	RunE:   seedBrandsFn,
	PreRun: loadApp,
}

var seedTagsCmd = &cobra.Command{
	Use:    "tags",
	Short:  "seed tags",
	RunE:   seedTagsFn,
	PreRun: loadApp,
}

var seedReviewsCmd = &cobra.Command{
	Use:    "reviews",
	Short:  "seed reviews",
	RunE:   seedReviewsFn,
	PreRun: loadApp,
}

func init() {
	seedCmd.AddCommand(seedUsersCmd, seedProductsCmd, seedCategoriesCmd, seedBrandsCmd, seedTagsCmd)
	rootCmd.AddCommand(seedCmd)
}

func seedUsersFn(command *cobra.Command, args []string) error {
	return seedUsers()
}

func seedCategoriesFn(command *cobra.Command, args []string) error {
	return seedCategories()
}

func seedBrandsFn(command *cobra.Command, args []string) error {
	return seedBrands()
}

func seedTagsFn(command *cobra.Command, args []string) error {
	return seedTags()
}

func seedReviewsFn(command *cobra.Command, args []string) error {
	return seedReviews()
}

func seedPromotionsFn(command *cobra.Command, args []string) error {
	return seedPromotions()
}

func seedProductsFn(command *cobra.Command, args []string) error {
	return seedProducts()
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
	if err := seedReviews(); err != nil {
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

// seedCategories populates the categories table
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

// seedBrands populates the brand table
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

// seedTags populates the tag table
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

// seedReviews populates the product_review table
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

// seedPromotions populates the product_promotions table
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

// seedProducts populates the product table
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
