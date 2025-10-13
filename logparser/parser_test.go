package logparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseAndFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		contains []string // Substrings that should be present in output
	}{
		{
			name:  "basic log entry",
			input: `{"level":"debug","program":"synthetic-monitoring-agent","subsystem":"secretstore","time":1749975482337,"caller":"github.com/grafana/synthetic-monitoring-agent/internal/secrets/tenant.go:125","message":"🐛 NewCachedSecretProvider"}`,
			contains: []string{
				"DEBUG",
				"🐛 NewCachedSecretProvider",
				"program=synthetic-monitoring-agent",
				"subsystem=secretstore",
			},
		},
		{
			name:  "info level",
			input: `{"level":"info","time":1609459200,"message":"Server started"}`,
			contains: []string{
				"INFO",
				"Server started",
			},
		},
		{
			name:  "error level",
			input: `{"level":"error","time":1609459200,"message":"Connection failed","error":"timeout"}`,
			contains: []string{
				"ERROR",
				"Connection failed",
				"error=timeout",
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
		{
			name:  "minimal log",
			input: `{"message":"test"}`,
			contains: []string{
				"test",
			},
		},
		{
			name:  "string time format",
			input: `{"level":"warn","time":"2023-01-01T12:00:00Z","message":"Warning message"}`,
			contains: []string{
				"WARN",
				"Warning message",
				"2023-01-01 12:00:00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAndFormat(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseAndFormat() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseAndFormat() unexpected error: %v", err)
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("ParseAndFormat() result missing expected substring %q\nGot: %s", substr, result)
				}
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil time",
			input:    nil,
			expected: "",
		},
		{
			name:     "millisecond timestamp",
			input:    float64(1749975482337),
			expected: "2025-06-15 08:18:02",
		},
		{
			name:     "second timestamp",
			input:    float64(1609459200),
			expected: "2021-01-01 00:00:00",
		},
		{
			name:     "int64 millisecond timestamp",
			input:    int64(1749975482337),
			expected: "2025-06-15 08:18:02",
		},
		{
			name:     "RFC3339 string",
			input:    "2023-01-01T12:00:00Z",
			expected: "2023-01-01 12:00:00",
		},
		{
			name:     "unparseable string",
			input:    "invalid-time",
			expected: "invalid-time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTime(tt.input)
			if result != tt.expected {
				t.Errorf("formatTime() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"debug", "debug", "DEBUG"},
		{"info", "info", "INFO"},
		{"warn", "warn", "WARN"},
		{"error", "error", "ERROR"},
		{"unknown", "custom", "CUSTOM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLevel(tt.input)
			// Just check that the level text is present (ignoring ANSI color codes)
			if !strings.Contains(result, tt.want) {
				t.Errorf("formatLevel() = %q, should contain %q", result, tt.want)
			}
		})
	}
}

func TestLogEntryExtractionAndSorting(t *testing.T) {
	input := `{"level":"info","z_last":"last","a_first":"first","message":"test message","time":1609459200}`

	result, err := ParseAndFormat(input)
	if err != nil {
		t.Fatalf("ParseAndFormat() error: %v", err)
	}

	// Check that keys are sorted (a_first should come before z_last)
	firstPos := strings.Index(result, "a_first")
	lastPos := strings.Index(result, "z_last")

	if firstPos == -1 || lastPos == -1 {
		t.Errorf("Expected both a_first and z_last in output: %s", result)
	}

	if firstPos > lastPos {
		t.Errorf("Expected a_first to come before z_last in output: %s", result)
	}
}

func TestParseAndFormatWithColors(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		customColors map[string]string
		contains     []string
		notContains  []string
	}{
		{
			name:         "no custom colors",
			input:        `{"message":"Test PASS and FAIL"}`,
			customColors: nil,
			contains:     []string{"Test PASS and FAIL"},
		},
		{
			name:         "custom color for PASS",
			input:        `{"message":"Test PASS result"}`,
			customColors: map[string]string{"PASS": "green"},
			contains:     []string{"Test", "result"},
		},
		{
			name:  "multiple custom colors",
			input: `{"message":"Test PASS and FAIL results"}`,
			customColors: map[string]string{
				"PASS": "green",
				"FAIL": "red",
			},
			contains: []string{"Test", "and", "results"},
		},
		{
			name:  "custom colors in values",
			input: `{"message":"Testing","status":"PASS","error":"FAIL"}`,
			customColors: map[string]string{
				"PASS": "green",
				"FAIL": "red",
			},
			contains: []string{"Testing"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAndFormatWithColors(tt.input, tt.customColors)
			if err != nil {
				t.Errorf("ParseAndFormatWithColors() error: %v", err)
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("ParseAndFormatWithColors() result missing expected substring %q\nGot: %s", substr, result)
				}
			}

			for _, substr := range tt.notContains {
				if strings.Contains(result, substr) {
					t.Errorf("ParseAndFormatWithColors() result contains unexpected substring %q\nGot: %s", substr, result)
				}
			}
		})
	}
}

