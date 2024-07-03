package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParams(t *testing.T) {
	testCases := []struct {
		Input struct {
			Keywords   string
			StartingTs string
			EndingTs   string
			Sources    string
		}
		ExpectedOutput *FilteringParams
	}{
		{
			Input: struct {
				Keywords   string
				StartingTs string
				EndingTs   string
				Sources    string
			}{
				Keywords:   "Ukraine",
				StartingTs: "2024-05-12",
				EndingTs:   "2024-05-15",
				Sources:    "abc",
			},
			ExpectedOutput: &FilteringParams{
				Keywords:          "Ukraine",
				StartingTimestamp: "2024-05-12",
				EndingTimestamp:   "2024-05-15",
				Sources:           "abc",
			},
		},
	}

	for _, testCase := range testCases {
		params := NewFilteringParams(
			testCase.Input.Keywords,
			testCase.Input.StartingTs,
			testCase.Input.EndingTs,
			testCase.Input.Sources,
		)
		assert.Equal(t, params, testCase.ExpectedOutput)
	}
}
