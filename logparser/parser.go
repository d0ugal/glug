package logparser

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a parsed log entry
type LogEntry struct {
	Level   string      `json:"level"`
	Time    interface{} `json:"time"`
	Message string      `json:"message"`
	Other   map[string]interface{}
}

// ShouldShowLogLevel determines if a log entry should be shown based on minimum level
func ShouldShowLogLevel(jsonLine, minLevelStr string) (bool, error) {
	var rawLog map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLine), &rawLog); err != nil {
		return true, nil // If we can't parse JSON, show the line
	}

	// Extract level from the log entry
	levelInterface, exists := rawLog["level"]
	if !exists {
		return true, nil // If no level field, show the line
	}

	levelStr, ok := levelInterface.(string)
	if !ok {
		return true, nil // If level is not a string, show the line
	}

	logLevel := parseLogLevel(levelStr)
	minLevel := parseLogLevel(minLevelStr)

	// Show if log level is >= minimum level
	return logLevel >= minLevel, nil
}

// parseLogLevel converts a string to a LogLevel, handling common aliases
func parseLogLevel(levelStr string) LogLevel {
	levelStr = strings.ToUpper(strings.TrimSpace(levelStr))

	switch levelStr {
	case "TRACE", "TRC":
		return LevelTrace
	case "DEBUG", "DBG":
		return LevelDebug
	case "INFO", "INF":
		return LevelInfo
	case "WARN", "WARNING", "WRN":
		return LevelWarn
	case "ERROR", "ERR", "FATAL", "CRIT", "CRITICAL":
		return LevelError
	default:
		// If we don't recognize the level, treat it as INFO
		return LevelInfo
	}
}

// ParseAndFormat parses a JSON log line and returns a formatted colored string
func ParseAndFormat(jsonLine string) (string, error) {
	return ParseAndFormatWithColors(jsonLine, nil)
}

// ParseAndFormatWithColors parses a JSON log line and returns a formatted colored string with custom color rules
func ParseAndFormatWithColors(jsonLine string, customColors map[string]string) (string, error) {
	return ParseAndFormatWithOptions(jsonLine, customColors, false, nil)
}

// ParseAndFormatWithOptions parses a JSON log line with full configuration options
func ParseAndFormatWithOptions(jsonLine string, customColors map[string]string, convertTimestamps bool, timestampFields []string) (string, error) {
	var rawLog map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLine), &rawLog); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	entry := LogEntry{
		Other: make(map[string]interface{}),
	}

	// Extract known fields
	for key, value := range rawLog {
		switch key {
		case "level":
			if str, ok := value.(string); ok {
				entry.Level = str
			}
		case "time":
			entry.Time = value
		case "message":
			if str, ok := value.(string); ok {
				entry.Message = str
			}
		default:
			entry.Other[key] = value
		}
	}

	return formatEntryWithOptions(entry, customColors, convertTimestamps, timestampFields), nil
}

// formatEntry formats a LogEntry into a colored string
func formatEntry(entry LogEntry) string {
	return formatEntryWithColors(entry, nil)
}

// formatEntryWithColors formats a LogEntry into a colored string with custom color rules
func formatEntryWithColors(entry LogEntry, customColors map[string]string) string {
	return formatEntryWithOptions(entry, customColors, false, nil)
}

