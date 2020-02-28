package cmd

import (
	"log"

	api "github.com/dankobgd/ecommerce-shop/api/v1"
	"github.com/dankobgd/ecommerce-shop/app"
	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/store/postgres"
	"github.com/dankobgd/ecommerce-shop/store/redis"
	"github.com/dankobgd/ecommerce-shop/store/supplier"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run the API web server",
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.RunE = serverCmdFn
}

func serverCmdFn(command *cobra.Command, args []string) error {
	db, err := postgres.Connect()
	if err != nil {
		return err
	}

	pgStore := postgres.NewStore(db)
	redisClient := redis.NewClient()
	redisStore := redis.NewStore(redisClient)

	storage := &supplier.Supplier{Pgst: pgStore, Rdst: redisStore}
	server, err := app.NewServer(storage)

	if err != nil {
		log.Fatalf("could not create the app server: %v\n", err)
	}

	cfg := config.New()

	logger := zlog.NewLogger(&zlog.LoggerConfig{
		EnableConsole: true,
		ConsoleLevel:  "debug",
		ConsoleJSON:   true,
		EnableFile:    true,
		FileLevel:     "info",
		FileJSON:      true,
		FileLocation:  "./logs/app.log",
	})

	zlog.InitGlobalLogger(logger)

	locale.InitTranslations()

	appOpts := []app.Option{
		app.SetConfig(cfg),
		app.SetServer(server),
		app.SetLogger(logger),
	}

	a := app.New(appOpts...)

	api.Init(a, server.Router)

	if srvErr := server.Start(); srvErr != nil {
		log.Fatalf("could not start the server: %v\n", srvErr)
	}

	return nil
}
