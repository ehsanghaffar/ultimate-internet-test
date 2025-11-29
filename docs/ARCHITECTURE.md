# Architecture & Design

## Project Overview

The all-internet-tests project is a modular Go application designed to perform comprehensive network diagnostics. It follows a clean architecture with separation of concerns across functional domains.

## Architecture Layers

### 1. **Configuration Layer** (`config/`)

Centralized configuration management providing:

- **config.go**: Defines the Config struct and exports all constants
  - HTTP/Ping/Speed test timeouts
  - Ping packet count
  - Result file path
  - Conversion factors (bytes to bits/megabytes)
  - File permissions

**Key benefit**: All magic numbers are centralized and easily configurable without code changes.

### 2. **Utility Layer** (`utils/`)

Shared utilities and data structures:

- **structs.go**: Complete type system
  - `TestResults`: Aggregates all test results with timestamp
  - `HTTPTest`: HTTP request results with TLS info
  - `SpeedTest`: Speed metrics and download information
  - `VPNTest`: VPN/Proxy detection results
  - `PingTest`: ICMP ping statistics

- **errors.go**: Custom error hierarchy
  - `TestError`: Base error type with context
  - `NetworkError`: Network-specific failures
  - `TimeoutError`: Timeout conditions
  - `ValidationError`: Input/data validation errors
  - `ParseError`: Parsing/marshaling errors

- **storage.go**: Persistent data management
  - Mutex-protected concurrent access to data.json
  - `LoadResults()`: Load previous test results
  - `SaveResults()`: Atomically save test results
  - `AppendResult()`: Append to existing results

### 3. **Module Layer** (`modules/`)

Discrete test implementations:

- **httpTest.go**: HTTP connectivity testing
  - Protocol and TLS information logging
  - Configurable timeouts
  - Response analysis

- **speedTest.go**: Internet speed measurement
  - Download speed calculation
  - Elapsed time tracking
  - Supports custom URLs

- **checkVPN.go**: VPN/Proxy detection
  - External IP detection
  - IP validation and comparison
  - Robust error handling

- **pingTest.go**: ICMP ping diagnostics
  - Packet loss calculation
  - Round-trip timing
  - Signal handling for graceful termination

### 4. **Application Layer** (`main.go`)

Orchestration and execution:

- Parallel test execution using goroutines and sync.WaitGroup
- Configuration initialization
- Command-line argument parsing
- Aggregated result collection and persistence

## Design Patterns Used

### Concurrent Execution

```sh
main() 
  ├─ HTTP Tests (parallel)
  ├─ Speed Tests (parallel)
  ├─ VPN Check (sequential)
  └─ Ping Tests (parallel)
  └─ Aggregated Save
```

**Benefit**: All independent tests run concurrently, reducing total execution time.

### Error Handling Chain

```sh
Operation
  ├─ Success → Return typed result
  ├─ Failure → Return result with Error field set
  └─ Critical → Return result with error field + logging
```

**Benefit**: No panic-driven crashes; graceful degradation with context preservation.

### Mutex-Protected Resource Access

```sh
File I/O
  ├─ Lock acquired
  ├─ Read/Write operation
  ├─ Lock released
  └─ Safe concurrent access
```

**Benefit**: Data integrity during concurrent test execution and result updates.

## Data Flow

### Test Execution Flow

```sh
┌─────────────┐
│   main()    │
│ init config │
└──────┬──────┘
       │
       ├─────────────┬─────────────┬──────────────┐
       │             │             │              │
       ▼             ▼             ▼              ▼
   ┌─────────┐ ┌──────────┐ ┌─────────┐ ┌──────────┐
   │TestHTTP │ │CheckSpeed│ │CheckVPN │ │PingCheck │
   └────┬────┘ └────┬─────┘ └────┬────┘ └────┬─────┘
        │           │            │           │
        └───────────┴────────────┴───────────┘
                     │
                     ▼
          ┌─────────────────────┐
          │  utils.TestResults  │
          │  [struct union]     │
          └──────────┬──────────┘
                     │
                     ▼
          ┌─────────────────────┐
          │  utils.SaveResults()│
          │  Mutex protected    │
          └──────────┬──────────┘
                     │
                     ▼
              data.json (file)
```

### Result Aggregation

```sh
Individual Results
  ├─ HTTPTest[]
  ├─ SpeedTest[]
  ├─ VPNTest
  ├─ PingTest
  └─ Timestamp

        ↓ [combine]

TestResults {
  HTTPTests:  []HTTPTest
  SpeedTests: []SpeedTest
  VPNTest:    VPNTest
  PingTest:   PingTest
  Timestamp:  time.Time
}
```

## Package Dependencies

```sh
main
  ├── config
  ├── modules
  │   ├── config
  │   ├── utils
  │   └── github.com/go-ping/ping
  └── utils
      ├── config (for FilePermissions)
      └── (stdlib: encoding/json, os, sync, time)

modules (all files)
  ├── config
  ├── utils
  ├── stdlib (io, log, net, net/http, regexp)
  └── go-ping/ping (ping only)
```

## Extension Points

### Adding a New Test Type

1. **Define result type** in `utils/structs.go`:

   ```go
   type CustomTest struct {
       URL    string
       Status string
       Error  string
   }
   ```

2. **Add to TestResults** in `utils/structs.go`:

   ```go
   type TestResults struct {
       CustomTests []CustomTest `json:"custom_tests,omitempty"`
       // ... other fields
   }
   ```

3. **Create module** in `modules/customTest.go`:

   ```go
   func CustomCheck(url string, cfg *config.Config) *utils.CustomTest { ... }
   ```

4. **Integrate in main.go**:

   ```go
   wg.Add(1)
   go func() {
       defer wg.Done()
       result := modules.CustomCheck(url, cfg)
       mu.Lock()
       customTests = append(customTests, result)
       mu.Unlock()
   }(url)
   ```

### Adding Configuration Options

1. Add to `config.go` struct:

   ```go
   type Config struct {
       NewOption time.Duration
   }
   ```

2. Add constant:

   ```go
   const DefaultNewOption = 10 * time.Second
   ```

3. Initialize in `New()`:

   ```go
   func New() *Config {
       return &Config{
           NewOption: DefaultNewOption,
       }
   }
   ```

## Performance Considerations

- **Concurrent Execution**: All tests run in parallel using goroutines
- **Single File Write**: One atomic write to data.json per run (vs. multiple writes)
- **Mutex Protection**: Minimal lock contention; short critical sections
- **Memory Efficiency**: Results streamed, not stored in intermediate buffers
- **Network Timeouts**: Configurable per-test to prevent hanging

## Security Considerations

- **Input Validation**: URL parsing, IP address validation
- **Error Contexts**: Errors preserve context without exposing sensitive data
- **File Permissions**: Explicit file permissions (0644) on JSON output
- **No Hardcoded Credentials**: Configuration-driven endpoint URLs
- **Graceful Degradation**: Individual test failures don't crash the application

## Testing Strategy

Tests use:

- **Mock HTTP servers** (httptest) for isolated testing
- **Table-driven tests** for multiple scenarios
- **Error path verification** for robustness
- **Concurrent access testing** for data integrity
