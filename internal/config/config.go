package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultDeadtime is the default folder expiration time
	DefaultDeadtime = 4 * time.Hour

	// DefaultTempDir is the default directory for temporary folders
	DefaultTempDir = "C:\\Temp"

	// RegistryFileName is the name of the registry file
	RegistryFileName = "registry.json"

	// AppName is the application name
	AppName = "Tempi"
)

// GetRegistryPath returns the full path to the registry file
func GetRegistryPath() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback to user profile
		userProfile := os.Getenv("USERPROFILE")
		localAppData = filepath.Join(userProfile, "AppData", "Local")
	}
	return filepath.Join(localAppData, AppName, RegistryFileName)
}

// GetTempDir returns the configured temp directory
// Falls back to DefaultTempDir if not configured
func GetTempDir() string {
	// For now, we always use the default
	// In the future, this could read from a config file
	return DefaultTempDir
}

// ParseDuration parses a duration string (e.g., "2h", "30m")
// Returns the parsed duration or an error
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// FormatDuration formats a duration into a human-readable string
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "expired"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		if minutes > 0 {
			return formatTimeUnit(hours, "h") + formatTimeUnit(minutes, "m")
		}
		return formatTimeUnit(hours, "h")
	}

	if minutes > 0 {
		return formatTimeUnit(minutes, "m")
	}

	seconds := int(d.Seconds())
	return formatTimeUnit(seconds, "s")
}

func formatTimeUnit(value int, unit string) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d%s", value, unit)
}

