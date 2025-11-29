# Documentation Index

Welcome to the ultimate-internet-test documentation. This project provides comprehensive network diagnostics through parallel testing of HTTP connectivity, speed, VPN detection, and ping.

## Quick Navigation

### 🚀 Getting Started

Start here if you're new to the project:

1 **[Quick Reference](QUICK_REFERENCE.md)** - Fast commands and code examples

- Installation and setup
- Basic usage examples
- Common commands
- Troubleshooting quick links

### 📚 Comprehensive Guides

For detailed information:

2 **[Architecture](ARCHITECTURE.md)** - System design and structure

- Project overview
- Architecture layers
- Design patterns
- Data flow diagrams
- Extension points

3 **[API Reference](API.md)** - Complete API documentation

- Config package
- Utils package and types
- Error handling
- Module functions
- Code examples

4 **[Development Guide](DEVELOPMENT.md)** - Building and extending

- Setup and build instructions
- Project structure
- Adding new test types
- Configuration management
- Concurrency patterns
- Debugging techniques

5 **[Troubleshooting](TROUBLESHOOTING.md)** - Issues and solutions

- Common problems
- Platform-specific issues
- Debugging techniques
- Performance optimization

## Documentation Overview

| Document                                 | Purpose                         | Audience               |
| ---------------------------------------- | ------------------------------- | ---------------------- |
| [QUICK_REFERENCE.md](QUICK_REFERENCE.md) | Fast lookup and examples        | Everyone               |
| [ARCHITECTURE.md](ARCHITECTURE.md)       | System design and patterns      | Developers, Architects |
| [API.md](API.md)                         | Function and type documentation | Developers             |
| [DEVELOPMENT.md](DEVELOPMENT.md)         | Build and extend guide          | Contributors           |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Problem solving                 | Operators, Developers  |

## Key Concepts

### Core Components

- **config**: Centralized configuration with defaults
- **modules**: Independent test implementations (HTTP, Speed, VPN, Ping)
- **utils**: Shared types, errors, and persistence
- **main**: Orchestration and parallel execution

### Design Highlights

- ✅ **Type Safety**: Strongly-typed result structures
- ✅ **Error Resilience**: Graceful handling with no panics
- ✅ **Concurrency**: All tests run in parallel
- ✅ **Persistence**: Atomic JSON I/O with mutex protection
- ✅ **Configurability**: All magic numbers as constants

### Testing Results

All tests generate structured JSON results:

```json
{
  "http_tests": [
    {
      "url": "https://...",
      "status": "200 OK",
      "proto": "HTTP/2.0"
    }
  ],
  "speed_tests": [...],
  "vpn_test": {...},
  "ping_test": {...},
  "timestamp": "2024-01-15T..."
}
```

## Getting Help

### Quick Issues

**Build fails?** → [QUICK_REFERENCE.md - Troubleshooting](QUICK_REFERENCE.md#troubleshooting)

**Permission denied?** → [TROUBLESHOOTING.md - Ping Requires Elevated Privileges](TROUBLESHOOTING.md#ping-requires-elevated-privileges)

**Network timeouts?** → [TROUBLESHOOTING.md - Network Timeout Errors](TROUBLESHOOTING.md#network-timeout-errors)

### Detailed Help

**Want to extend?** → [DEVELOPMENT.md - Adding New Test Type](DEVELOPMENT.md#adding-a-new-test-type)

**Need API docs?** → [API.md - Complete Reference](API.md)

**Understand design?** → [ARCHITECTURE.md - Full Overview](ARCHITECTURE.md)

## Installation

```bash
git clone https://github.com/ehsanghaffar/ultimate-internet-test.git
cd ultimate-internet-test
go build
./ultimate-internet-test
```

See [QUICK_REFERENCE.md](QUICK_REFERENCE.md#installation--setup) for more options.

## Development Workflow

### Setup Development Environment

```bash
# Clone and navigate
git clone https://github.com/ehsanghaffar/ultimate-internet-test.git
cd ultimate-internet-test

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build

# Run
./ultimate-internet-test
```

See [DEVELOPMENT.md](DEVELOPMENT.md#getting-started) for detailed setup.

### Common Development Tasks

```bash
# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o ultimate-internet-test-linux

# Run tests with coverage
go test -cover ./...

# Debug with race detector
go test -race ./...

# Format and lint
go fmt ./...
go vet ./...
```

See [QUICK_REFERENCE.md - Building](QUICK_REFERENCE.md#building) for more commands.

## Testing Examples

### Run All Tests

```bash
./ultimate-internet-test
```

### Test Specific URLs

```bash
./ultimate-internet-test https://example.com https://test.com
```

### Check Results

```bash
cat data.json | jq '.'
```

See [QUICK_REFERENCE.md - Basic Usage](QUICK_REFERENCE.md#basic-usage) for more.

## API Quick Start

```go
import (
    "github.com/ehsanghaffar/ultimate-internet-test/config"
    "github.com/ehsanghaffar/ultimate-internet-test/modules"
)

// Create config
cfg := config.New()

// Run tests
httpResult := modules.TestHTTP("https://example.com", cfg)
speedResult := modules.CheckSpeed("https://example.com", cfg)
vpnResult := modules.CheckVPN("http://checkip.dyndns.org/")
pingResult := modules.PingCheck("example.com", cfg)

// Check results
if httpResult.Error == "" {
    log.Println("HTTP test passed:", httpResult.Status)
}
```

See [API.md](API.md) for complete documentation.

## Architecture Layers

```sh
┌─────────────────────────────────────────┐
│ Application (main.go)                   │
│ - Orchestration, parallel execution     │
├─────────────────────────────────────────┤
│ Modules (modules/)                      │
│ - HTTP, Speed, VPN, Ping tests          │
├─────────────────────────────────────────┤
│ Utils (utils/)                          │
│ - Types, errors, persistence            │
├─────────────────────────────────────────┤
│ Config (config/)                        │
│ - Configuration, constants              │
└─────────────────────────────────────────┘
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed information.

## FAQ

**Q: Do I need elevated privileges?**
A: Yes, for ping tests on Linux. Use `sudo` or set capabilities. See [TROUBLESHOOTING.md](TROUBLESHOOTING.md#ping-requires-elevated-privileges).

**Q: How do I customize test URLs?**
A: Edit `main.go` or modify the application to accept configuration. See [DEVELOPMENT.md](DEVELOPMENT.md#configuration-management).

**Q: Can I use this as a library?**
A: Yes! All packages are importable. See [API.md](API.md) and [QUICK_REFERENCE.md - Code Examples](QUICK_REFERENCE.md#code-examples).

**Q: How do I add a new test type?**
A: See [DEVELOPMENT.md - Adding a New Test Type](DEVELOPMENT.md#adding-a-new-test-type).

**Q: What about results persistence?**
A: Results are saved to `data.json` in JSON format. See [ARCHITECTURE.md - Data Flow](ARCHITECTURE.md#data-flow).

## Contributing

We welcome contributions! See [DEVELOPMENT.md - Contributing](DEVELOPMENT.md#contributing) for guidelines.

## Support

- 📧 Email: <ghafari.5000@gmail.com>
- 🐙 GitHub: <https://github.com/ehsanghaffar/ultimate-internet-test>
- 🌐 Website: <https://ehsanghaffarii.ir>

## License

MIT License - see [LICENSE](../LICENSE) file

---

**Last Updated**: 2024
**Go Version Required**: 1.19+
**Status**: Active Development
