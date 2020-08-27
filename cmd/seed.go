package cmd

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/zlog"
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

func init() {
	seedCmd.AddCommand(seedUsersCmd, seedProductsCmd, seedCategoriesCmd, seedBrandsCmd, seedTagsCmd)
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
	if err := seedProducts(); err != nil {
		return err
	}
	cmdApp.Log().Info("database seed completed successfully")
	return nil
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

func seedProductsFn(command *cobra.Command, args []string) error {
	return seedProducts()
}

// seedUsers seeds the user tables
func seedUsers() error {
	users := parseUsers()
	for _, u := range users {
		u.PreSave()
	}
	if err := cmdApp.Srv().Store.User().BulkInsert(users); err != nil {
		cmdApp.Log().Error("could not seed users", zlog.String("err: ", err.Message))
		return err
	}
	cmdApp.Log().Info("users seeded")
	return nil
}

// seedProducts populates the product table
func seedProducts() error {
	productData := parseProducts()

	for _, d := range productData {
		d.P.PreSave()
		newProd, err := cmdApp.Srv().Store.Product().Save(d.P)
		if err != nil {
			cmdApp.Log().Error("could not seed save product", zlog.String("err: ", err.Message))
		}

		for _, tag := range d.Tags {
			tag.ProductID = &newProd.ID
		}
		if _, err := cmdApp.Srv().Store.ProductTag().BulkInsert(d.Tags); err != nil {
			cmdApp.Log().Error("could not seed bulk insert tags", zlog.String("err: ", err.Message))
		}

		for _, img := range d.Imgs {
			img.ProductID = &newProd.ID
			img.PreSave()
		}
		if _, err := cmdApp.Srv().Store.ProductImage().BulkInsert(d.Imgs); err != nil {
			cmdApp.Log().Error("could not seed bulk insert images", zlog.String("err: ", err.Message))
		}
	}

	cmdApp.Log().Info("products seeded")
	return nil
}

// seedCategories populates the categories table
func seedCategories() error {
	categories := parseCategories()
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
	brands := parseBrands()
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
	tags := parseTags()
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

// parseUsers parses csv lines and creates the list of users
func parseUsers() []*model.User {
	records, err := readCSVFile("./data/seeds/users.csv")
	if err != nil {
		log.Fatalf("error parsing users from CSV: %v", err)
	}
	userList := make([]*model.User, 0)

	for _, line := range records {
		u := &model.User{}
		u.Username = line[0]
		u.Email = line[1]
		u.Password = line[2]
		u.Role = line[3]
		userList = append(userList, u)
	}

	return userList
}

type prodData struct {
	P    *model.Product
	Tags []*model.ProductTag
	Imgs []*model.ProductImage
}

// parseProducts parses csv lines and creates the list of products
func parseProducts() []*prodData {
	records, err := readCSVFile("./data/seeds/products.csv")
	if err != nil {
		log.Fatalf("error parsing products from CSV: %v", err)
	}

	prodList := make([]*prodData, 0)

	for _, line := range records {
		p := &model.Product{}
		data := &prodData{}

		p.Name = line[0]
		p.BrandID, err = strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			log.Fatalf("brand_id err: %v", err)
		}
		p.CategoryID, err = strconv.ParseInt(line[2], 10, 64)
		if err != nil {
			log.Fatalf("category_id err: %v", err)
		}
		p.Slug = line[3]
		p.ImageURL = line[4]
		p.Description = line[5]
		p.Price, err = strconv.Atoi(line[6])
		if err != nil {
			log.Fatalf("price err: %v", err)
		}
		p.InStock, err = strconv.ParseBool(line[7])
		if err != nil {
			log.Fatalf("in_stock err: %v", err)
		}
		p.SKU = line[8]
		p.IsFeatured, err = strconv.ParseBool(line[9])
		if err != nil {
			log.Fatalf("is_featured err: %v", err)
		}

		tagSlice := strings.Split(line[10], " ")
		imgSlice := strings.Split(line[11], " ")

		tagList := make([]*model.ProductTag, 0)
		imgList := make([]*model.ProductImage, 0)

		for _, tid := range tagSlice {
			id, err := strconv.ParseInt(tid, 10, 64)
			if err != nil {
				log.Fatalf("tag_id err: %v", err)
			}
			tagList = append(tagList, &model.ProductTag{TagID: model.NewInt64(id)})
		}
		for _, img := range imgSlice {
			now := time.Now()
			imgList = append(imgList, &model.ProductImage{
				URL:       model.NewString(img),
				CreatedAt: &now,
				UpdatedAt: &now,
			})
		}

		data.P = p
		data.Tags = tagList
		data.Imgs = imgList

		prodList = append(prodList, data)
	}

	return prodList
}

// parseCategories parses csv lines and creates the list of categories
func parseCategories() []*model.Category {
	records, err := readCSVFile("./data/seeds/categories.csv")
	if err != nil {
		log.Fatalf("error parsing categoies from CSV: %v", err)
	}
	categoryList := make([]*model.Category, 0)

	for _, line := range records {
		c := &model.Category{}
		c.Name = line[0]
		c.Slug = line[1]
		c.Description = line[2]
		c.Logo = line[3]
		categoryList = append(categoryList, c)
	}

	return categoryList
}

// parseBrands parses csv lines and creates the list of brands
func parseBrands() []*model.Brand {
	records, err := readCSVFile("./data/seeds/brands.csv")
	if err != nil {
		log.Fatalf("error parsing brands from CSV: %v", err)
	}
	brandList := make([]*model.Brand, 0)

	for _, line := range records {
		b := &model.Brand{}
		b.Name = line[0]
		b.Slug = line[1]
		b.Type = line[2]
		b.Description = line[3]
		b.Email = line[4]
		b.WebsiteURL = line[5]
		b.Logo = line[6]
		brandList = append(brandList, b)
	}

	return brandList
}

// parseTags parses csv lines and creates the list of tags
func parseTags() []*model.Tag {
	records, err := readCSVFile("./data/seeds/tags.csv")
	if err != nil {
		log.Fatalf("error parsing tags from CSV: %v", err)
	}
	tagList := make([]*model.Tag, 0)

	for _, line := range records {
		t := &model.Tag{}
		t.Name = line[0]
		t.Slug = line[1]
		t.Description = line[2]
		tagList = append(tagList, t)
	}

	return tagList
}

// readCSVFile is a halper to read csv file
func readCSVFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("could not read input file: %v, %v", filePath, err)
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	if _, err := csvr.Read(); err != nil {
		return nil, err
	}

	records, err := csvr.ReadAll()
	if err != nil {
		log.Fatalf("could not parse file as CSV for path: %v ,%v", filePath, err)
		return nil, err
	}

	return records, nil
}
