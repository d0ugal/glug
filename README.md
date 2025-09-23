# Glug - JSON Log Parser and Colorizer

A simple Go tool that parses JSON log lines and displays them in a colorized, human-readable format.

## Features

- Parses JSON log entries from stdin
- Colorizes output for better readability
- Formats timestamps into human-readable dates
- Displays log levels with appropriate colors
- Shows additional fields as key=value pairs
- Handles various timestamp formats (Unix seconds, milliseconds, RFC3339)
- **Proper signal handling** - Works correctly with Docker containers and other piped commands
- Custom word coloring with CLI flags

## Installation

### Build from source

```bash
go mod tidy
go build -o glug .
```

### Install directly

```bash
go install github.com/yourusername/glug@latest
```

Note: Replace `yourusername` with the actual GitHub username/organization if publishing to a Git repository.

## Usage

### Basic Usage

Pipe JSON log lines to the tool:

```bash
echo '{"level":"debug","program":"synthetic-monitoring-agent","subsystem":"secretstore","time":1749975482337,"caller":"github.com/grafana/synthetic-monitoring-agent/internal/secrets/tenant.go:125","message":"🐛 NewCachedSecretProvider"}' | ./glug
```

Or read from a file:

```bash
cat logs.json | ./glug
```

### Custom Word Coloring

Color specific words using the `--colour` or `--color` flags:

```bash
echo '{"message":"Test PASS and FAIL results"}' | ./glug --colour green:PASS --colour red:FAIL
```

```bash
cat test-logs.json | ./glug --color green:SUCCESS --color red:ERROR --color yellow:WARNING
```

### Level Filtering

Filter logs by minimum level using the `--level` flag:

```bash
# Show only warnings and errors
cat app.log | ./glug --level warning

# Show only errors (and fatal/critical)
docker logs container | ./glug --level error

# Combined with custom colors
tail -f service.log | ./glug --level info --colour green:PASS --colour red:FAIL
```

### Pager Support

Use a pager for better viewing of large log files:

```bash
# Use pager for large log files
cat large-logs.json | ./glug --pager

# Combine with filtering and colors
cat logs.json | ./glug --pager --level error --colour red:ERROR

# Short form
cat logs.json | ./glug -p
```

**Pager behavior:**
- Auto-detects available pagers: `less` (preferred), `more`, or `cat` (fallback)
- Preserves colors and formatting in the pager
- Works with all other glug features (filtering, custom colors, etc.)
- Use `q` to quit the pager, arrow keys to navigate

**Supported levels** (from lowest to highest):
- `trace` (aliases: `trc`)
- `debug` (aliases: `dbg`)
- `info` (aliases: `inf`)
- `warn` / `warning` (aliases: `wrn`)
- `error` (aliases: `err`, `fatal`, `critical`, `crit`)

**Level filtering behavior:**
- `--level warning` shows: WARNING, ERROR, FATAL, CRITICAL
- `--level error` shows: ERROR, FATAL, CRITICAL  
- `--level debug` shows: DEBUG, INFO, WARNING, ERROR, etc.
- Logs without a level field are always shown
- Invalid JSON lines are always shown

### Help

```bash
./glug --help
```

## Docker Compatibility

Glug handles signals properly, making it work seamlessly with Docker containers:

```bash
# This will work correctly - Ctrl+C will stop both glug and the container
docker run --rm your-app | glug --colour green:PASS --colour red:FAIL

# Same with docker-compose
docker-compose logs -f | glug --color green:UP --color red:DOWN
```

When you press `Ctrl+C`, glug will:
1. Receive the interrupt signal (SIGINT)
2. Gracefully shut down and exit
3. Allow the Docker container to also receive the signal and stop properly
4. Prevent "broken pipe" errors

## Output Format

The tool transforms JSON logs into this format:
```
DATE LEVEL MESSAGE key=val key=val
```

Example:
```
2025-06-12 09:31:22 DEBUG 🐛 NewCachedSecretProvider caller=github.com/grafana/synthetic-monitoring-agent/internal/secrets/tenant.go:125 program=synthetic-monitoring-agent subsystem=secretstore
```

## Color Scheme

### Default Colors

- **Time**: Cyan
- **ERROR/ERR**: Red
- **WARN/WARNING**: Yellow  
- **INFO**: Green
- **DEBUG**: Blue
- **TRACE**: Magenta
- **Message**: White
- **Keys**: Magenta
- **Values**: Yellow

### Custom Colors

Use `--colour color:word` or `--color color:word` to color specific words.

**Supported colors:**
- `red`
- `green` 
- `yellow`
- `blue`
- `magenta`
- `cyan`
- `white`

**Examples:**
```bash
# Highlight test results
./glug --colour green:PASS --colour red:FAIL

# Color service statuses  
./glug --color green:UP --color red:DOWN --color yellow:DEGRADED

# Multiple words with same color
./glug --colour red:ERROR --colour red:CRITICAL --colour red:FATAL
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Examples

Debug log:
```json
{"level":"debug","time":1749975482337,"message":"Debug message","component":"auth"}
```
Output: `2025-06-12 09:31:22 DEBUG Debug message component=auth`

Error log:
```json
{"level":"error","time":"2023-01-01T12:00:00Z","message":"Connection failed","error":"timeout","retries":3}
```
Output: `2023-01-01 12:00:00 ERROR Connection failed error=timeout retries=3` 