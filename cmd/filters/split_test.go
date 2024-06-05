package filters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitString(t *testing.T) {
	testCases := []struct {
		Input struct {
			S   string
			Sep string
		}
		ExpectedOutput []string
	}{
		{
			Input: struct {
				S   string
				Sep string
			}{
				S:   "aaa,bbb",
				Sep: ",",
			},
			ExpectedOutput: []string{"aaa", "bbb"},
		},
	}

	for _, testCase := range testCases {
		result := SplitString(testCase.Input.S, testCase.Input.Sep)
		assert.Equal(t, result, testCase.ExpectedOutput)
	}
}
