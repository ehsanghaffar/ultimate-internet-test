# Quick Reference

## Installation & Setup

```bash
# Clone
git clone https://github.com/ehsanghaffar/ultimate-internet-test.git
cd ultimate-internet-test

# Build
go build -o ultimate-internet-test

# Run
./ultimate-internet-test
```

## Basic Usage

```bash
# Run all tests
go run .

# Build and run
go build && ./ultimate-internet-test

# Test specific URLs
go run . https://example.com https://test.com

# Redirect output
./ultimate-internet-test > results.txt

# With elevated privileges (for ping)
sudo go run .
```

## Code Examples

### Basic Test Execution

```go
package main

import (
    "github.com/ehsanghaffar/ultimate-internet-test/config"
    "github.com/ehsanghaffar/ultimate-internet-test/modules"
)

func main() {
    cfg := config.New()

    // HTTP test
    httpResult := modules.TestHTTP("https://example.com", cfg)
    if httpResult.Error != "" {
        log.Println("HTTP test failed:", httpResult.Error)
    }

    // Speed test
    speedResult := modules.CheckSpeed("https://example.com", cfg)
    if speedResult.Error == "" {
        log.Printf("Speed: %.2f Mbps\n", speedResult.DownloadMbps)
    }

    // VPN test
    vpnResult := modules.CheckVPN("http://checkip.dyndns.org/")
    if vpnResult.Error == "" {
        log.Println(vpnResult.Status)
    }

    // Ping test
    pingResult := modules.PingCheck("example.com", cfg)
    if pingResult.Error == "" {
        log.Printf("Packet loss: %.1f%%\n", pingResult.Loss)
    }
}
```

### Configuration

```go
cfg := config.New()

// Override defaults
cfg.HTTPTimeout = 10 * time.Second
cfg.PingCount = 10
cfg.SpeedTestTimeout = 20 * time.Second
cfg.ResultsFilePath = "results.json"

result := modules.TestHTTP("https://example.com", cfg)
```

### Save Results

```go
import (
    "github.com/ehsanghaffar/ultimate-internet-test/config"
    "github.com/ehsanghaffar/ultimate-internet-test/utils"
)

results := &utils.TestResults{
    HTTPTests:  []utils.HTTPTest{...},
    SpeedTests: []utils.SpeedTest{...},
    VPNTest:    utils.VPNTest{Status: "..."},
    PingTest:   utils.PingTest{...},
}

err := utils.SaveResults(results, "data.json", config.FilePermissions)
if err != nil {
    log.Println("Save failed:", err)
}
```

## Testing

```bash
# Run all tests
go test ./...

# Specific package
go test ./modules
go test ./utils

# Verbose output
go test -v ./...

# Coverage
go test -cover ./...

# Race detector
go test -race ./...

# Benchmarks
go test -bench=. ./...
```

## Building

```bash
# Development build
go build

# Production build (smaller)
go build -ldflags="-s -w" -o ultimate-internet-test

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o ultimate-internet-test-linux

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o ultimate-internet-test.exe

# Cross-compile for macOS
GOOS=darwin GOARCH=amd64 go build -o ultimate-internet-test-macos
```

## Troubleshooting

| Problem                        | Solution                                                  |
| ------------------------------ | --------------------------------------------------------- |
| Build error                    | `go mod tidy`                                             |
| Ping permission denied         | `sudo go run .` or `sudo setcap cap_net_raw=+ep ./binary` |
| Network timeout                | Increase timeout: `cfg.HTTPTimeout = 15 * time.Second`    |
| Connection refused             | Check if target is online: `curl https://target.com`      |
| Permission denied writing JSON | `chmod 644 .`                                             |
| Corrupted data.json            | Delete it, re-run program                                 |

## Project Structure

```sh
.
├── config/              # Configuration
├── modules/             # Test implementations
│   ├── httpTest.go
│   ├── speedTest.go
│   ├── checkVPN.go
│   └── pingTest.go
├── utils/               # Shared utilities
│   ├── structs.go       # Types
│   ├── errors.go        # Errors
│   └── storage.go       # Persistence
├── docs/                # Documentation
└── main.go              # Entry point
```

## Key Packages

### config

```go
cfg := config.New()  // Create config with defaults
```

Constants:

- `DefaultHTTPTimeout` = 5s
- `DefaultPingCount` = 5
- `BytesToBits` = 8
- `BytesToMegabytes` = 1000000
- `FilePermissions` = 0644

### modules

```go
TestHTTP(url, cfg)      // HTTP connectivity
CheckSpeed(url, cfg)    // Speed measurement
CheckVPN(url)           // VPN detection
PingCheck(domain, cfg)  // ICMP ping
```

### utils

Types:

- `TestResults` - Aggregated results
- `HTTPTest` - HTTP test result
- `SpeedTest` - Speed test result
- `VPNTest` - VPN detection result
- `PingTest` - Ping test result

Functions:

- `LoadResults(filePath)` - Load JSON results
- `SaveResults(results, filePath, perms)` - Save JSON results
- `AppendResult(...)` - Append to existing results

## Command Reference

```bash
# Build and run
go build && ./ultimate-internet-test

# Run with specific URLs
./ultimate-internet-test https://example.com https://test.com

# Run tests
go test ./... -v

# Check for issues
go vet ./...

# Format code
go fmt ./...

# Update dependencies
go mod tidy

# Build for deployment
go build -ldflags="-s -w"

# Run with profiling
go run -cpuprofile=cpu.prof .

# Debug with delve
dlv debug
```

## Error Handling Pattern

```go
result := modules.TestHTTP(url, cfg)

// Check Error field instead of returning error
if result.Error != "" {
    log.Printf("Test failed: %s\n", result.Error)
    // Continue with next test
    continue
}

// Process successful result
log.Printf("Status: %s\n", result.Status)
```

## Concurrency Pattern

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
// All done, results collected
```

## JSON Result Format

```json
{
  "http_tests": [...],
  "speed_tests": [...],
  "vpn_test": {...},
  "ping_test": {...},
  "timestamp": "2024-01-15T10:30:45Z"
}
```

## Documentation Files

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - System design and patterns
- **[API.md](API.md)** - Complete API documentation
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Development guide
- **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** - Issues and solutions

## External Resources

- [Go Documentation](https://golang.org/doc/)
- [go-ping/ping package](https://github.com/go-ping/ping)
- [httptest package](https://golang.org/pkg/net/http/httptest/)
- [sync.WaitGroup](https://golang.org/pkg/sync/#WaitGroup)

## Version Info

Check current Go version:

```bash
go version
```

Minimum required: Go 1.19

## Author

- [@ehsanghaffar](https://github.com/ehsanghaffar)
- Email: (mailto:ghafari.5000@gmail.com)