// formatEntryWithOptions formats a LogEntry with full configuration options
func formatEntryWithOptions(entry LogEntry, customColors map[string]string, convertTimestamps bool, timestampFields []string) string {
	var parts []string

	// Format timestamp
	timeStr := formatTime(entry.Time)
	if timeStr != "" {
		parts = append(parts, color.CyanString(timeStr))
	}

	// Format level with appropriate color
	if entry.Level != "" {
		levelStr := formatLevel(entry.Level)
		parts = append(parts, levelStr)
	}

	// Add message with custom coloring
	if entry.Message != "" {
		messageStr := applyCustomColors(entry.Message, customColors)
		parts = append(parts, messageStr)
	}

	// Add other fields as key=value pairs
	var otherParts []string
	var keys []string
	for key := range entry.Other {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Sort for consistent output

	for _, key := range keys {
		value := entry.Other[key]
		keyStr := color.MagentaString(key)
		
		// Check if this field should be converted to a timestamp
		var convertedValue string
		if convertTimestamps {
			convertedValue = convertTimestampFieldWithConfig(key, value, timestampFields)
		} else {
			convertedValue = fmt.Sprintf("%v", value)
		}
		
		valueStr := applyCustomColors(color.YellowString(convertedValue), customColors)
		otherParts = append(otherParts, fmt.Sprintf("%s=%s", keyStr, valueStr))
	}

	if len(otherParts) > 0 {
		parts = append(parts, strings.Join(otherParts, " "))
	}

	return strings.Join(parts, " ")
}

// formatTime converts various time formats to a readable string
func formatTime(timeVal interface{}) string {
	if timeVal == nil {
		return ""
	}

	switch t := timeVal.(type) {
	case float64:
		// Assume milliseconds if > 1e10, otherwise seconds
		if t > 1e10 {
			return time.Unix(0, int64(t)*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
		}
		return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	case int64:
		// Assume milliseconds if > 1e10, otherwise seconds
		if t > 1e10 {
			return time.Unix(0, t*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
		}
		return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	case string:
		// Try to parse as RFC3339 or other common formats
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
		if parsed, err := time.Parse("2006-01-02T15:04:05", t); err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
		// If parsing fails, return as-is
		return t
	default:
		return fmt.Sprintf("%v", timeVal)
	}
}

// formatLevel returns a colored level string
func formatLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case "ERROR", "ERR":
		return color.RedString(level)
	case "WARN", "WARNING":
		return color.YellowString(level)
	case "INFO":
		return color.GreenString(level)
	case "DEBUG":
		return color.BlueString(level)
	case "TRACE":
		return color.MagentaString(level)
	default:
		return color.WhiteString(level)
	}
}

// applyCustomColors applies custom color rules to a string
func applyCustomColors(text string, customColors map[string]string) string {
	if customColors == nil || len(customColors) == 0 {
		return color.WhiteString(text)
	}

	result := text
	for word, colorName := range customColors {
		if strings.Contains(result, word) {
			coloredWord := getColorFunc(colorName)(word)
			result = strings.ReplaceAll(result, word, coloredWord)
		}
	}

	// If no custom colors were applied, use default white
	if result == text {
		result = color.WhiteString(text)
	}

	return result
}

// getColorFunc returns the appropriate color function based on color name
func getColorFunc(colorName string) func(string) string {
	switch strings.ToLower(colorName) {
	case "red":
		return func(s string) string { return color.RedString(s) }
	case "green":
		return func(s string) string { return color.GreenString(s) }
	case "yellow":
		return func(s string) string { return color.YellowString(s) }
	case "blue":
		return func(s string) string { return color.BlueString(s) }
	case "magenta":
		return func(s string) string { return color.MagentaString(s) }
	case "cyan":
		return func(s string) string { return color.CyanString(s) }
	case "white":
		return func(s string) string { return color.WhiteString(s) }
	default:
		return func(s string) string { return color.WhiteString(s) }
	}
}

// isTimestampField checks if a field name suggests it contains a timestamp
func isTimestampField(fieldName string) bool {
	fieldName = strings.ToLower(fieldName)
	
	// Common timestamp field patterns - be more specific to avoid false positives
	timestampPatterns := []string{
		"time", "timestamp", "ts", "date", "created", "updated", "modified",
		"expires", "expiry", "expire", "validuntil", "valid_until",
		"starttime", "start_time", "endtime", "end_time", "begintime", "begin_time",
		"lastseen", "last_seen", "lastlogin", "last_login", "lastaccess", "last_access",
		"issued", "issuedat", "issued_at", "notbefore", "not_before", "notafter", "not_after",
		"since", "until", "from", "to", "when",
	}
	
	// Check for exact matches or specific patterns
	for _, pattern := range timestampPatterns {
		if fieldName == pattern || strings.HasPrefix(fieldName, pattern+"_") || strings.HasSuffix(fieldName, "_"+pattern) {
			return true
		}
	}
	
	// Special cases for common patterns
	if strings.Contains(fieldName, "time") && !strings.Contains(fieldName, "status") {
		return true
	}
	if strings.Contains(fieldName, "at") && (strings.Contains(fieldName, "time") || strings.Contains(fieldName, "date")) {
		return true
	}
	
	return false
}

// convertTimestampField converts a field value to human-readable date if it looks like a timestamp
func convertTimestampField(fieldName string, value interface{}) string {
	return convertTimestampFieldWithConfig(fieldName, value, nil)
}

// convertTimestampFieldWithConfig converts a field value to human-readable date with custom field configuration
func convertTimestampFieldWithConfig(fieldName string, value interface{}, customFields []string) string {
	// Check if this field should be converted - only if it's in the custom fields list
	shouldConvert := false
	
	for _, field := range customFields {
		if strings.EqualFold(fieldName, field) {
			shouldConvert = true
			break
		}
	}
	
	if !shouldConvert {
		return fmt.Sprintf("%v", value)
	}
	
	// Try to convert the value to a timestamp
	converted := formatTime(value)
	originalStr := fmt.Sprintf("%v", value)
	
	if converted != "" && converted != originalStr {
		// If conversion was successful and different from original, return both
		return fmt.Sprintf("%s (%s)", converted, originalStr)
	}
	
	// If conversion failed or wasn't different, return original
	return originalStr
}
