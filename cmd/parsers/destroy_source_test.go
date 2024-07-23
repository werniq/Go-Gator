package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDestroySource(t *testing.T) {
	tests := []struct {
		name       string
		argsDate   string
		argsSource string
		expectErr  bool
	}{
		{
			name:       "Successful deletion",
			argsDate:   "2024-07-22",
			argsSource: "bbc",
			expectErr:  false,
		},
		{
			name:       "Wrong filename",
			argsDate:   "2024-09-22",
			argsSource: "bbc",
			expectErr:  true,
		},
		{
			name:       "Wrong source",
			argsDate:   "2024-07-22",
			argsSource: "new-source",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DestroySource(tt.argsSource, []string{tt.argsDate})
			assert.Nil(t, err)
		})
	}

}
