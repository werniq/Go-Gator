package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractFileData(t *testing.T) {
	filename := "sources" + JsonExtension

	data, err := extractFileData(filename)

	assert.NotEqual(t, 0, len(data))
	assert.Equal(t, err, nil, "Expected err to be nil. Got: ", err)
}
