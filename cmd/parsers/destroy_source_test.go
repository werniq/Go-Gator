package parsers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/types"
	"os"
	"testing"
)

func TestDestroySource(t *testing.T) {
	tests := []struct {
		name         string
		argsDate     string
		argsSource   string
		prepData     []types.Article
		expectedData []types.Article
		expectErr    bool
	}{
		{
			name:       "Successful deletion",
			argsDate:   "2024-07-22",
			argsSource: "bbc",
			prepData: []types.Article{
				{Publisher: "bbc", Title: "Article 1"},
				{Publisher: "cnn", Title: "Article 2"},
			},
			expectedData: []types.Article{
				{Publisher: "cnn", Title: "Article 2"},
			},
			expectErr: false,
		},
		{
			name:       "Invalid date",
			argsDate:   "2024-06-35",
			argsSource: "bbc",
			expectErr:  true,
		},
		{
			name:       "Not-existent data record",
			argsDate:   "2024-06-20",
			argsSource: "bbc",
			expectErr:  true,
		},
		{
			name:       "No articles to delete (empty data)",
			argsDate:   "2024-06-20",
			argsSource: "not-source",
			expectErr:  true,
		},
		{
			name:       "No articles to delete (non-existent source)",
			argsDate:   "2024-07-22",
			argsSource: "new-source",
			prepData: []types.Article{
				{Publisher: "bbc", Title: "Article 1"},
				{Publisher: "abc", Title: "Article 2"},
			},
			expectedData: []types.Article{
				{Publisher: "bbc", Title: "Article 1"},
				{Publisher: "abc", Title: "Article 2"},
			},
			expectErr: false,
		},
		{
			name:         "Run without any data",
			argsDate:     "2024-07-22",
			argsSource:   "bbc",
			prepData:     []types.Article{},
			expectedData: []types.Article{},
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.prepData) > 0 {
				fileData, err := json.Marshal(tt.prepData)
				assert.Nil(t, err)
				err = os.WriteFile(tt.argsDate+JsonExtension, fileData, 0644)
				assert.Nil(t, err)
			}

			err := DestroySource(tt.argsSource, []string{tt.argsDate})
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)

				if len(tt.prepData) > 0 {
					fileData, err := os.ReadFile(tt.argsDate + JsonExtension)
					assert.Nil(t, err)

					var resultData []types.Article
					err = json.Unmarshal(fileData, &resultData)
					assert.Nil(t, err)

					assert.Equal(t, tt.expectedData, resultData)
				}
			}

			err = os.Remove(tt.argsDate + JsonExtension)
			assert.Nil(t, err)
		})
	}
}
