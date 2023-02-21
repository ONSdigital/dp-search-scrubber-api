package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig(t *testing.T) {
	// Call the Get function to get the default configuration
	config, err := Get()
	assert.Nil(t, err)

	// Assert that the configuration has the default values
	assert.Equal(t, ":3002", config.BindAddr)
	assert.Equal(t, 5*time.Second, config.GracefulShutdownTimeout)
	assert.Equal(t, 30*time.Second, config.HealthCheckInterval)
	assert.Equal(t, 90*time.Second, config.HealthCheckCriticalTimeout)
	assert.Equal(t, "data/2011 OAC Clusters and Names csv v2.csv", config.AreaDataFile)
	assert.Equal(t, "data/SIC07_CH_condensed_list_en.csv", config.IndustryDataFile)
}

func TestGetConfigFromEnv(t *testing.T) {
	// Set environment variables to modify the default configuration
	os.Setenv("BIND_ADDR", ":8080")
	os.Setenv("GRACEFUL_SHUTDOWN_TIMEOUT", "10s")
	os.Setenv("HEALTHCHECK_INTERVAL", "60s")
	os.Setenv("HEALTHCHECK_CRITICAL_TIMEOUT", "180s")
	os.Setenv("AREA_DATA_FILE", "data/areas.csv")
	os.Setenv("INDUSTRY_DATA_FILE", "data/industries.csv")

	// Call the Get function to get the modified configuration
	config, err := Get()
	assert.Nil(t, err)

	// Assert that the configuration has the modified values
	assert.Equal(t, ":8080", config.BindAddr)
	assert.Equal(t, 10*time.Second, config.GracefulShutdownTimeout)
	assert.Equal(t, 60*time.Second, config.HealthCheckInterval)
	assert.Equal(t, 180*time.Second, config.HealthCheckCriticalTimeout)
	assert.Equal(t, "data/areas.csv", config.AreaDataFile)
	assert.Equal(t, "data/industries.csv", config.IndustryDataFile)

	// Unset the environment variables
	os.Unsetenv("BIND_ADDR")
	os.Unsetenv("GRACEFUL_SHUTDOWN_TIMEOUT")
	os.Unsetenv("HEALTHCHECK_INTERVAL")
	os.Unsetenv("HEALTHCHECK_CRITICAL_TIMEOUT")
	os.Unsetenv("AREA_DATA_FILE")
	os.Unsetenv("INDUSTRY_DATA_FILE")
}
