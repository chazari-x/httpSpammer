package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Config struct {
	Address string `yaml:"address"`
}

type Statist interface {
	Add(t float64)
}

type Metrics struct {
	r prometheus.Histogram
}

func NewPrometheus(cfg *Config) Statist {
	reg := prometheus.NewRegistry()

	responseTime := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "HTTP response time in seconds.",
	})

	reg.MustRegister(responseTime)

	go func(reg *prometheus.Registry) {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
		log.Fatal(http.ListenAndServe(cfg.Address, nil))
	}(reg)

	return &Metrics{r: responseTime}
}

func (c *Metrics) Add(t float64) {
	c.r.Observe(t)
}
