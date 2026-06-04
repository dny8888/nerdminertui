package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	t.Run("Valid Config", func(t *testing.T) {
		t.Parallel()
		c := &Config{
			CPUTarget:  0.75,
			MockMining: false,
			BTCAddress: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		}
		assert.NoError(t, c.Validate())
	})

	t.Run("Invalid CPU Target Low", func(t *testing.T) {
		t.Parallel()
		c := &Config{
			CPUTarget:  0.02,
			MockMining: true,
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cpu_target 0.02 is out of bounds")
	})

	t.Run("Invalid CPU Target High", func(t *testing.T) {
		t.Parallel()
		c := &Config{
			CPUTarget:  0.80,
			MockMining: true,
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cpu_target 0.80 is out of bounds")
	})

	t.Run("Missing BTC Address", func(t *testing.T) {
		t.Parallel()
		c := &Config{
			CPUTarget:  0.5,
			MockMining: false,
			BTCAddress: "",
		}
		err := c.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "btc_address is required")
	})

	t.Run("Missing BTC Address but MockMining is true", func(t *testing.T) {
		t.Parallel()
		c := &Config{
			CPUTarget:  0.5,
			MockMining: true,
			BTCAddress: "",
		}
		assert.NoError(t, c.Validate())
	})
}

func TestExpandPath(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	assert.NoError(t, err)

	// With tilde
	expanded, err := ExpandPath("~/.nerdtui/test.db")
	assert.NoError(t, err)
	expected := filepath.Join(home, ".nerdtui/test.db")
	assert.Equal(t, expected, expanded)

	// Without tilde
	regular := "/tmp/test.db"
	expandedRegular, err := ExpandPath(regular)
	assert.NoError(t, err)
	assert.Equal(t, regular, expandedRegular)
}

func TestLoadDefaultsAndEnv(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv
	// Set mock env vars
	t.Setenv("NM_POOL_ADDRESS", "test-pool.io")
	t.Setenv("NM_CPU_TARGET", "0.70")
	t.Setenv("NM_MOCK_MINING", "true")

	// Temporarily override NM_CONFIG_DIR to a fake directory so it doesn't read the real user config
	t.Setenv("NM_CONFIG_DIR", "/tmp/nerdtui-mock-config-dir")

	c, err := Load("")
	assert.NoError(t, err)

	// Env overrides
	assert.Equal(t, "test-pool.io", c.PoolAddress)
	assert.Equal(t, 0.70, c.CPUTarget)
	assert.True(t, c.MockMining)

	// Defaults
	assert.Equal(t, 21496, c.PoolPort)
	assert.Equal(t, 5*time.Second, c.PollInterval)
	assert.Equal(t, "dark", c.Theme)

	// Path expansion default
	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, ".nerdtui/metrics.db")
	assert.Equal(t, expectedPath, c.StorePath)
}