func TestGetColorFunc(t *testing.T) {
	tests := []struct {
		colorName string
		testWord  string
	}{
		{"red", "ERROR"},
		{"green", "PASS"},
		{"yellow", "WARN"},
		{"blue", "INFO"},
		{"magenta", "TRACE"},
		{"cyan", "TIME"},
		{"white", "DEFAULT"},
		{"invalid", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.colorName, func(t *testing.T) {
			colorFunc := getColorFunc(tt.colorName)
			result := colorFunc(tt.testWord)

			// Just check that the function returns something and includes the original word
			if !strings.Contains(result, tt.testWord) {
				t.Errorf("getColorFunc(%q)(%q) should contain the original word, got: %s", tt.colorName, tt.testWord, result)
			}
		})
	}
}

func TestShouldShowLogLevel(t *testing.T) {
	tests := []struct {
		name       string
		jsonLine   string
		minLevel   string
		shouldShow bool
	}{
		{
			name:       "error level with warning filter",
			jsonLine:   `{"level":"error","message":"Error occurred"}`,
			minLevel:   "warning",
			shouldShow: true,
		},
		{
			name:       "info level with warning filter",
			jsonLine:   `{"level":"info","message":"Info message"}`,
			minLevel:   "warning",
			shouldShow: false,
		},
		{
			name:       "warning level with warning filter",
			jsonLine:   `{"level":"warning","message":"Warning message"}`,
			minLevel:   "warning",
			shouldShow: true,
		},
		{
			name:       "debug level with error filter",
			jsonLine:   `{"level":"debug","message":"Debug message"}`,
			minLevel:   "error",
			shouldShow: false,
		},
		{
			name:       "error level with debug filter",
			jsonLine:   `{"level":"error","message":"Error message"}`,
			minLevel:   "debug",
			shouldShow: true,
		},
		{
			name:       "no level field",
			jsonLine:   `{"message":"No level field"}`,
			minLevel:   "error",
			shouldShow: true, // Should show when no level field
		},
		{
			name:       "invalid JSON",
			jsonLine:   `{invalid json}`,
			minLevel:   "error",
			shouldShow: true, // Should show when JSON is invalid
		},
		{
			name:       "trace level with trace filter",
			jsonLine:   `{"level":"trace","message":"Trace message"}`,
			minLevel:   "trace",
			shouldShow: true,
		},
		{
			name:       "warn alias with warning filter",
			jsonLine:   `{"level":"warn","message":"Warning message"}`,
			minLevel:   "warning",
			shouldShow: true,
		},
		{
			name:       "fatal level with error filter",
			jsonLine:   `{"level":"fatal","message":"Fatal error"}`,
			minLevel:   "error",
			shouldShow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ShouldShowLogLevel(tt.jsonLine, tt.minLevel)
			if err != nil {
				t.Errorf("ShouldShowLogLevel() error: %v", err)
				return
			}

			if result != tt.shouldShow {
				t.Errorf("ShouldShowLogLevel() = %v, want %v", result, tt.shouldShow)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"trace", LevelTrace},
		{"TRACE", LevelTrace},
		{"TRC", LevelTrace},
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"DBG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"INF", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"WARNING", LevelWarn},
		{"WRN", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"ERR", LevelError},
		{"fatal", LevelError},
		{"FATAL", LevelError},
		{"critical", LevelError},
		{"CRITICAL", LevelError},
		{"CRIT", LevelError},
		{"unknown", LevelInfo}, // Default to INFO for unknown levels
		{"", LevelInfo},        // Default to INFO for empty string
		{" WARN ", LevelWarn},  // Should handle whitespace
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelTrace, "TRACE"},
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LogLevel(999), "UNKNOWN"}, // Invalid level
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("LogLevel(%d).String() = %q, want %q", tt.level, result, tt.expected)
			}
		})
	}
}

