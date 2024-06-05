package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonNewsToNews(t *testing.T) {
	testCases := []struct {
		Input          []JsonNews
		ExpectedOutput []News
	}{
		{
			Input: []JsonNews{
				{
					Title: "Json article 1",
				},
				{
					Title: "Json article 2",
				},
			},
			ExpectedOutput: []News{
				{
					Title: "Json article 1",
				},
				{
					Title: "Json article 2",
				},
			},
		},
	}

	for _, testCase := range testCases {
		news := JsonNewsToNews(testCase.Input)

		assert.Equal(t, news, testCase.ExpectedOutput)
	}
}
