package cmd

import (
	"awesomeProject/prometheus"
	"awesomeProject/worker"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var configFile = "config/config.yaml"

type Config struct {
	WorkerConfig  worker.Config     `yaml:"worker"`
	StatistConfig prometheus.Config `yaml:"statist"`
}

func getConfig(cmd *cobra.Command) *Config {
	filePath, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatalf("get flag err: %s", err)
	}

	if filePath != "" {
		configFile = fmt.Sprintf("config/%s", filePath)
	}

	var cfg Config
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("open config file \"%s\" err: %s", configFile, err)
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("decode config file err: %s", err)
	}

	return &cfg
}
