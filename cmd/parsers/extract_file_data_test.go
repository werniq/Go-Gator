package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractFileData(t *testing.T) {
	filename := "bbc.xml"

	data, err := extractFileData(filename)

	assert.Equal(t, len(data), 41350)
	assert.Equal(t, err, nil, "Expected err to be nil. Got: ", err)
}
