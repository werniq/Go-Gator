package cmd

import (
	"fmt"
	"strings"
	"time"
)

// validateDate checks if the date string is in the correct format YYYY-MM-DD
func validateDate(dateStr string) error {
	if dateStr == "" {
		return nil
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format for %s, expected YYYY-MM-DD", dateStr)
	}

	return nil
}

// validateSources checks if the provided sources are within the supported list
func validateSources(sources string) error {
	supportedSources := []string{"abc", "bbc", "nbc", "usatoday", "washingtontimes", "all"}

	for _, source := range strings.Split(sources, ",") {
		if !contains(supportedSources, source) {
			return fmt.Errorf("unsupported source: %s. Supported sources are: %v", source, supportedSources)
		}
	}
	return nil
}

// contains checks if a slice contains a given string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
