# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1](https://github.com/d0ugal/glug/compare/v1.1.0...v1.1.1) (2025-10-13)


### Bug Fixes

* add manual trigger to Homebrew update workflow for testing ([2e06d33](https://github.com/d0ugal/glug/commit/2e06d3387a42597c07ca834e660d158a8a19a470))
* use dedicated token for Homebrew tap access ([3aafef1](https://github.com/d0ugal/glug/commit/3aafef1ba8c9a28c992cdb0cfb3349ec0a57c46b))

## [1.1.0](https://github.com/d0ugal/glug/compare/v1.0.1...v1.1.0) (2025-10-13)


### Features

* add automated Homebrew tap update workflow ([9477c64](https://github.com/d0ugal/glug/commit/9477c64ab0d9602ba078dab2eb784604e4bed003))
* add Homebrew tap for glug ([009ea66](https://github.com/d0ugal/glug/commit/009ea660b9299fdd77a3f00cade718c8c9582762))


### Bug Fixes

* remove unused cache.go file ([9f20e82](https://github.com/d0ugal/glug/commit/9f20e82c95fad98be7a77b47ce11de0c2afa283e))
* update golang.org/x/sys to v0.30.0 (compatible with go 1.21) ([5f8d395](https://github.com/d0ugal/glug/commit/5f8d395d01a058a3dd7ab14ceb7bea72c8756bec))

## [1.0.1](https://github.com/d0ugal/glug/compare/v1.0.0...v1.0.1) (2025-10-13)


### Bug Fixes

* update release workflow to use softprops/action-gh-release ([bcaebe4](https://github.com/d0ugal/glug/commit/bcaebe4f2512fefbfdc362132ae93e0dc6b9b0bb))

## 1.0.0 (2025-10-13)


### Features

* add GitHub workflows for CI and releases ([64eb858](https://github.com/d0ugal/glug/commit/64eb858f4717ab7d36d156c9b6258e098d6156db))
* add release management and renovate configuration ([a69a987](https://github.com/d0ugal/glug/commit/a69a9877ca4142d6d35988804ede357f3ab8ffa6))


### Bug Fixes

* lint ([ea42169](https://github.com/d0ugal/glug/commit/ea421693ddefca35b9b59ac485ef1cf877b5a6c9))
* make tests timezone-agnostic by calculating expected times dynamically ([aef0cf7](https://github.com/d0ugal/glug/commit/aef0cf7b1fe07f13195852d62f2853660f351363))
* resolve linting issues ([091228a](https://github.com/d0ugal/glug/commit/091228a1f2018d9807db17bfd2289f5d9d274bd0))
* resolve remaining linting issues ([72d2225](https://github.com/d0ugal/glug/commit/72d222558b007cb048d68e8d3cf3e2991410d035))
* run go mod tidy to fix linting issues ([9152292](https://github.com/d0ugal/glug/commit/91522923689d69ea061b09156c2d59d2b238ae33))
* update go.mod dependencies for compatibility ([6c69f59](https://github.com/d0ugal/glug/commit/6c69f59f4fc0bc226324a25c9aa1ffb96aee2516))
* update tests to use UTC timezone for CI compatibility ([0836681](https://github.com/d0ugal/glug/commit/083668162e9efc455ba54ee84f4225c9a4c1d081))
* update workflows to trigger on master branch ([9b51a3c](https://github.com/d0ugal/glug/commit/9b51a3c12c485557fc09b5ad730f1d3b4565b02e))

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
