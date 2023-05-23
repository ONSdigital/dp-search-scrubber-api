package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-nlp-search-scrubber
type Config struct {
	AreaDataFile               string        `envconfig:"AREA_DATA_FILE"`
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	IndustryDataFile           string        `envconfig:"INDUSTRY_DATA_FILE"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	cfg = &Config{
		AreaDataFile:               "data/2011 OAC Clusters and Names csv v2.csv",
		BindAddr:                   ":28700",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		IndustryDataFile:           "data/SIC07_CH_condensed_list_en.csv",
	}

	return cfg, envconfig.Process("", cfg)
}