func TestLevelHierarchy(t *testing.T) {
	// Test that level hierarchy is correct
	if LevelTrace >= LevelDebug {
		t.Error("TRACE should be < DEBUG")
	}
	if LevelDebug >= LevelInfo {
		t.Error("DEBUG should be < INFO")
	}
	if LevelInfo >= LevelWarn {
		t.Error("INFO should be < WARN")
	}
	if LevelWarn >= LevelError {
		t.Error("WARN should be < ERROR")
	}
}

func TestIsTimestampField(t *testing.T) {
	tests := []struct {
		fieldName string
		expected  bool
	}{
		{"time", true},
		{"timestamp", true},
		{"ts", true},
		{"date", true},
		{"created", true},
		{"updated", true},
		{"modified", true},
		{"expires", true},
		{"expiry", true},
		{"validUntil", true},
		{"valid_until", true},
		{"startTime", true},
		{"start_time", true},
		{"endTime", true},
		{"end_time", true},
		{"lastSeen", true},
		{"last_seen", true},
		{"lastLogin", true},
		{"last_login", true},
		{"issued", true},
		{"issuedAt", true},
		{"issued_at", true},
		{"notBefore", true},
		{"not_before", true},
		{"notAfter", true},
		{"not_after", true},
		{"since", true},
		{"until", true},
		{"from", true},
		{"to", true},
		{"when", true},
		{"at", false}, // Now more restrictive - "at" alone is not considered a timestamp field
		{"on", false}, // Now more restrictive - "on" alone is not considered a timestamp field
		{"message", false},
		{"level", false},
		{"user", false},
		{"id", false},
		{"name", false},
		{"status", false},
		{"count", false},
		{"value", false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			result := isTimestampField(tt.fieldName)
			if result != tt.expected {
				t.Errorf("isTimestampField(%q) = %v, want %v", tt.fieldName, result, tt.expected)
			}
		})
	}
}

func TestConvertTimestampField(t *testing.T) {
	tests := []struct {
		fieldName string
		value     interface{}
		expected  string
	}{
		// Without custom fields, no conversion should happen
		{"validUntil", int64(1760134416629), "1760134416629"},
		{"expires", float64(1609459200), "1.6094592e+09"},
		{"created", "2023-01-01T12:00:00Z", "2023-01-01T12:00:00Z"},
		{"timestamp", int64(1749975482337), "1749975482337"},

		// Non-timestamp fields should not be converted
		{"message", "test message", "test message"},
		{"level", "info", "info"},
		{"user", "john", "john"},
		{"count", 42, "42"},

		// Timestamp fields with invalid values
		{"expires", "invalid-date", "invalid-date"},
		{"timestamp", nil, "<nil>"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%v", tt.fieldName, tt.value), func(t *testing.T) {
			result := convertTimestampField(tt.fieldName, tt.value)
			if result != tt.expected {
				t.Errorf("convertTimestampField(%q, %v) = %q, want %q", tt.fieldName, tt.value, result, tt.expected)
			}
		})
	}
}

