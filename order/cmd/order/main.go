package main

import (
	"log"
	"shop-graphql-demo/order"

	"github.com/spf13/cobra"
)

func main() {
	var cleanUp func()
	cmd := cobra.Command{
		Use: "order-srv",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			cfg, err := order.NewConfig(configPath)
			if err != nil {
				return err
			}
			repo, err := order.NewOrderRepository(cfg.DBString)
			if err != nil {
				log.Printf("db string %s", cfg.DBString)
				return err
			}
			cleanUp = repo.Close
			ser := order.NewOrderService(repo)
			return order.ListenGRPC(ser, cfg.AccountService, cfg.CatalogService, cfg.Port)
		},
	}
	defer cleanUp()
	cmd.Flags().String("config", "", "config file")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
