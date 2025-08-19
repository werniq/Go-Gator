package parsers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/filters"
	"gogator/cmd/types"
	"sync"
	"testing"
)

func TestParseWithParams(t *testing.T) {
	sources := []string{"2024-06-20" + JsonExtension}
	testCases := []struct {
		Input          *types.FilteringParams
		ExpectedOutput int
	}{
		{
			Input: &types.FilteringParams{
				Keywords: "glide",
			},
			ExpectedOutput: 0,
		},
		{
			Input: &types.FilteringParams{
				Keywords: "Ukraine",
			},
			ExpectedOutput: 0,
		},
	}

	for _, testCase := range testCases {
		var news []types.Article
		var err error

		for _, source := range sources {
			news, err = ParseBySource(source)
			assert.Equal(t, err, nil, fmt.Sprintf("Expected: %v, Got: %v", nil, err))
		}
		news = filters.Apply(news, testCase.Input)
		assert.Equal(t, testCase.ExpectedOutput, len(news))
	}
}

func Test_fetchNews(t *testing.T) {

	tests := []struct {
		name        string
		p           Parser
		news        *[]types.Article
		wg          *sync.WaitGroup
		mu          *sync.Mutex
		errChannel  chan error
		expectedErr bool
	}{
		{
			name:        "Successful execution",
			p:           g.XmlParser(sourceToEndpoint[ABC]),
			news:        new([]types.Article),
			wg:          &sync.WaitGroup{},
			mu:          &sync.Mutex{},
			errChannel:  make(chan error, 1),
			expectedErr: false,
		},
		{
			name:        "Incorrect endpoint (bad parser)",
			p:           g.XmlParser(""),
			news:        new([]types.Article),
			wg:          &sync.WaitGroup{},
			mu:          &sync.Mutex{},
			errChannel:  make(chan error, 1),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchNews(tt.p, tt.news, tt.wg, tt.mu, tt.errChannel)

			err := <-tt.errChannel
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
