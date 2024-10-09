package main

import (
	"fmt"
	"log"
	"shop-graphql-demo/account"

	"github.com/spf13/cobra"
)

func main() {
	var cleanUp func()
	rootCmd := &cobra.Command{
		Use:  "account-srv",
		Long: "account service ",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}
			cfg, err := account.NewConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}
			repo, err := account.NewPostgresRepository(cfg.DBString)
			if err != nil {
				return fmt.Errorf("failed to create postgres repository: %w", err)
			}
			cleanUp = repo.Close
			srv := account.NewService(repo)
			return account.ListenGRPC(srv, cfg.Port)
		},
	}
	defer cleanUp()
	rootCmd.Flags().String("config", "", "path to config file")
	fmt.Println("account-srv starting")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
