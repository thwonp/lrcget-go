package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	DatabasePath string `json:"database_path" validate:"required"`
	MaxWorkers   int    `json:"max_workers" validate:"min=1,max=100"`
	Timeout      int    `json:"timeout" validate:"min=1"`
	LRCLibURL    string `json:"lrclib_url" validate:"required,url"`
	LogLevel     string `json:"log_level" validate:"oneof=debug info warn error"`
	DataDir      string `json:"data_dir"`
	CacheSize    int    `json:"cache_size" validate:"min=0"`
	EnableMetrics bool  `json:"enable_metrics"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DatabasePath:  "~/.lrcget/db.sqlite3",
		MaxWorkers:    10,
		Timeout:       30,
		LRCLibURL:     "https://lrclib.net",
		LogLevel:      "info",
		DataDir:       "~/.lrcget",
		CacheSize:     100,
		EnableMetrics: true,
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()
	
	// Override with environment variables if they exist
	if dbPath := os.Getenv("LRCGET_DB_PATH"); dbPath != "" {
		config.DatabasePath = dbPath
	}
	
	if maxWorkers := os.Getenv("LRCGET_MAX_WORKERS"); maxWorkers != "" {
		if val, err := strconv.Atoi(maxWorkers); err == nil {
			config.MaxWorkers = val
		}
	}
	
	if timeout := os.Getenv("LRCGET_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			config.Timeout = val
		}
	}
	
	if lrclibURL := os.Getenv("LRCGET_LRCLIB_URL"); lrclibURL != "" {
		config.LRCLibURL = lrclibURL
	}
	
	if logLevel := os.Getenv("LRCGET_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}
	
	if dataDir := os.Getenv("LRCGET_DATA_DIR"); dataDir != "" {
		config.DataDir = dataDir
	}
	
	if cacheSize := os.Getenv("LRCGET_CACHE_SIZE"); cacheSize != "" {
		if val, err := strconv.Atoi(cacheSize); err == nil {
			config.CacheSize = val
		}
	}
	
	if enableMetrics := os.Getenv("LRCGET_ENABLE_METRICS"); enableMetrics != "" {
		if val, err := strconv.ParseBool(enableMetrics); err == nil {
			config.EnableMetrics = val
		}
	}
	
	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DatabasePath == "" {
		return fmt.Errorf("database_path is required")
	}
	
	if c.MaxWorkers < 1 || c.MaxWorkers > 100 {
		return fmt.Errorf("max_workers must be between 1 and 100")
	}
	
	if c.Timeout < 1 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if c.LRCLibURL == "" {
		return fmt.Errorf("lrclib_url is required")
	}
	
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.LogLevel) {
		return fmt.Errorf("log_level must be one of: %v", validLogLevels)
	}
	
	if c.CacheSize < 0 {
		return fmt.Errorf("cache_size must be non-negative")
	}
	
	return nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	return nil
}

// GetTimeoutDuration returns the timeout as a duration
func (c *Config) GetTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// GetCacheSize returns the cache size
func (c *Config) GetCacheSize() int {
	if c.CacheSize <= 0 {
		return 100 // Default cache size
	}
	return c.CacheSize
}

// GetMaxWorkers returns the maximum number of workers
func (c *Config) GetMaxWorkers() int {
	if c.MaxWorkers <= 0 {
		return 10 // Default worker count
	}
	return c.MaxWorkers
}

// GetLogLevel returns the log level
func (c *Config) GetLogLevel() string {
	if c.LogLevel == "" {
		return "info" // Default log level
	}
	return c.LogLevel
}

// GetDataDir returns the data directory
func (c *Config) GetDataDir() string {
	if c.DataDir == "" {
		return "~/.lrcget" // Default data directory
	}
	return c.DataDir
}

// GetDatabasePath returns the database path
func (c *Config) GetDatabasePath() string {
	if c.DatabasePath == "" {
		return "~/.lrcget/db.sqlite3" // Default database path
	}
	return c.DatabasePath
}

// GetLRCLibURL returns the LRCLIB URL
func (c *Config) GetLRCLibURL() string {
	if c.LRCLibURL == "" {
		return "https://lrclib.net" // Default LRCLIB URL
	}
	return c.LRCLibURL
}

// IsMetricsEnabled returns whether metrics are enabled
func (c *Config) IsMetricsEnabled() bool {
	return c.EnableMetrics
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
