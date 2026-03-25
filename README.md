# CLIServiceWatcher

A simple CLI utility for monitoring the availability of external services. Built with Go.

CLIServiceWatcher performs a parallel check of the availability of web services from a configuration file, measures the
response time and generates a report in JSON format.

### Prerequisites

- Go 1.21 or later

### What it does

- Reads a list of targets from a JSON config file
- Supports multiple check types: HTTP, TCP, and DNS
- Checks targets in parallel using a worker pool
- Measures response time (latency) for each check
- Handles graceful shutdown on Ctrl+C (finishes current checks, saves results)
- Saves a detailed JSON report
- Uses idiomatic Go error handling with sentinel errors and error wrapping
- Implements panic recovery in all goroutines for crash protection
- Renders colored output table with go-pretty and fatih/color

### Configuration file format

The program expects a `cfg.json` file in the working directory. The file should contain:

- `timeout` – global timeout in seconds for each check (can be overridden per target)
- `targets` – an array of targets to check, each with:
    - `name` – display name (optional, defaults to URL if empty)
    - `Address` – address to check (for HTTP: full URL, for TCP: `host:port`, for DNS: domain name)
    - `type` – check type: `"http"`, `"tcp"`, or `"dns"`
    - `timeout` – optional custom timeout in seconds (overrides global)

### flag

- `-config`  Path to configuration file. Default: `cfg.json`
- `-output`  Path to output report file. Default: `report.json`
- `-workers` Maximum number of concurrent workers. Default: `10`

### Running

```bash
git clone https://github.com/Mrilki/CLIServicesWatcher.git
go build -o watcher ./cmd/CLIServicesWatcher

# Windows PowerShell
.\watcher.exe -config="cfg.json" -output="report.json" -workers=4

# Linux/Mac
./watcher -config=cfg.json -output=report.json -workers=4
```

