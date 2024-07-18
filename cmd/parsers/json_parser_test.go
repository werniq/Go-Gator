package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonParser_ParseWithArgs(t *testing.T) {
	parser := JsonParser{
		Source: time.Now().Format(time.DateOnly) + JsonExtension,
	}

	testCases := []struct {
		Name string
	}{
		{
			Name: "Default parse",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			news, err := parser.Parse()
			assert.NoError(t, err)

			assert.NotNil(t, news)
		})
	}
}
