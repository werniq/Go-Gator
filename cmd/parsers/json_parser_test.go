package parsers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonParser_ParseWithArgs(t *testing.T) {
	parser := JsonParser{
		Source: "sources.json",
	}

	testCases := []struct {
		Name        string
		setupMock   func()
		expectError bool
	}{
		{
			Name: "Default parse",
			setupMock: func() {
				mockFileContent := `[{"Title":"Test News","Description":"This is a test news.","PubDate":"2024-07-23","Publisher":"Test Source","Link":"http://example.com"}]`
				mockExtractFileData := func(filename string) ([]byte, error) {
					return []byte(mockFileContent), nil
				}
				openFile = mockExtractFileData
			},
			expectError: false,
		},
		{
			Name: "File read failure",
			setupMock: func() {
				mockExtractFileData := func(filename string) ([]byte, error) {
					return nil, errors.New("file read error")
				}
				openFile = mockExtractFileData
			},
			expectError: true,
		},
		{
			Name: "Invalid JSON format",
			setupMock: func() {
				mockFileContent := `[{Title: "Test News", "Description": "Invalid JSON", "PubDate": "2024-07-23", "Publisher": "Test Source", "Link": "http://example.com"}]`
				mockExtractFileData := func(filename string) ([]byte, error) {
					return []byte(mockFileContent), nil
				}
				openFile = mockExtractFileData
			},
			expectError: true,
		},
		{
			Name: "Empty JSON",
			setupMock: func() {
				mockFileContent := `[]`
				mockExtractFileData := func(filename string) ([]byte, error) {
					return []byte(mockFileContent), nil
				}
				openFile = mockExtractFileData
			},
			expectError: false,
		},
		{
			Name: "JSON with unexpected structure",
			setupMock: func() {
				mockFileContent := `[{"UnexpectedField":"Some value"}]`
				mockExtractFileData := func(filename string) ([]byte, error) {
					return []byte(mockFileContent), nil
				}
				openFile = mockExtractFileData
			},
			expectError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testCase.setupMock()
			news, err := parser.Parse()

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, news)
			}
		})
	}
}
