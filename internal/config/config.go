package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	PoolAddress  string        `mapstructure:"pool_address"`
	PoolPort     int           `mapstructure:"pool_port"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
	CPUTarget    float64       `mapstructure:"cpu_target"`
	StorePath    string        `mapstructure:"store_path"`
	Theme        string        `mapstructure:"theme"`
	MockMining   bool          `mapstructure:"mock_mining"`
	BTCAddress   string        `mapstructure:"btc_address"`
	WorkerName   string        `mapstructure:"worker_name"`
}

// Validate ensures that the config values fall within acceptable boundaries.
func (c *Config) Validate() error {
	if c.CPUTarget < 0.05 || c.CPUTarget > 0.75 {
		return fmt.Errorf("cpu_target %.2f is out of bounds [0.05, 0.75]", c.CPUTarget)
	}
	if !c.MockMining && c.BTCAddress == "" {
		return errors.New("btc_address is required when mock_mining is false")
	}
	return nil
}

// Load loads the configuration from environment variables, files, and defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Apply defaults
	v.SetDefault("pool_address", "public-pool.io")
	v.SetDefault("pool_port", 21496)
	v.SetDefault("worker_name", ".nerdtui")
	v.SetDefault("poll_interval", 5*time.Second)
	v.SetDefault("cpu_target", 0.75)
	v.SetDefault("store_path", "~/.nerdtui/metrics.db")
	v.SetDefault("theme", "dark")
	v.SetDefault("mock_mining", false)

	v.SetEnvPrefix("NM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	cfgDir := os.Getenv("NM_CONFIG_DIR")
	if cfgDir == "" {
		cfgDir, _ = ExpandPath("~/.nerdtui")
	}
	var errRead error
	if cfgDir != "" {
		v.AddConfigPath(cfgDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		errRead = v.ReadInConfig() // capture error to know if config was loaded
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	// If config was successfully read from a file and we have a BTCAddress,
	// force mock_mining to false to start real mining immediately.
	if errRead == nil && c.BTCAddress != "" {
		c.MockMining = false
	}

	expandedPath, err := ExpandPath(c.StorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand store path: %w", err)
	}
	c.StorePath = expandedPath

	return &c, nil
}

// Save saves the configuration to the default config file path.
func Save(c *Config) error {
	v := viper.New()
	v.Set("pool_address", c.PoolAddress)
	v.Set("pool_port", c.PoolPort)
	v.Set("poll_interval", c.PollInterval)
	v.Set("cpu_target", c.CPUTarget)
	v.Set("store_path", c.StorePath)
	v.Set("theme", c.Theme)
	v.Set("mock_mining", c.MockMining)
	v.Set("btc_address", c.BTCAddress)
	v.Set("worker_name", c.WorkerName)
	
	// Create ~/.nerdtui directory if it doesn't exist
	cfgDir, _ := ExpandPath("~/.nerdtui")
	_ = os.MkdirAll(cfgDir, 0755)
	
	v.AddConfigPath(cfgDir)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	
	return v.WriteConfigAs(cfgDir + "/config.yaml")
}
