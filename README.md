# All-in-one Internet Test

A comprehensive Go application for network diagnostics and internet connectivity testing. Performs parallel HTTP, speed, VPN detection, and ping tests with structured result storage and robust error handling.

## Features

- **HTTP Testing**: Test HTTP/HTTPS connectivity with TLS information
- **Speed Testing**: Measure download speed in Mbps
- **VPN Detection**: Detect if connection uses VPN or proxy
- **Ping Testing**: ICMP ping with packet loss statistics
- **Parallel Execution**: All tests run concurrently for faster execution
- **Structured Results**: Results saved to JSON with timestamps
- **Error Resilience**: Individual test failures don't crash the application
- **Type Safety**: Strongly-typed result structures
- **Configuration Management**: Centralized config with sensible defaults
- **Minimal Dependencies**: Only one external package (go-ping/ping)

## Quick Start

### Prerequisites

- Go 1.19 or higher
- Network access for tests
- Ping tests may require elevated privileges on Linux

### Installation

```bash
# Clone the repository
git clone https://github.com/ehsanghaffar/all-internet-tests.git
cd all-internet-tests

# Build
go build -o all-internet-tests

# Run
./all-internet-tests
```

### Usage

```bash
# Run all tests
go run .

# Build and execute binary
go build && ./all-internet-tests

# Run specific URL tests
go run . https://example.com https://test.com

# View results
cat data.json
```

## Documentation

- **[Architecture](docs/ARCHITECTURE.md)** - Design patterns, data flow, and system architecture
- **[API Reference](docs/API.md)** - Complete API documentation for all packages
- **[Development Guide](docs/DEVELOPMENT.md)** - Setup, building, testing, and extending
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions

## Project Structure

```bash
.
├── config/              # Configuration management
│   └── config.go        # Config struct, constants
├── modules/             # Test implementations
│   ├── httpTest.go      # HTTP connectivity tests
│   ├── speedTest.go     # Speed measurement
│   ├── checkVPN.go      # VPN/proxy detection
│   └── pingTest.go      # ICMP ping tests
├── utils/               # Shared utilities
│   ├── structs.go       # Result data types
│   ├── errors.go        # Error hierarchy
│   └── storage.go       # JSON persistence
├── docs/                # Documentation
│   ├── ARCHITECTURE.md
│   ├── API.md
│   ├── DEVELOPMENT.md
│   └── TROUBLESHOOTING.md
├── main.go              # Application entry
├── go.mod               # Module definition
└── data.json            # Results (generated)
```

## Example Output

```json
{
  "http_tests": [
    {
      "url": "https://www.google.com/",
      "status": "200 OK",
      "proto": "HTTP/2.0",
      "tls_version": "771",
      "response_length": 15678
    }
  ],
  "speed_tests": [
    {
      "url": "https://google.com",
      "download_mbps": 45.32,
      "elapsed_time": "2.345s",
      "bytes_received": 1048576
    }
  ],
  "vpn_test": {
    "status": "Not using VPN or proxy."
  },
  "ping_test": {
    "url": "www.google.com",
    "transmitted_packets": 5,
    "received_packets": 5,
    "loss_packets": 0
  },
  "timestamp": "2024-01-15T10:30:45Z"
}
```

## Configuration

The application uses sensible defaults but can be customized:

```go
cfg := config.New()
cfg.HTTPTimeout = 10 * time.Second
cfg.PingCount = 10
cfg.SpeedTestTimeout = 20 * time.Second
```

See [Configuration Constants](docs/API.md#constants) for all options.

## Development

### Build

```bash
go build -o all-internet-tests
```

### Test

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# With race detector
go test -race ./...
```

### Add Tests

Tests are located in `*_test.go` files. See [Development Guide](docs/DEVELOPMENT.md) for patterns and examples.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make changes with tests
4. Run: `go test ./... && go build`
5. Commit with clear messages
6. Push and create a Pull Request

## Architecture Highlights

- **Modular Design**: Independent test modules with clear interfaces
- **Type Safety**: Strongly-typed result structures for all tests
- **Error Handling**: Custom error hierarchy with proper context
- **Concurrency**: Parallel execution using goroutines and sync.WaitGroup
- **Data Persistence**: Mutex-protected JSON I/O for thread safety
- **Configuration**: Centralized config with all magic numbers as constants

See [Architecture Documentation](docs/ARCHITECTURE.md) for detailed information.

## Troubleshooting

### Common Issues

**Ping:Permission Denied**

```bash
# Linux: Use sudo
sudo go run .

# Or set capabilities
sudo setcap cap_net_raw=+ep ./all-internet-tests
./all-internet-tests
```

Network Timeouts

```go
cfg := config.New()
cfg.HTTPTimeout = 15 * time.Second
```

**VPN Detection Fails**
Use alternative IP detection service or check network connectivity.

See [Troubleshooting Guide](docs/TROUBLESHOOTING.md) for more solutions.

## Performance

- Tests run in parallel for faster execution
- Single atomic write to data.json per run
- Configurable timeouts to prevent hanging
- Minimal memory footprint with streaming I/O

## License

This project is licensed under the MIT License — see [LICENSE](https://choosealicense.com/licenses/mit/) for details.

## Authors and Contact

- **Original Author**: [@ehsanghaffar](https://github.com/ehsanghaffar)
- **Email**: (mailto:ghafari.5000@gmail.com)

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Project Issues](https://github.com/ehsanghaffar/all-internet-tests/issues)
