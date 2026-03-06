# CLIServiceWatcher 

A simple CLI utility for monitoring the availability of external services. Built with Go.

CLIServiceWatcher performs a parallel check of the availability of web services from a configuration file, measures the response time and generates a report in JSON format.

### Prerequisites
- Go 1.21 or later

### Running
```bash
git clone https://github.com/Mrilki/CLIServicesWatcher.git
go run ./cmd/CLIServicesWatcher