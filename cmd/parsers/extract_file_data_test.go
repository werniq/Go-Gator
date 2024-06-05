package parsers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractFileData(t *testing.T) {
	filename := "bbc.xml"

	data := extractFileData(filename)

	assert.Equal(t, len(data), 41350)
}
