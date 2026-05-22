# logslice

Fast log parser that extracts time-range slices from large compressed log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

## Usage

```
logslice [flags] <logfile>

Flags:
  --from   string   Start of time range (e.g. "2024-01-15T08:00:00")
  --to     string   End of time range   (e.g. "2024-01-15T09:00:00")
  --format string   Log timestamp format (default: RFC3339)
  --out    string   Output file (default: stdout)
```

### Example

Extract all log entries from a compressed log file between 8am and 9am:

```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" /var/log/app.log.gz
```

Pipe the output to another tool:

```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log.gz | grep "ERROR"
```

## Features

- Supports `.gz` and `.zst` compressed log files
- Binary search on indexed timestamps for fast seeking
- Streams output — works on files of any size
- Handles multiple timestamp formats

## Requirements

- Go 1.21 or later

## License

MIT — see [LICENSE](LICENSE) for details.