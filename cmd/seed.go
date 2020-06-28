package cmd

import (
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

func init() {
	rootCmd.AddCommand(seedCmd)

}

func seedDatabaseFn(command *cobra.Command, args []string) error {
	users := []*model.User{
		{
			Username: "Test",
			Password: "Test_123",
			Email:    "test@test.com",
		}, {
			Username: "John",
			Password: "John$123",
			Email:    "john@gmail.com",
		},
		{
			Username: "Ana",
			Password: "Ana55$7",
			Email:    "ana@gmail.com",
		},
	}

	for _, u := range users {
		u.PreSave()
	}

	if err := cmdApp.Srv().Store.User().BulkInsert(users); err != nil {
		cmdApp.Log().Error("could not seed database", zlog.String("err: ", err.Message))
	}
	cmdApp.Log().Info("db seed completed")
	return nil
}
