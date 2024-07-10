package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		Source: UsaToday,
	}

	testCases := []struct {
		name string
	}{
		{
			name: "Default run",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse()
			assert.NoError(t, err)

		})
	}
}
