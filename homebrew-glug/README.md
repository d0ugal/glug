# Homebrew Tap for Glug

This is a Homebrew tap for [Glug](https://github.com/d0ugal/glug) - a JSON log parser and colorizer.

## Installation

```bash
# Add this tap
brew tap d0ugal/glug

# Install glug
brew install glug
```

## Usage

```bash
# Parse JSON logs with colorization
echo '{"level":"info","message":"Hello World"}' | glug

# Use custom colors
echo '{"message":"Test PASS"}' | glug --colour green:PASS

# Filter by log level
echo '{"level":"error","message":"Something went wrong"}' | glug --level error
```

## What is Glug?

Glug is a command-line tool that parses JSON log files and outputs them with beautiful colorization and formatting. It's perfect for:

- Docker logs: `docker logs container | glug`
- Application logs: `cat app.log | glug`
- Real-time monitoring: `tail -f logs.json | glug`

## Features

- 🎨 **Colorized output** with customizable word highlighting
- 📊 **Log level filtering** (trace, debug, info, warn, error)
- ⏰ **Timestamp conversion** for specified fields
- 📄 **Pager support** for large outputs
- 🐳 **Docker-friendly** with proper signal handling

## License

MIT License - see the [main repository](https://github.com/d0ugal/glug) for details.
