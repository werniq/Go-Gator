package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractFileData(t *testing.T) {

	testCases := []struct {
		name        string
		filename    string
		expectedNil bool
	}{
		{
			name:        "Successful execution",
			filename:    "sources" + JsonExtension,
			expectedNil: false,
		},
		{
			name:        "Invalid filename",
			filename:    "not-file",
			expectedNil: true,
		},
	}

	for _, tt := range testCases {
		data, err := extractFileData(tt.filename)

		if tt.expectedNil {
			assert.Nil(t, err)
			assert.Equal(t, 0, len(data))
		} else {
			assert.NotEqual(t, 0, len(data))
			assert.Equal(t, err, nil, "Expected err to be nil. Got: ", err)
		}
	}

}
