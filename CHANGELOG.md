# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Version information support with `--version` flag
- Makefile for build automation
- Renovate configuration for dependency management
- CHANGELOG.md for release tracking

### Changed
- Enhanced project structure with internal/version package
- Improved build process with version embedding

## [1.0.0] - 2025-01-XX

### Added
- Initial release of glug JSON log parser and colorizer
- Support for parsing JSON log entries from stdin
- Colorized output for better readability
- Timestamp formatting into human-readable dates
- Log level filtering with `--level` flag
- Custom word coloring with `--colour`/`--color` flags
- Pager support (enabled by default)
- Timestamp conversion for specific fields
- Proper signal handling for Docker containers
- Support for various timestamp formats (Unix seconds, milliseconds, RFC3339)
- Level filtering with configurable minimum levels
- Pager auto-detection (less, more, cat)
- Graceful shutdown on SIGINT/SIGTERM

### Features
- **JSON Log Parsing**: Parses JSON log lines and displays them in human-readable format
- **Colorization**: Colorizes output with appropriate colors for different log levels
- **Timestamp Conversion**: Converts Unix timestamps to human-readable dates
- **Level Filtering**: Filter logs by minimum level (trace, debug, info, warn, error)
- **Custom Colors**: Color specific words using CLI flags
- **Pager Support**: Automatic pager detection and usage for better viewing
- **Docker Compatibility**: Proper signal handling for containerized environments
- **Flexible Input**: Works with stdin, files, and piped commands

### Supported Log Levels
- `trace` (aliases: `trc`)
- `debug` (aliases: `dbg`)
- `info` (aliases: `inf`)
- `warn` / `warning` (aliases: `wrn`)
- `error` (aliases: `err`, `fatal`, `critical`, `crit`)

### Supported Colors
- `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`

### Usage Examples
```bash
# Basic usage
echo '{"level":"debug","message":"Test"}' | ./glug

# With custom colors
echo '{"message":"Test PASS and FAIL"}' | ./glug --colour green:PASS --colour red:FAIL

# With level filtering
cat logs.json | ./glug --level warning

# With timestamp conversion
cat logs.json | ./glug --convert-timestamps validUntil,expires

# Disable pager
echo '{"message":"Quick output"}' | ./glug --no-pager
```
