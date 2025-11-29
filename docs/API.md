# API Reference

## Config Package

### Types

#### Config

```go
type Config struct {
    HTTPTimeout      time.Duration
    PingCount        int
    PingTimeout      time.Duration
    SpeedTestTimeout time.Duration
    ResultsFilePath  string
}
```

Configuration for test execution with all timing and path settings.

### Constants

- `DefaultHTTPTimeout` = 5 seconds
- `DefaultPingCount` = 5 packets
- `DefaultPingTimeout` = 10 seconds
- `DefaultSpeedTestTimeout` = 10 seconds
- `DefaultResultsFilePath` = "data.json"
- `BytesToBits` = 8
- `BytesToMegabytes` = 1000000
- `FilePermissions` = 0644

### Functions

#### New() \*Config

```go
func New() *Config
```

Creates a new Config with default values. Use this to initialize the configuration for your test run.

**Example:**

```go
cfg := config.New()
// Modify as needed
cfg.HTTPTimeout = 10 * time.Second
```

---

## Utils Package

### Type

#### TestResults

```go
type TestResults struct {
    HTTPTests  []HTTPTest
    SpeedTests []SpeedTest
    VPNTest    VPNTest
    PingTest   PingTest
    Timestamp  time.Time
}
```

Aggregated results from all tests run in a single execution.

#### HTTPTest

```go
type HTTPTest struct {
    URL            string
    Status         string
    Proto          string
    TLSVersion     string
    CipherSuite    string
    ServerName     string
    ResponseLength int
    Error          string
}
```

Results from an individual HTTP test including protocol and TLS details.

#### SpeedTest

```go
type SpeedTest struct {
    URL         string
    DownloadMbps float64
    ElapsedTime  time.Duration
    BytesReceived int
    Error        string
}
```

Internet speed measurement results.

#### VPNTest

```go
type VPNTest struct {
    Status string
    Error  string
}
```

VPN/Proxy detection results. Status indicates "Using VPN or proxy." or "Not using VPN or proxy."

#### PingTest

```go
type PingTest struct {
    URL         string
    Transmitted int
    Received    int
    Loss        float64
    Error       string
}
```

ICMP ping statistics and packet loss information.

### Error Types

#### TestError

Base error type for all test failures. Implements `error` interface with `Unwrap()` support for error chain inspection.

```go
type TestError struct {
    TestType string
    Message  string
    Err      error
}
```

#### NetworkError

Network-related failures (connection refused, DNS resolution failed, etc.)

```go
type NetworkError struct {
    *TestError
}
```

#### TimeoutError

Timeout during test execution.

```go
type TimeoutError struct {
    *TestError
}
```

#### ValidationError

Input or data validation failures.

```go
type ValidationError struct {
    *TestError
}
```

#### ParseError

JSON marshaling/unmarshaling errors.

```go
type ParseError struct {
    *TestError
}
```

### Function

#### LoadResults(filePath string) (\*TestResults, error)

```go
func LoadResults(filePath string) (*TestResults, error)
```

Loads test results from a JSON file. Returns empty results if file doesn't exist.

**Parameters:**

- `filePath`: Path to the results JSON file

**Returns:**

- Pointer to TestResults struct
- Error if read/parse fails

**Example:**

```go
results, err := utils.LoadResults("data.json")
if err != nil {
    log.Fatal(err)
}
```

#### SaveResults(results \*TestResults, filePath string, filePermissions os.FileMode) error

```go
func SaveResults(results *TestResults, filePath string, filePermissions os.FileMode) error
```

Atomically saves test results to a JSON file with mutex protection for concurrent access.

**Parameters:**

- `results`: TestResults to save
- `filePath`: Destination file path
- `filePermissions`: File mode (e.g., 0644)

**Returns:**

- Error if operation fails

**Example:**

```go
results := &utils.TestResults{
    HTTPTests: httpTests,
    PingTest:  pingResult,
}
err := utils.SaveResults(results, "data.json", config.FilePermissions)
```

#### AppendResult(...) error

```go
func AppendResult(httpTests []HTTPTest, speedTests []SpeedTest,
    vpnTest *VPNTest, pingTest *PingTest,
    filePath string, filePermissions os.FileMode) error
```

Appends new test results to existing results file and saves atomically.

---

## Modules Package

### HTTP Testing

#### TestHTTP(url string, cfg *config.Config)*utils.HTTPTest

```go
func TestHTTP(url string, cfg *config.Config) *utils.HTTPTest
```

Performs an HTTP GET request and returns response details including TLS information.

**Parameters:**

- `url`: HTTP or HTTPS URL to test
- `cfg`: Config with timeout settings

**Returns:**

- HTTPTest result (never nil; check Error field)

**Example:**

```go
cfg := config.New()
result := modules.TestHTTP("https://example.com", cfg)
if result.Error != "" {
    log.Printf("Test failed: %s\n", result.Error)
} else {
    log.Printf("Status: %s, TLS: %s\n", result.Status, result.TLSVersion)
}
```

### Speed Testing

#### CheckSpeed(url string, cfg *config.Config)*utils.SpeedTest

```go
func CheckSpeed(url string, cfg *config.Config) *utils.SpeedTest
```

Downloads from URL and calculates download speed in Mbps.

**Parameters:**

- `url`: URL to download from
- `cfg`: Config with timeout settings

**Returns:**

- SpeedTest result (never nil; check Error field)

**Example:**

```go
result := modules.CheckSpeed("https://example.com/largefile", cfg)
if result.Error == "" {
    log.Printf("Speed: %.2f Mbps\n", result.DownloadMbps)
}
```

### VPN Detection

#### CheckVPN(ipChecker string) \*utils.VPNTest

```go
func CheckVPN(ipChecker string) *utils.VPNTest
```

Detects VPN/proxy by comparing external IP with local IP.

**Parameters:**

- `ipChecker`: URL of IP detection service (e.g., "checkip.dyndns.org")

**Returns:**

- VPNTest result (never nil; check Error field)

**Example:**

```go
result := modules.CheckVPN("http://checkip.dyndns.org/")
if result.Error == "" {
    log.Println(result.Status)
}
```

### Ping Testing

#### PingCheck(domain string, cfg *config.Config)*utils.PingTest

```go
func PingCheck(domain string, cfg *config.Config) *utils.PingTest
```

Sends ICMP ping packets and collects statistics.

**Parameters:**

- `domain`: Domain or IP address to ping
- `cfg`: Config with ping count setting

**Returns:**

- PingTest result (never nil; check Error field)

**Example:**

```go
result := modules.PingCheck("example.com", cfg)
if result.Error == "" {
    log.Printf("Packets: %d sent, %d received (%.1f%% loss)\n",
        result.Transmitted, result.Received, result.Loss)
}
```

---

## Error Handling

All module functions return typed result pointers, never errors directly. Check the `Error` field on the result:

```go
result := modules.TestHTTP(url, cfg)
if result.Error != "" {
    log.Printf("HTTP test failed: %s\n", result.Error)
    return
}
// Process successful result
```

For storage operations, errors are returned directly:

```go
err := utils.SaveResults(results, "data.json", config.FilePermissions)
if err != nil {
    log.Printf("Save failed: %v\n", err)
}
```
