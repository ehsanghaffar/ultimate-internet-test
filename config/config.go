package config

import "time"

// Config holds all configuration for internet tests
type Config struct {
	HTTPTimeout      time.Duration
	PingCount        int
	PingTimeout      time.Duration
	SpeedTestTimeout time.Duration
	ResultsFilePath  string
}

// Default configuration constants
const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 5 * time.Second

	// DefaultPingCount is the number of ping packets to send
	DefaultPingCount = 5

	// DefaultPingTimeout is the timeout for a single ping operation
	DefaultPingTimeout = 10 * time.Second

	// DefaultSpeedTestTimeout is the timeout for speed test requests
	DefaultSpeedTestTimeout = 10 * time.Second

	// DefaultResultsFilePath is the default path for storing test results
	DefaultResultsFilePath = "data.json"

	// BytesToBits conversion factor (for Mbps calculation)
	BytesToBits = 8

	// BytesToMegabytes conversion factor
	BytesToMegabytes = 1000000

	// FilePermissions for created JSON files (rw-r--r--)
	FilePermissions = 0644
)

// New creates a new Config with default values
func New() *Config {
	return &Config{
		HTTPTimeout:      DefaultHTTPTimeout,
		PingCount:        DefaultPingCount,
		PingTimeout:      DefaultPingTimeout,
		SpeedTestTimeout: DefaultSpeedTestTimeout,
		ResultsFilePath:  DefaultResultsFilePath,
	}
}
