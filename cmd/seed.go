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

func init() {
	seedCmd.AddCommand(seedUsersCmd, seedProductsCmd)
	rootCmd.AddCommand(seedCmd)
}

func seedDatabaseFn(command *cobra.Command, args []string) error {
	if err := seedUsers(); err != nil {
		return err
	}
	if err := seedProducts(); err != nil {
		return err
	}
	cmdApp.Log().Info("database seed completed")
	return nil
}

func seedUsersFn(command *cobra.Command, args []string) error {
	return seedUsers()
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
	cmdApp.Log().Info("users seed completed")
	return nil
}

// seedProducts populates the product tables together with related brand, categories, imgs, and tags
func seedProducts() error {
	data := parseProducts()

	for _, d := range data {
		d.P.PreSave()
		newProd, err := cmdApp.Srv().Store.Product().Save(d.P)
		if err != nil {
			cmdApp.Log().Error("could not seed save product", zlog.String("err: ", err.Message))
		}

		for _, tag := range d.Tags {
			tag.ProductID = &newProd.ID
			tag.PreSave()
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

	cmdApp.Log().Info("products seed completed")
	return nil
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

// parseUsers parses read csv lines and creates the list of users
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

// parseUsers parses read csv lines and creates the list of products
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
		p.Slug = line[1]
		p.ImageURL = line[2]
		p.Description = line[3]
		p.Price, err = strconv.Atoi(line[4])
		if err != nil {
			log.Fatalf("price err: %v", err)
		}
		p.InStock, err = strconv.ParseBool(line[5])
		if err != nil {
			log.Fatalf("in_stock err: %v", err)
		}
		p.SKU = line[6]
		p.IsFeatured, err = strconv.ParseBool(line[7])
		if err != nil {
			log.Fatalf("is_featured err: %v", err)
		}

		p.Category = &model.ProductCategory{
			Name:        line[8],
			Slug:        line[9],
			Description: line[10],
		}

		p.Brand = &model.ProductBrand{
			Name:        line[11],
			Slug:        line[12],
			Type:        line[13],
			Description: line[14],
			Email:       line[15],
			WebsiteURL:  line[16],
		}

		tagSlice := strings.Split(line[17], " ")
		imgSlice := strings.Split(line[18], " ")

		tagList := make([]*model.ProductTag, 0)
		imgList := make([]*model.ProductImage, 0)

		for _, tag := range tagSlice {
			now := time.Now()
			tagList = append(tagList, &model.ProductTag{
				Name:      model.NewString(tag),
				CreatedAt: &now,
				UpdatedAt: &now,
			})
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
