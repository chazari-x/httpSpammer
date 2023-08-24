package cmd

import (
	"awesomeProject/prometheus"
	"awesomeProject/worker"
	"context"
)

var cfgKey = contextKey("config")

type Config struct {
	WorkerConfig  worker.Config     `yaml:"worker"`
	StatistConfig prometheus.Config `yaml:"statist"`
}

func configFromContext(ctx context.Context) *Config {
	return ctx.Value(cfgKey).(*Config)
}