func TestConvertTimestampFieldWithConfig(t *testing.T) {
	tests := []struct {
		fieldName    string
		value        interface{}
		customFields []string
		expected     string
	}{
		// With custom fields specified
		{"validUntil", int64(1760134416629), []string{"validUntil"}, "2025-10-10 22:13:36 (1760134416629)"},
		{"expires", float64(1609459200), []string{"expires"}, "2021-01-01 00:00:00 (1.6094592e+09)"},
		{"created", "2023-01-01T12:00:00Z", []string{"created"}, "2023-01-01 12:00:00 (2023-01-01T12:00:00Z)"},
		{"timestamp", int64(1749975482337), []string{"timestamp"}, "2025-06-15 08:18:02 (1749975482337)"},

		// Fields not in custom list should not be converted
		{"validUntil", int64(1760134416629), []string{"expires"}, "1760134416629"},
		{"expires", float64(1609459200), []string{"validUntil"}, "1.6094592e+09"},

		// Case insensitive matching
		{"ValidUntil", int64(1760134416629), []string{"validuntil"}, "2025-10-10 22:13:36 (1760134416629)"},
		{"EXPIRES", float64(1609459200), []string{"expires"}, "2021-01-01 00:00:00 (1.6094592e+09)"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%v_%v", tt.fieldName, tt.value, tt.customFields), func(t *testing.T) {
			result := convertTimestampFieldWithConfig(tt.fieldName, tt.value, tt.customFields)
			if result != tt.expected {
				t.Errorf("convertTimestampFieldWithConfig(%q, %v, %v) = %q, want %q", tt.fieldName, tt.value, tt.customFields, result, tt.expected)
			}
		})
	}
}

func TestTimestampFieldConversionInLogs(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		convertTimestamps bool
		timestampFields   []string
		contains          []string
		notContains       []string
	}{
		{
			name:              "no timestamp conversion by default",
			input:             `{"level":"info","message":"Token created","validUntil":1760134416629}`,
			convertTimestamps: false,
			contains: []string{
				"Token created",
				"validUntil=1.760134416629e+12",
			},
			notContains: []string{
				"2025-10-10",
			},
		},
		{
			name:              "validUntil timestamp conversion with explicit field",
			input:             `{"level":"info","message":"Token created","validUntil":1760134416629}`,
			convertTimestamps: true,
			timestampFields:   []string{"validUntil"},
			contains: []string{
				"Token created",
				"validUntil=",
				"2025-10-10 22:13:36",
				"(1.760134416629e+12)",
			},
		},
		{
			name:              "expires timestamp conversion",
			input:             `{"level":"warn","message":"Session expires soon","expires":1609459200}`,
			convertTimestamps: true,
			timestampFields:   []string{"expires"},
			contains: []string{
				"Session expires soon",
				"expires=",
				"2021-01-01 00:00:00",
				"(1.6094592e+09)",
			},
		},
		{
			name:              "multiple timestamp fields",
			input:             `{"level":"info","message":"User action","created":1609459200,"expires":1760134416629}`,
			convertTimestamps: true,
			timestampFields:   []string{"created", "expires"},
			contains: []string{
				"User action",
				"created=",
				"expires=",
				"2021-01-01 00:00:00",
				"2025-10-10 22:13:36",
			},
		},
		{
			name:              "non-timestamp fields unchanged",
			input:             `{"level":"info","message":"Test","user":"john","count":42}`,
			convertTimestamps: true,
			timestampFields:   []string{"validUntil"},
			contains: []string{
				"Test",
				"user=john",
				"count=42",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAndFormatWithOptions(tt.input, nil, tt.convertTimestamps, tt.timestampFields)
			if err != nil {
				t.Errorf("ParseAndFormatWithOptions() error: %v", err)
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("ParseAndFormatWithOptions() result missing expected substring %q\nGot: %s", substr, result)
				}
			}

			for _, substr := range tt.notContains {
				if strings.Contains(result, substr) {
					t.Errorf("ParseAndFormatWithOptions() result contains unexpected substring %q\nGot: %s", substr, result)
				}
			}
		})
	}
}
