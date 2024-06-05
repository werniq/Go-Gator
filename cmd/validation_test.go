package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validationDate(t *testing.T) {
	testCases := []struct {
		Input          string
		ExpectedOutput error
	}{
		{
			Input:          "2024-05-12",
			ExpectedOutput: nil,
		},
		{
			Input:          "2024-13-05",
			ExpectedOutput: fmt.Errorf("invalid date format for %s, expected YYYY-MM-DD", "2024-13-05"),
		},
		{
			Input:          "2024-05-92",
			ExpectedOutput: fmt.Errorf("invalid date format for %s, expected YYYY-MM-DD", "2024-05-92"),
		},
	}

	for _, testCase := range testCases {
		err := validateDate(testCase.Input)
		assert.Equal(t, err, testCase.ExpectedOutput)
	}
}

func Test_validateSources(t *testing.T) {
	testCases := []struct {
		Input          []string
		ExpectedOutput error
	}{
		{
			Input: []string{"abc", "aaaa", "bbc"},
			ExpectedOutput: fmt.Errorf("unsupported source: %s. Supported sources are: %v",
				"aaaa",
				[]string{"abc", "bbc", "nbc", "usatoday", "washingtontimes"}),
		},
		{
			Input: []string{"abcbbc"},
			ExpectedOutput: fmt.Errorf("unsupported source: %s. Supported sources are: %v",
				"abcbbc",
				[]string{"abc", "bbc", "nbc", "usatoday", "washingtontimes"}),
		},
	}

	for _, testCase := range testCases {
		err := validateSources(testCase.Input)
		assert.Equal(t, err, testCase.ExpectedOutput)
	}
}
