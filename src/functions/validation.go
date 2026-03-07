package functions

import (
	"os"
	"path/filepath"
	"strings"
)

// SplitAndTrim splits a string by sep and trims spaces from each element.
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// IsAllowedMimeType checks if the MIME type is allowed.
func IsAllowedMimeType(mime string, allowed []string) bool {
	for _, m := range allowed {
		if m == "*" {
			return true
		}
		if m == mime {
			return true
		}
	}
	return false
}

// GetTotalStorageSize returns the total size of files in the storage path.
func GetTotalStorageSize(path string) (int64, error) {
	var total int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
}

// SplitComma splits a comma-separated string into a slice of trimmed strings.
func SplitComma(s string) []string {
	if s == "" {
		return nil
	}
	return SplitAndTrim(s, ",")
}
