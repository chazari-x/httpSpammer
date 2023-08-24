package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "balance_checker",
	Short: "Balance checker",
	Long:  "Balance checker",
	Run:   func(cmd *cobra.Command, args []string) {},
}

type contextKey string

func Execute() {
	ctx := context.Background()
	f, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("open config file: %v", err)
	}

	var cfg *Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("decode config file: %v", err)
	}

	ctx = context.WithValue(ctx, cfgKey, cfg)
	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		panic(err)
	}
}
