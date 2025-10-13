package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
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

// detectPager finds the best available pager
func detectPager() string {
	// Check for common pagers in order of preference
	pagers := []string{"less", "more", "cat"}

	for _, pager := range pagers {
		if _, err := exec.LookPath(pager); err == nil {
			return pager
		}
	}

	// Fallback to cat if no pager is found
	return "cat"
}

// executeWithPager runs the pager with the given content
func executeWithPager(content string, pagerName string) error {
	var cmd *exec.Cmd

	// Configure pager with appropriate flags for color support
	switch pagerName {
	case "less":
		// -R: enable raw control characters (colors)
		// -X: don't clear screen on exit
		// -F: quit if one screen
		cmd = exec.Command("less", "-R", "-X", "-F")
	case "more":
		// more doesn't need special flags for colors
		cmd = exec.Command("more")
	case "cat":
		// cat just outputs everything
		cmd = exec.Command("cat")
	default:
		cmd = exec.Command(pagerName)
	}

	// Set up stdin for the pager
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the pager
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start pager %s: %v", pagerName, err)
	}

	// Wait for the pager to complete
	return cmd.Wait()
}

func main() {
	var colorRules colorFlags
	flag.Var(&colorRules, "colour", "Color specific words (format: color:word, e.g., green:PASS)")
	flag.Var(&colorRules, "color", "Color specific words (format: color:word, e.g., green:PASS)")

	var minLevel string
	flag.StringVar(&minLevel, "level", "", "Minimum log level to show (trace, debug, info, warn/warning, error)")

	var usePager bool
	flag.BoolVar(&usePager, "pager", true, "Use pager for output (auto-detects less/more) [default: true]")
	flag.BoolVar(&usePager, "p", true, "Use pager for output (auto-detects less/more) [default: true]")
	
	var noPager bool
	flag.BoolVar(&noPager, "no-pager", false, "Disable pager (output directly to stdout)")
	flag.BoolVar(&noPager, "n", false, "Disable pager (output directly to stdout)")

	var timestampFields string
	flag.StringVar(&timestampFields, "convert-timestamps", "", "Comma-separated list of field names to convert as timestamps")
	flag.StringVar(&timestampFields, "t", "", "Comma-separated list of field names to convert as timestamps")

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
		fmt.Fprintf(os.Stderr, "  cat large-logs.json | glug --level error\n")
		fmt.Fprintf(os.Stderr, "  echo '{\"message\":\"Quick output\"}' | glug --no-pager\n")
		fmt.Fprintf(os.Stderr, "  cat logs.json | glug --convert-timestamps validUntil,expires\n")
		fmt.Fprintf(os.Stderr, "  cat logs.json | glug --convert-timestamps created,updated\n")
		fmt.Fprintf(os.Stderr, "\nSupported colors: red, green, yellow, blue, magenta, cyan, white\n")
		fmt.Fprintf(os.Stderr, "Supported levels: trace, debug, info, warn/warning, error\n")
		fmt.Fprintf(os.Stderr, "Pager: Enabled by default, use --no-pager to disable\n")
		fmt.Fprintf(os.Stderr, "Timestamps: Use --convert-timestamps to specify which fields to convert\n")
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

	// Handle pager logic: default to true, but can be disabled with --no-pager
	if noPager {
		usePager = false
	}

	// Parse timestamp fields - conversion is enabled only if fields are specified
	var timestampFieldList []string
	var convertTimestamps bool
	if timestampFields != "" {
		convertTimestamps = true
		timestampFieldList = strings.Split(timestampFields, ",")
		// Trim whitespace from each field name
		for i, field := range timestampFieldList {
			timestampFieldList[i] = strings.TrimSpace(field)
		}
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

	// Collect output if using pager
	var outputLines []string

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

		formatted, err := logparser.ParseAndFormatWithOptions(line, customColors, convertTimestamps, timestampFieldList)
		if err != nil {
			// If parsing fails, just print the original line
			if usePager {
				outputLines = append(outputLines, line)
			} else {
				fmt.Println(line)
			}
			continue
		}

		// Apply level filtering if specified
		if minLevel != "" {
			shouldShow, err := logparser.ShouldShowLogLevel(line, minLevel)
			if err != nil {
				// If level parsing fails, show the line (fail open)
				if usePager {
					outputLines = append(outputLines, formatted)
				} else {
					fmt.Println(formatted)
				}
				continue
			}
			if !shouldShow {
				continue
			}
		}

		if usePager {
			outputLines = append(outputLines, formatted)
		} else {
			fmt.Println(formatted)
		}
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

	// If using pager, execute it with collected output
	if usePager {
		pagerName := detectPager()
		content := strings.Join(outputLines, "\n")

		if err := executeWithPager(content, pagerName); err != nil {
			fmt.Fprintf(os.Stderr, "Error running pager: %v\n", err)
			os.Exit(1)
		}
	}
}
