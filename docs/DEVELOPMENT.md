# Development Guide

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Git
- Network access for tests

### Setup

```bash
# Clone the repository
git clone https://github.com/ehsanghaffar/ultimate-internet-test.git
cd ultimate-internet-test

# Install dependencies
go mod download

# Build the project
go build -o ultimate-internet-test

# Run tests
go test ./...
```

## Project Structure

```bash
.
├── config/                 # Configuration management
│   └── config.go          # Config struct and constants
├── modules/               # Test implementations
│   ├── httpTest.go        # HTTP connectivity tests
│   ├── speedTest.go       # Speed measurement
│   ├── checkVPN.go        # VPN detection
│   └── pingTest.go        # ICMP ping
├── utils/                 # Shared utilities
│   ├── structs.go         # Data types
│   ├── errors.go          # Error hierarchy
│   └── storage.go         # JSON persistence
├── docs/                  # Documentation
├── main.go               # Application entry point
├── go.mod                # Module definition
└── data.json             # Test results (generated)
```

## Building

### Development Build

```bash
go build
```

### Production Build

```bash
# Build optimized binary
go build -ldflags="-s -w" -o ultimate-internet-test
```

### Run Directly

```bash
go run .
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Specific Package Tests

```bash
go test ./modules
go test ./utils
```

### Run with Verbose Output

```bash
go test -v ./...
```

### Run with Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Benchmarks

```bash
go test -bench=. -benchmem ./...
```

## Code Organization

### Adding a New Test Type

1. **Define the result type** in `utils/structs.go`:

```go
type CustomTest struct {
    URL    string      `json:"url"`
    Result string      `json:"result"`
    Error  string      `json:"error,omitempty"`
}
```

2 **Add to TestResults** in `utils/structs.go`:

```go
type TestResults struct {
    // ... existing fields
    CustomTests []CustomTest `json:"custom_tests,omitempty"`
    Timestamp   time.Time    `json:"timestamp"`
}
```

3 **Create module file** `modules/customTest.go`:

```go
package modules

import (
    "github.com/ehsanghaffar/ultimate-internet-test/config"
    "github.com/ehsanghaffar/ultimate-internet-test/utils"
)

// CustomCheck performs a custom test and returns results
func CustomCheck(url string, cfg *config.Config) *utils.CustomTest {
    result := &utils.CustomTest{URL: url}

    // Implement test logic

    return result
}
```

4 **Integrate in main.go**:

```go
func runAllTests(cfg *config.Config) {
    // ... existing code ...

    var customTests []*utils.CustomTest

    // Run custom tests
    for _, url := range customURLs {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            result := modules.CustomCheck(u, cfg)
            mu.Lock()
            customTests = append(customTests, result)
            mu.Unlock()
        }(url)
    }

    // Convert results and add to TestResults
}
```

5 **Create tests** in `modules/customTest_test.go`:

```go
package modules

import (
    "testing"
    "github.com/ehsanghaffar/ultimate-internet-test/config"
)

func TestCustomCheck(t *testing.T) {
    cfg := config.New()
    result := CustomCheck("http://example.com", cfg)

    if result == nil {
        t.Fatal("Expected non-nil result")
    }
}
```

## Error Handling Best Practices

### Module Functions

All module test functions return a typed result pointer. Always check the `Error` field:

```go
result := modules.TestHTTP(url, cfg)
if result.Error != "" {
    log.Printf("HTTP test failed: %s\n", result.Error)
    // Handle gracefully, test continues
}
```

### Storage Functions

Storage functions return errors directly. Handle appropriately:

```go
if err := utils.SaveResults(results, path, perms); err != nil {
    log.Printf("Failed to save results: %v\n", err)
    return err
}
```

## Configuration Management

### Using Config

```go
// Create with defaults
cfg := config.New()

// Override specific settings
cfg.HTTPTimeout = 10 * time.Second
cfg.PingCount = 10

// Pass to test functions
result := modules.TestHTTP("https://example.com", cfg)
```

### Adding New Configuration Options

1 Add to `Config` struct in `config/config.go`:

```go
type Config struct {
    // ... existing fields
    CustomTimeout time.Duration
}
```

2 Add constant:

```go
const DefaultCustomTimeout = 15 * time.Second
```

3 Update `New()` function:

```go
func New() *Config {
    return &Config{
        // ... existing fields
        CustomTimeout: DefaultCustomTimeout,
    }
}
```

## Concurrency Patterns

### WaitGroup for Parallel Execution

```go
var wg sync.WaitGroup
var mu sync.Mutex

results := make([]Result, 0)

for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        result := process(i)

        mu.Lock()
        results = append(results, result)
        mu.Unlock()
    }(item)
}

wg.Wait()
// All goroutines completed, results collected
```

### Mutex Protection for Shared Resources

```go
var mu sync.Mutex

mu.Lock()
// Critical section - file I/O, shared data modification
data := readFile()
mu.Unlock()
```

## Debugging

### Enable Verbose Logging

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Print Debug Information

```go
log.Printf("Debug: %#v\n", result)
```

### Use Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug with breakpoints
dlv debug
```

## Performance Tips

1. **Reuse HTTP Client**: Create once, use many times
2. **Concurrent Tests**: Run independent tests in parallel (already implemented)
3. **Single File Write**: Aggregate results, write once (already implemented)
4. **Timeout Configuration**: Set appropriate timeouts to prevent hanging

## Common Issues

### Ping Requires Elevated Privileges

```bash
# Run with sudo
sudo go run .

# Or set capability on binary
sudo setcap cap_net_raw=+ep ./ultimate-internet-test
./ultimate-internet-test
```

### Network Timeout Errors

Increase timeout in config:

```go
cfg := config.New()
cfg.HTTPTimeout = 15 * time.Second
cfg.SpeedTestTimeout = 20 * time.Second
```

### Port Already in Use (Testing)

Ensure test servers are properly cleaned up:

```go
server := httptest.NewServer(handler)
defer server.Close()  // Important!
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-test`
3. Make changes with tests
4. Run: `go test ./...` and `go build`
5. Commit with clear messages
6. Push and create a Pull Request

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
