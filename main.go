package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dougalmatthews/glug/logparser"
)

type colorFlags []string

func (c *colorFlags) String() string {
	return strings.Join(*c, ", ")
}

func (c *colorFlags) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func main() {
	var colorRules colorFlags
	flag.Var(&colorRules, "colour", "Color specific words (format: color:word, e.g., green:PASS)")
	flag.Var(&colorRules, "color", "Color specific words (format: color:word, e.g., green:PASS)")

	var minLevel string
	flag.StringVar(&minLevel, "level", "", "Minimum log level to show (trace, debug, info, warn/warning, error)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&help, "h", false, "Show help message")

	flag.Parse()

	if help {
		fmt.Fprintf(os.Stderr, "Glug - JSON Log Parser and Colorizer\n\n")
		fmt.Fprintf(os.Stderr, "Usage: glug [options] < logfile.json\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  echo '{\"message\":\"Test PASS\"}' | glug --colour green:PASS\n")
		fmt.Fprintf(os.Stderr, "  cat logs.json | glug --colour green:PASS --colour red:FAIL\n")
		fmt.Fprintf(os.Stderr, "  docker logs container | glug --level warning --color red:ERROR\n")
		fmt.Fprintf(os.Stderr, "\nSupported colors: red, green, yellow, blue, magenta, cyan, white\n")
		fmt.Fprintf(os.Stderr, "Supported levels: trace, debug, info, warn/warning, error\n")
		return
	}

	// Parse color rules
	customColors := make(map[string]string)
	for _, rule := range colorRules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid color rule format: %s (expected color:word)\n", rule)
			os.Exit(1)
		}
		color, word := parts[0], parts[1]
		customColors[word] = color
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// Check if we should exit due to signal
		select {
		case <-ctx.Done():
			return
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		formatted, err := logparser.ParseAndFormatWithColors(line, customColors)
		if err != nil {
			// If parsing fails, just print the original line
			fmt.Println(line)
			continue
		}

		// Apply level filtering if specified
		if minLevel != "" {
			shouldShow, err := logparser.ShouldShowLogLevel(line, minLevel)
			if err != nil {
				// If level parsing fails, show the line (fail open)
				fmt.Println(formatted)
				continue
			}
			if !shouldShow {
				continue
			}
		}

		fmt.Println(formatted)
	}

	if err := scanner.Err(); err != nil {
		// Don't report error if context was cancelled (user pressed Ctrl+C)
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	}
}
