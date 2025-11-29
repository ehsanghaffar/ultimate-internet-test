# Troubleshooting Guide

## Common Issues and Solutions

### Build Issues

#### Error: cannot find module

**Symptom:** `cannot find module declared in go.mod`

**Solution:**

```bash
go mod download
go mod tidy
```

#### Error: go version too old

**Symptom:** `go version go1.17 does not support go1.19`

**Solution:**
Update Go to version 1.19 or higher:

```bash
# macOS (Homebrew)
brew upgrade go

# Or download from golang.org
```

### Runtime Issues

#### Ping: Permission Denied

**Symptom:** `operation not permitted` when running ping tests

**Solution Option 1 - Use sudo:**

```bash
sudo go run .
```

**Solution Option 2 - Set Linux capabilities:**

```bash
sudo setcap cap_net_raw=+ep ./ultimate-internet-test
./ultimate-internet-test
```

**Solution Option 3 - Use macOS/Windows:**
Ping works without special permissions on macOS and Windows.

#### Network Timeout Errors

**Symptom:** `context deadline exceeded` or tests hanging

**Cause:** Network timeouts or blocked connections

**Solution:**

```go
cfg := config.New()
cfg.HTTPTimeout = 15 * time.Second      // Increase timeout
cfg.SpeedTestTimeout = 20 * time.Second
result := modules.TestHTTP(url, cfg)
```

Or from command line, modify the code and rebuild.

#### Connection Refused

**Symptom:** Tests fail with "connection refused"

**Cause:** Target server is down or unreachable

**Solution:**

1. Check if target is accessible: `curl https://target.com`
2. Verify firewall settings
3. Use a different target URL

### File System Issues

#### Permission Denied Writing data.json

**Symptom:** `permission denied` when saving results

**Solution:**

```bash
# Check current permissions
ls -la data.json

# Make writable
chmod 644 data.json

# Check directory permissions
ls -ld .

# Ensure directory is writable
chmod 755 .
```

#### data.json File Corrupted

**Symptom:** `invalid character` or JSON unmarshal errors

**Solution:**

```bash
# Backup corrupted file
cp data.json data.json.backup

# Remove corrupted file
rm data.json

# Re-run program to create fresh file
go run .
```

### VPN Detection Issues

#### VPN Detection Always Says "Not Using VPN"

**Symptom:** Even with VPN enabled, shows "Not using VPN"

**Cause:** Local network configuration

**Solution:**

1. Verify the IP detection service is working:

```bash
curl http://checkip.dyndns.org/
```

2 Check your actual external IP:

```bash
curl ifconfig.me
```

3 Verify localhost resolution:

```bash
nslookup localhost
```

#### VPN Detection Service Unreachable

**Symptom:** "could not extract IP address from response"

**Solution:**
Use alternative IP detection service:

```go
result := modules.CheckVPN("http://ifconfig.me")  // Alternative service
```

### Speed Test Issues

#### Speed Test Shows 0 Mbps

**Symptom:** `DownloadMbps: 0.00`

**Cause:** Test URL might be blocked or connection too fast

**Solution:**

1. Use a larger file: `https://example.com/largefile.iso`
2. Increase timeout: `cfg.SpeedTestTimeout = 30 * time.Second`
3. Use different URL with better tracking

#### Speed Test Hangs

**Symptom:** Program hangs during speed test

**Cause:** Large file or slow connection

**Solution:**

1. Increase timeout in config
2. Use smaller/faster file URL
3. Check network connectivity: `ping 8.8.8.8`

### JSON Output Issues

#### data.json Has Empty Results

**Symptom:** Tests ran but results are empty

**Cause:** Test failures or parsing issues

**Solution:**

1. Check logs for error messages
2. Run with verbose output: `go run . 2>&1 | grep -i error`
3. Verify individual tests work

#### data.json Structure Unexpected

**Symptom:** Missing fields or different structure

**Cause:** Version mismatch or incomplete save

**Solution:**

1. Delete existing data.json
2. Run program again
3. Check if all tests completed successfully

### Concurrency Issues

#### Race Conditions in Tests

**Symptom:** Intermittent failures, data inconsistency

**Solution - Run with race detector:**

```bash
go test -race ./...
go run -race .
```

### Debugging

#### Verbose Logging

Add to `main.go`:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
```

#### Print Debug Information

```go
result := modules.TestHTTP(url, cfg)
log.Printf("DEBUG: %#v\n", result)
```

#### Use Delve Debugger

```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug
```

In debugger:

```sh
(dlv) break main.main
(dlv) continue
(dlv) print cfg
(dlv) next
```

### Platform-Specific Issues

#### macOS: Slow Ping Tests

**Symptom:** Ping tests take much longer on macOS

**Solution:**
Default behavior on macOS; increase `PingCount` in config if needed.

#### Windows: Firewall Blocks Tests

**Symptom:** All tests fail or hang

**Solution:**

1. Add program to Windows Defender exception list
2. Check Windows Firewall settings
3. Run as Administrator: `go run .` (right-click, Run as Administrator)

#### Linux: Permission Errors

**Symptom:** Tests fail with permission errors

**Solution:**
Run with appropriate privileges:

```bash
sudo go run .

# Or set capabilities once
go build -o ultimate-internet-test
sudo setcap cap_net_raw=+ep ./ultimate-internet-test
./ultimate-internet-test
```

### Getting Help

If issues persist:

1. [**Check existing issues:**](https://github.com/ehsanghaffar/ultimate-internet-test/issues)
2. **Review logs:** Capture full output with timestamps
3. **Provide system info:**

   ```bash
   go version
   uname -a
   ```

4. **Create detailed report** with:
   - Error messages
   - Steps to reproduce
   - System information
   - Network configuration

### Performance Troubleshooting

#### Program Runs Slowly

**Solution:**

1. Verify network connectivity
2. Increase timeouts (might be retrying failed requests)
3. Check system resources: `top` or `Activity Monitor`
4. Run with profiling:

```bash
go run -cpuprofile=cpu.prof .
go tool pprof cpu.prof
```

#### High Memory Usage

**Solution:**

1. Check for goroutine leaks: `runtime.NumGoroutine()`
2. Ensure result channels are closed
3. Verify mutex locks are released

### Testing Issues

#### Tests Fail Intermittently

**Symptom:** `go test ./...` fails sometimes

**Cause:** Network-dependent tests or race conditions

**Solution:**

```bash
# Run with race detection
go test -race ./...

# Run multiple times
for i in {1..5}; do go test ./...; done
```

## Quick Reference

| Issue                    | Command                     |
| ------------------------ | --------------------------- |
| Update dependencies      | `go mod tidy`               |
| Check for issues         | `go vet ./...`              |
| Format code              | `go fmt ./...`              |
| Run tests                | `go test ./...`             |
| Build release            | `go build -ldflags="-s -w"` |
| Debug with race detector | `go test -race ./...`       |
| Check code coverage      | `go test -cover ./...`      |
