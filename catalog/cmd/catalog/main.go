package main

import (
	"log"
	"shop-graphql-demo/catalog"

	"github.com/spf13/cobra"
)

func main() {
	var cleanUp func()
	cmd := cobra.Command{
		Use: "catalog-srv",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			cfg, err := catalog.NewConfig(configPath)
			if err != nil {
				return err
			}
			repo, err := catalog.NewElasticRepository(cfg.DBString)
			if err != nil {
				log.Printf("db string %s", cfg.DBString)
				return err
			}
			cleanUp = repo.Close

			srv, _ := catalog.NewService(repo)
			return catalog.ListenGRPC(srv, cfg.Port)
		},
	}
	defer cleanUp()
	cmd.Flags().String("config", "", "config file")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
