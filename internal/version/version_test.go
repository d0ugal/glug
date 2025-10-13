package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	// Test that we get valid info
	if info.Version == "" {
		t.Error("Version should not be empty")
	}
	if info.Commit == "" {
		t.Error("Commit should not be empty")
	}
	if info.BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
	// Platform is not part of the Info struct
}

func TestString(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		Commit:    "abc123",
		BuildDate: "2025-01-01T00:00:00Z",
		GoVersion: "go1.21.0",
	}

	expected := "glug 1.0.0 (commit: abc123, built: 2025-01-01T00:00:00Z, go: go1.21.0)"
	result := info.String()

	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestStringWithDevVersion(t *testing.T) {
	info := Info{
		Version:   "dev",
		Commit:    "none",
		BuildDate: "unknown",
		GoVersion: "go1.21.0",
	}

	result := info.String()

	// Should contain all the expected parts
	if !contains(result, "glug") {
		t.Error("String should contain 'glug'")
	}
	if !contains(result, "dev") {
		t.Error("String should contain 'dev'")
	}
	if !contains(result, "none") {
		t.Error("String should contain 'none'")
	}
	if !contains(result, "unknown") {
		t.Error("String should contain 'unknown'")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
