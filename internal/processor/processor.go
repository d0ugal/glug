package processor

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/dougalmatthews/glug/logparser"
)

// LogProcessor handles the processing of log input
type LogProcessor struct {
	config       *Config
	customColors map[string]string
	output       *OutputHandler
}

// Config represents the application configuration
type Config struct {
	MinLevel           string
	UsePager           bool
	ConvertTimestamps  bool
	TimestampFieldList []string
}

// NewLogProcessor creates a new log processor
func NewLogProcessor(config *Config, customColors map[string]string) *LogProcessor {
	return &LogProcessor{
		config:       config,
		customColors: customColors,
		output:       NewOutputHandler(config.UsePager),
	}
}

// Process reads from stdin and processes log entries
func (lp *LogProcessor) Process(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// Check if we should exit due to signal
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		formatted, err := lp.processLine(line)
		if err != nil {
			// If parsing fails, just print the original line
			lp.output.AddLine(line)
			continue
		}

		// Apply level filtering if specified
		if lp.config.MinLevel != "" {
			shouldShow, err := logparser.ShouldShowLogLevel(line, lp.config.MinLevel)
			if err != nil {
				// If level parsing fails, show the line (fail open)
				lp.output.AddLine(formatted)
				continue
			}
			if !shouldShow {
				continue
			}
		}

		lp.output.AddLine(formatted)
	}

	if err := scanner.Err(); err != nil {
		// Don't report error if context was cancelled (user pressed Ctrl+C)
		select {
		case <-ctx.Done():
			return nil
		default:
			return fmt.Errorf("error reading input: %v", err)
		}
	}

	// Flush output
	return lp.output.Flush()
}

// processLine processes a single log line
func (lp *LogProcessor) processLine(line string) (string, error) {
	return logparser.ParseAndFormatWithOptions(
		line,
		lp.customColors,
		lp.config.ConvertTimestamps,
		lp.config.TimestampFieldList,
	)
}
