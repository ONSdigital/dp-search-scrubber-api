package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-nlp-search-scrubber
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	AreaDataFile               string
	IndustryDataFile           string
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	cfg := &Config{}

	// default arg for .Load() is .env
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg = &Config{
		BindAddr:                   ":3002",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		AreaDataFile:               "data/2011 OAC Clusters and Names csv v2.csv",
		IndustryDataFile:           "data/SIC07_CH_condensed_list_en.csv",
	}

	return cfg, envconfig.Process("", cfg)
}
