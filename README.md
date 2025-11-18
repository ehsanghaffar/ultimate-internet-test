# All-in-one Internet Test

A lightweight Go application that performs a variety of internet and network checks (local network access, international access, ping, speed tests, interface listing, VPN detection, etc.). Designed to be easy to run locally for quick diagnostics.

## Features

- Test local network access (LAN)
- Test international network access (WAN)
- Run basic ping checks
- Measure internet speed (download/upload/latency) where applicable
- List available network interfaces
- Basic VPN detection checks (where possible)
- Minimal dependencies; written in Go

## Prerequisites

- Go (1.18+ recommended)
- Network access for tests
- Some checks may require elevated privileges (for ICMP/ping or interface inspection)

## Quick start

Clone and run:

```bash
git clone https://github.com/ghaffariidev/check-internet-tests.git
cd check-internet-tests

# Run directly
go run .

# Or build and run binary
go build -o check-internet-tests .
./check-internet-tests
```

If the program supports a help flag, run:

```bash
./check-internet-tests --help
```

## Configuration

- The repository includes a data.json file that contains default endpoints or settings used by the tests. Edit it to change targets or defaults.
- Additional modules live under modules/ and helper functions under utils/.

## Usage

- Run the program (go run . or the built binary) and follow the on-screen output.
- For automation, redirect output to a file:

```bash
./check-internet-tests > results.txt
```

- For tests that require low-level network access (ICMP/ping, interface queries), run with appropriate privileges if you encounter permission errors.

## Development

- Implemented in Go. Use `go build`, `go run` and `go test` as usual.
- To add features, create modules under modules/ and helpers under utils/.
- Run `go mod tidy` after adding dependencies.

## Contributing

Contributions are welcome!

1. Fork the repository.
2. Create a branch for your feature or fix.
3. Submit a pull request with a clear description and rationale.
4. Add or update tests when applicable.

## Authors and contact

- Original author: @ehsanghaffar (https://github.com/ehsanghaffar)
- For feedback: ghafari.5000@gmail.com

## Documentation

Project documentation and more details may be available at: https://ehsanghaffarii.ir

## License

This project is licensed under the MIT License — see https://choosealicense.com/licenses/mit/ for details.

## Troubleshooting

- If tests fail, verify local firewall settings and outbound internet access.
- Speed test results vary due to ISP, server location, and network load.
- If you need help understanding output, run with verbose/debug flags if available or open an issue including sample output and environment details.