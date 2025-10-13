package processor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// OutputHandler manages output to either stdout or pager
type OutputHandler struct {
	usePager    bool
	outputLines []string
}

// NewOutputHandler creates a new output handler
func NewOutputHandler(usePager bool) *OutputHandler {
	return &OutputHandler{
		usePager:    usePager,
		outputLines: make([]string, 0),
	}
}

// AddLine adds a line to the output buffer
func (oh *OutputHandler) AddLine(line string) {
	if oh.usePager {
		oh.outputLines = append(oh.outputLines, line)
	} else {
		fmt.Println(line)
	}
}

// Flush outputs all buffered lines
func (oh *OutputHandler) Flush() error {
	if !oh.usePager {
		return nil
	}

	pagerName := detectPager()
	content := strings.Join(oh.outputLines, "\n")

	if err := executeWithPager(content, pagerName); err != nil {
		return fmt.Errorf("error running pager: %v", err)
	}

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
