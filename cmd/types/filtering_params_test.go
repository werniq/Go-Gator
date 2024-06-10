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
		}
		ExpectedOutput *FilteringParams
	}{
		{
			Input: struct {
				Keywords   string
				StartingTs string
				EndingTs   string
			}{
				Keywords:   "Ukraine",
				StartingTs: "2024-05-12",
				EndingTs:   "2024-05-15",
			},
			ExpectedOutput: &FilteringParams{
				Keywords:          "Ukraine",
				StartingTimestamp: "2024-05-12",
				EndingTimestamp:   "2024-05-15",
			},
		},
	}

	for _, testCase := range testCases {
		params := NewFilteringParams(
			testCase.Input.Keywords,
			testCase.Input.StartingTs,
			testCase.Input.EndingTs,
		)
		assert.Equal(t, params, testCase.ExpectedOutput)
	}
}
