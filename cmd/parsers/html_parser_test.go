package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		Source: UsaToday,
	}

	testCases := []struct {
		Name  string
		Input *types.FilteringParams
	}{
		{Name: "Default run"},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			news, err := parser.Parse()
			assert.NoError(t, err)

			assert.NotNil(t, news)
		})
	}
}
