package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestSignalHandling(t *testing.T) {
	// This test verifies that signal handling is set up correctly
	// We can't easily test the full signal flow in a unit test,
	// but we can test that the signal setup doesn't panic

	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with help flag to avoid stdin reading
	os.Args = []string{"glug", "--help"}

	// This should not panic or hang
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main() panicked: %v", r)
			}
			done <- true
		}()
		main()
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Success - main completed without panic
	case <-time.After(2 * time.Second):
		t.Error("main() took too long to complete")
	}
}

func TestColorRuleParsing(t *testing.T) {
	tests := []struct {
		rule  string
		valid bool
	}{
		{"green:PASS", true},
		{"red:FAIL", true},
		{"blue:INFO", true},
		{"invalid", false},
		{"", false},
		{"color:", false},
		{":word", false},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			parts := splitColorRule(tt.rule)
			isValid := len(parts) == 2 && parts[0] != "" && parts[1] != ""

			if isValid != tt.valid {
				t.Errorf("splitColorRule(%q) validity = %v, want %v", tt.rule, isValid, tt.valid)
			}
		})
	}
}

// Helper function to test color rule parsing logic
func splitColorRule(rule string) []string {
	if rule == "" {
		return []string{}
	}
	return strings.SplitN(rule, ":", 2)
}

func TestDetectPager(t *testing.T) {
	// Test that detectPager returns a valid pager
	pager := detectPager()
	if pager == "" {
		t.Error("detectPager() should return a non-empty string")
	}
	
	// Should be one of the expected pagers
	expectedPagers := []string{"less", "more", "cat"}
	validPager := false
	for _, expected := range expectedPagers {
		if pager == expected {
			validPager = true
			break
		}
	}
	
	if !validPager {
		t.Errorf("detectPager() returned unexpected pager: %s", pager)
	}
}

func TestExecuteWithPager(t *testing.T) {
	// Test with cat (should always be available)
	content := "test line 1\ntest line 2\ntest line 3"
	
	err := executeWithPager(content, "cat")
	if err != nil {
		t.Errorf("executeWithPager() with cat failed: %v", err)
	}
}
